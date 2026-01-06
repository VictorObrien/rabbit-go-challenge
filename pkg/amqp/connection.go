package amqp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/VictorObrien/rabbitmq-go-challenge/pkg/log"
)

// ConnectionConfig configurações de conexão
type ConnectionConfig struct {
	URL             string
	ReconnectDelay  time.Duration
	MaxReconnects   int
	ChannelPoolSize int
}

// DefaultConnectionConfig retorna configuração padrão
func DefaultConnectionConfig(url string) *ConnectionConfig {
	return &ConnectionConfig{
		URL:             url,
		ReconnectDelay:  5 * time.Second,
		MaxReconnects:   10,
		ChannelPoolSize: 1,
	}
}

// Connection wrapper para conexão AMQP
type Connection struct {
	conn   *amqp.Connection
	config *ConnectionConfig
	logger *slog.Logger
}

// NewConnection cria uma nova conexão
func NewConnection(config *ConnectionConfig) (*Connection, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c := &Connection{
		conn:   conn,
		config: config,
		logger: log.GetLogger().With("component", "amqp"),
	}

	c.logger.Info("connected to RabbitMQ", "url", maskURL(config.URL))
	return c, nil
}

// Channel cria um novo canal
func (c *Connection) Channel() (*amqp.Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	return ch, nil
}

// Close fecha a conexão
func (c *Connection) Close() error {
	if c.conn != nil && !c.conn.IsClosed() {
		c.logger.Info("closing RabbitMQ connection")
		return c.conn.Close()
	}
	return nil
}

// IsClosed verifica se a conexão está fechada
func (c *Connection) IsClosed() bool {
	return c.conn == nil || c.conn.IsClosed()
}

// NotifyClose registra um listener para eventos de fechamento
func (c *Connection) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	return c.conn.NotifyClose(receiver)
}

// maskURL mascara a senha na URL para logs
func maskURL(url string) string {
	// Implementação simples - em produção, use uma lib de parse de URL
	return "amqp://***:***@***"
}

// SetupWithRetry tenta conectar e declarar topologia com retries
func SetupWithRetry(ctx context.Context, config *ConnectionConfig, topologyConfig *TopologyConfig) (*Connection, error) {
	var conn *Connection
	var err error

	for attempt := 0; attempt < config.MaxReconnects; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		conn, err = NewConnection(config)
		if err != nil {
			log.Warn("connection attempt failed",
				"attempt", attempt+1,
				"max", config.MaxReconnects,
				"error", err,
			)
			time.Sleep(config.ReconnectDelay)
			continue
		}

		// Declarar topologia
		ch, err := conn.Channel()
		if err != nil {
			conn.Close()
			log.Warn("failed to create channel", "error", err)
			time.Sleep(config.ReconnectDelay)
			continue
		}

		if err := DeclareTopology(ch, topologyConfig); err != nil {
			ch.Close()
			conn.Close()
			log.Warn("failed to declare topology", "error", err)
			time.Sleep(config.ReconnectDelay)
			continue
		}

		ch.Close()
		log.Info("RabbitMQ setup completed successfully")
		return conn, nil
	}

	return nil, fmt.Errorf("failed to setup RabbitMQ after %d attempts: %w",
		config.MaxReconnects, err)
}


