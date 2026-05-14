package main

import (
	"log"
	"machine-telemetry-service/internal/app"
	"machine-telemetry-service/internal/config"
	"machine-telemetry-service/internal/logger"
)

func main () {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	zapLog,closelog,err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer closelog()
	
	
	if err := app.Run(cfg, zapLog); err != nil {
		log.Fatalf("app stopped with error: %v", err)
	}
}
