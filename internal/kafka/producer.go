package kafka

import (
	"context"
	"encoding/json"
	"machine-telemetry-service/internal/domain"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type AlertProducer struct {
	writer 	*kafka.Writer
	logger 	*zap.SugaredLogger
}

func NewAlertProducer (brokers []string, topic string, logger *zap.SugaredLogger) *AlertProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}
	
	return &AlertProducer{
		writer: w,
		logger: logger,
	}
}

func (p *AlertProducer) PublishAlert(ctx context.Context, alert domain.KafkaAlertDTO) error {
	b, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(alert.MachineId),
		Value: b,
	})
	if err != nil {
		p.logger.Errorf("kafka write alert error %v",err)
		return err
	}
	
	return  nil
}

func (p *AlertProducer) Close() error {
	return p.writer.Close()
}
