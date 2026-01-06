package amqp

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Exchanges
	TasksExchange = "tasks.exchange"
	RetryExchange = "retry.exchange"
	DLQExchange   = "dlq.exchange"

	// Queues
	MainQueue     = "tasks.main"
	Retry5sQueue  = "tasks.retry.5s"
	Retry30sQueue = "tasks.retry.30s"
	Retry5mQueue  = "tasks.retry.5m"
	DLQQueue      = "tasks.dlq"

	// Routing Keys
	TasksRoutingKey = "tasks"
	DLQRoutingKey   = "tasks.dlq"
)

// TopologyConfig configurações da topologia
type TopologyConfig struct {
	// TTLs em milissegundos
	Retry5sTTL  int32
	Retry30sTTL int32
	Retry5mTTL  int32
}

// DefaultTopologyConfig retorna a configuração padrão
func DefaultTopologyConfig() *TopologyConfig {
	return &TopologyConfig{
		Retry5sTTL:  5000,   // 5 segundos
		Retry30sTTL: 30000,  // 30 segundos
		Retry5mTTL:  300000, // 5 minutos
	}
}

// DeclareTopology declara toda a topologia RabbitMQ
func DeclareTopology(ch *amqp.Channel, config *TopologyConfig) error {
	if config == nil {
		config = DefaultTopologyConfig()
	}

	// 1. Declarar Exchanges
	if err := declareExchanges(ch); err != nil {
		return fmt.Errorf("failed to declare exchanges: %w", err)
	}

	// 2. Declarar Queues
	if err := declareQueues(ch, config); err != nil {
		return fmt.Errorf("failed to declare queues: %w", err)
	}

	// 3. Criar Bindings
	if err := createBindings(ch); err != nil {
		return fmt.Errorf("failed to create bindings: %w", err)
	}

	return nil
}

// declareExchanges declara todos os exchanges
func declareExchanges(ch *amqp.Channel) error {
	exchanges := []string{
		TasksExchange,
		RetryExchange,
		DLQExchange,
	}

	for _, exchange := range exchanges {
		err := ch.ExchangeDeclare(
			exchange, // name
			"direct", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", exchange, err)
		}
	}

	return nil
}

// declareQueues declara todas as queues com suas configurações
func declareQueues(ch *amqp.Channel, config *TopologyConfig) error {
	// Main Queue - com DLX para retry
	_, err := ch.QueueDeclare(
		MainQueue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-dead-letter-exchange": RetryExchange,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare main queue: %w", err)
	}

	// Retry 5s Queue - volta para main após TTL
	_, err = ch.QueueDeclare(
		Retry5sQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":          config.Retry5sTTL,
			"x-dead-letter-exchange": TasksExchange,
			"x-dead-letter-routing-key": TasksRoutingKey,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare retry 5s queue: %w", err)
	}

	// Retry 30s Queue
	_, err = ch.QueueDeclare(
		Retry30sQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":          config.Retry30sTTL,
			"x-dead-letter-exchange": TasksExchange,
			"x-dead-letter-routing-key": TasksRoutingKey,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare retry 30s queue: %w", err)
	}

	// Retry 5m Queue
	_, err = ch.QueueDeclare(
		Retry5mQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":          config.Retry5mTTL,
			"x-dead-letter-exchange": TasksExchange,
			"x-dead-letter-routing-key": TasksRoutingKey,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare retry 5m queue: %w", err)
	}

	// DLQ - sem TTL, mensagens ficam aqui permanentemente
	_, err = ch.QueueDeclare(
		DLQQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	return nil
}

// createBindings cria todos os bindings entre exchanges e queues
func createBindings(ch *amqp.Channel) error {
	bindings := []struct {
		queue      string
		exchange   string
		routingKey string
	}{
		// Main queue recebe de tasks.exchange
		{MainQueue, TasksExchange, TasksRoutingKey},

		// Retry queues recebem de retry.exchange
		{Retry5sQueue, RetryExchange, Retry5sQueue},
		{Retry30sQueue, RetryExchange, Retry30sQueue},
		{Retry5mQueue, RetryExchange, Retry5mQueue},

		// DLQ recebe de dlq.exchange
		{DLQQueue, DLQExchange, DLQRoutingKey},
	}

	for _, b := range bindings {
		err := ch.QueueBind(
			b.queue,      // queue name
			b.routingKey, // routing key
			b.exchange,   // exchange
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue %s to exchange %s: %w",
				b.queue, b.exchange, err)
		}
	}

	return nil
}

// GetRetryQueue retorna o nome da fila de retry baseado na tentativa
func GetRetryQueue(attempt int) string {
	switch attempt {
	case 0:
		return Retry5sQueue
	case 1:
		return Retry30sQueue
	case 2:
		return Retry5mQueue
	default:
		return "" // Vai para DLQ
	}
}

// GetRetryExchangeAndKey retorna o exchange e routing key para retry
func GetRetryExchangeAndKey(attempt int) (exchange, routingKey string) {
	queue := GetRetryQueue(attempt)
	if queue == "" {
		return DLQExchange, DLQRoutingKey
	}
	return RetryExchange, queue
}


