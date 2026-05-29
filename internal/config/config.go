package config

import (
	"errors"
	"os"
	"strings"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort string
	NotifyGRPCAddr string
	LogLevel string
	PostgresDSN string
	RedisAddr string
	RedisPass string
	KafkaBrokers []string
	KafkaTopic string
	KafkaAlertTopic string
	KafkaGroup string
	
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	
	dsn := os.Getenv("POSTGRES_DSN")

	if dsn == "" {
		return nil, errors.New("POSTGRES_DSN environment variable is required")
	}
	
	brokersRaw := getEnv("KAFKA_BROKERS", "localhost:9092")


	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "50054"),
		NotifyGRPCAddr: getEnv("NOTIFY_GRPC_ADDR", "localhost:50055"),
		LogLevel: getEnv("LOG_LEVEL", "WARN"),
		
		PostgresDSN: dsn,
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: getEnv("REDIS_PASSWORD", ""),
		KafkaBrokers: strings.Split(brokersRaw, ","),
		KafkaTopic: getEnv("KAFKA_TOPIC", "machine.telemetry"),
		KafkaAlertTopic: getEnv("KAFKA_ALERT_TOPIC", "machine.alerts"),
		KafkaGroup: getEnv("KAFKA_GROUP", "telemetry-consumer-group"),
	}, nil
}
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		
		return val
	}
	return fallback
	
}
