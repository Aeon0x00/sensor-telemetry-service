package grpcserver

import (
	"machine-telemetry-service/gen/telemetry"
	"machine-telemetry-service/internal/usecase"

	"go.uber.org/zap"
)

type Server struct {
	telemetry.UnimplementedTelemetryServiceServer
	us     *usecase.UseCase
	logger *zap.SugaredLogger
}

func NewServer (us *usecase.UseCase,logger *zap.SugaredLogger) *Server {
	return &Server{
		us:     us,
		logger: logger,
	}
}
