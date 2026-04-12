package grpcclient

import (
	"context"
	"machine-telemetry-service/gen/notify"
	"machine-telemetry-service/internal/domain"

	"go.uber.org/zap"
)

type NotifyClient struct {
	client 	notify.NotifyServiceClient
	logger 	*zap.SugaredLogger
}

func NewNotifyClient (client notify.NotifyServiceClient, logger *zap.SugaredLogger) *NotifyClient {
	return &NotifyClient{
		client: client,
		logger: logger,
	}
}

func (nc *NotifyClient) SendOverloadAlert(ctx context.Context, machine *domain.Machine, inc domain.Incident) error {
	nc.client.SendAlert(ctx, mapAlertToProto(machine, inc))
	return nil
}
