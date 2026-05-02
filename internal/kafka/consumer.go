package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"machine-telemetry-service/internal/domain"
	"machine-telemetry-service/internal/usecase"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TelemetryConsumer struct {
	reader 	*kafka.Reader
	us     	*usecase.UseCase
	logger 	*zap.SugaredLogger
}

func NewTelemetryConsumer (brokers []string, topic, groupId string, us *usecase.UseCase,logger *zap.SugaredLogger) *TelemetryConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupId,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	
	return &TelemetryConsumer{
		reader: r,
		us:     us,
		logger: logger,
	}
}

// func (c *TelemetryConsumer) StartReadSimple(ctx context.Context) {
// 	for {
// 		m, err := c.reader.ReadMessage(ctx)
// 		if err != nil {
// 			break
// 		}
// 		fmt.Println("msg:", string(m.Value))
// 	}
// }

func (c *TelemetryConsumer) Start(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			c.logger.Errorf("kafka fetch error %v",err)
			continue
		}
		fmt.Println("read offset:", msg.Offset)


		var dto domain.KafkaTelemetryDTO
		if err := json.Unmarshal(msg.Value, &dto); err != nil {
			c.logger.Errorw(
				"failed to decode kafka message",
				"offset", 	msg.Offset,
				"error",  	err,
			)
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		if err := c.us.ProcessKafkaTelemetry(ctx, dto); err != nil {
			c.logger.Errorw(
				"failed to process telemetry event",
				"eventId", 	dto.EventId,
				"error",   	err,
			)
			continue
		}


		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Errorf("failed to commit offset %v",err)
		}
	}
}

func (c *TelemetryConsumer) Close() error {
	return c.reader.Close()
}
