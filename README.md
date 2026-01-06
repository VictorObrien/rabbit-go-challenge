# RabbitMQ Go Challenge - Sistema de Processamento Assíncrono

Sistema de processamento assíncrono em Go com RabbitMQ, implementando retries com backoff, DLQ, idempotência e observabilidade.

## Requisitos

- Go 1.23+
- Docker & Docker Compose
- Make (opcional)

## Estrutura do Projeto

- `services/publisher`: API HTTP para publicação de tarefas
- `services/consumer`: Worker para processamento assíncrono
- `pkg/schema`: Estruturas de dados compartilhadas
- `pkg/log`: Logger estruturado compartilhado

## Como Rodar

### 1. Subir RabbitMQ

```bash
cd deployments
docker-compose up -d rabbitmq
```

Acesse o management: http://localhost:15672 (rabbit/rabbit)

### 2. Rodar localmente (desenvolvimento)

```bash
# Terminal 1 - Publisher
cd services/publisher
go run main.go

# Terminal 2 - Consumer
cd services/consumer
go run main.go
```

### 3. Rodar com Docker Compose (completo)

```bash
cd deployments
docker-compose up --build
```

## Arquitetura

```
[API] -> [RabbitMQ] -> [Consumer]
           |              |
           v              v
       [Retry Queues]  [SQLite]
           |
           v
        [DLQ]
```

## Status do Projeto

- [x] Etapa 1: Setup inicial
- [x] Etapa 2: Estruturas de dados
- [x] Etapa 3: Topologia RabbitMQ
- [ ] Etapa 4: Publisher API
- [ ] Etapa 5: Consumer básico
- [ ] Etapa 6: Sistema de retries
- [ ] Etapa 7: Idempotência
- [ ] Etapa 8: Observabilidade

