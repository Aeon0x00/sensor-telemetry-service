# machine-telemetry-service

Сервис сбора потоковой телеметрии и мониторинга состояния промышленного оборудования (температура, вибрация, обороты).

## Стек (Stack)
- Go 1.23 (gRPC, Protobuf, pgx/v5, kafka-go, go-redis/v9, zap)
- Apache Kafka (Consumer Group + Alert Producer)
- Redis (идемпотентность событий SetNX, счетчик подряд идущих перегрузок, кэш состояния станка)
- PostgreSQL + миграции (migrate)

## Команды (Makefile)
- make env-up
- make migrate-up
- make run
