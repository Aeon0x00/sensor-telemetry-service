include .env
export

export PROJECT_ROOT=$(shell pwd)

.PHONY: proto generate lint format

tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate:
	buf generate
	cd gen && go mod tidy

lint:
	buf lint

format:
	buf format -w

env-up:
	docker compose up -d telemetry-postgres redis kafka

env-down:
	docker compose down 

env-cleanup:
	@echo "ВНИМАНИЕ: Опасность утери данных!"
	@echo -n "Очистить все volume файлы окружения? [y/N]: " && read ans && \
	if [ "$$ans" = "y" ]; then \
		docker compose down telemetry-postgres && \
		rm -rf _out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Нет параметра 'seq'. Используйте: make migrate-create seq=имя_миграции"; \
		exit 1; \
	fi
	docker compose run --rm telemetry-postgres-migrate \
	    create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"
migrate-action:
	docker compose run --rm telemetry-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@telemetry-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		$(action)

migrate-up:
	@make migrate-action action=up
migrate-down:
	@make migrate-action action=down

run:
	@go run cmd/app/main.go
