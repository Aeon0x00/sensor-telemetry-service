package app

import (
	"context"
	"machine-telemetry-service/gen/notify"
	"machine-telemetry-service/gen/telemetry"
	"machine-telemetry-service/internal/config"
	"machine-telemetry-service/internal/database"
	"machine-telemetry-service/internal/grpcclient"
	"machine-telemetry-service/internal/grpcserver"
	"machine-telemetry-service/internal/kafka"
	"machine-telemetry-service/internal/repository"
	"machine-telemetry-service/internal/usecase"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run (cfg *config.Config, zapLog *zap.SugaredLogger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool,err := database.NewPostgresPool(ctx, cfg.PostgresDSN)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	notifyConn,err := grpc.NewClient(cfg.NotifyGRPCAddr,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer notifyConn.Close()
	notifyClient := grpcclient.NewNotifyClient(notify.NewNotifyServiceClient(notifyConn),zapLog)



	pgRepo := repository.NewPostgresRepo(dbPool,zapLog)
	redisRepo := repository.NewRedisRepo(cfg.RedisAddr,cfg.RedisPass,zapLog)
	alertProd := kafka.NewAlertProducer(cfg.KafkaBrokers,cfg.KafkaAlertTopic,zapLog)
	defer alertProd.Close()

	us := usecase.NewUseCase(pgRepo,redisRepo,notifyClient,alertProd,zapLog)

	consumer := kafka.NewTelemetryConsumer(cfg.KafkaBrokers,cfg.KafkaTopic,cfg.KafkaGroup,us,zapLog)
	defer consumer.Close()

	go func() {
		zapLog.Infof("kafka consumer listening topic: %s", cfg.KafkaTopic)
		if err := consumer.Start(ctx); err != nil {
			zapLog.Errorf("kafka consumer error: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return err
	}
	
	grpcSrv := grpc.NewServer()
	telemetry.RegisterTelemetryServiceServer(grpcSrv,grpcserver.NewServer(us,zapLog))

	go func() {
		zapLog.Infof("machine-telemetry-service listening on :%s", cfg.GRPCPort)
		if err := grpcSrv.Serve(lis); err != nil {
			zapLog.Errorf("grpc serve error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
	grpcSrv.GracefulStop()
	return  nil
}
