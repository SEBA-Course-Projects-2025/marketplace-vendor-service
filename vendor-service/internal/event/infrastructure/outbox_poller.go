package infrastructure

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"marketplace-vendor-service/vendor-service/internal/event/domain"
	"time"
)

type OutboxPoller struct {
	OutboxRepo  domain.EventRepository
	AMQPChannel *amqp.Channel
	Interval    time.Duration
}

func NewOutboxPoller(repo domain.EventRepository, channel *amqp.Channel, interval time.Duration) *OutboxPoller {
	return &OutboxPoller{
		OutboxRepo:  repo,
		AMQPChannel: channel,
		Interval:    interval,
	}
}

func (op *OutboxPoller) processOutbox(ctx context.Context, confirmationsChannel <-chan amqp.Confirmation) error {

	records, err := op.OutboxRepo.FetchUnprocessed(ctx)
	if err != nil {
		return err
	}

	for _, record := range records {
		
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := op.AMQPChannel.Publish(record.Exchange, record.EventType, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         record.Payload,
			MessageId:    record.Id.String(),
			Timestamp:    time.Now(),
			Type:         record.EventType,
			DeliveryMode: amqp.Persistent,
		}); err != nil {
			return err
		}

		select {
		case confirmed := <-confirmationsChannel:
			if !confirmed.Ack {
				return fmt.Errorf("message not acknowledged by broker")
			}
		case <-ctx.Done():
			return ctx.Err()
		}

		if err := op.OutboxRepo.MarkProcessed(ctx, record.Id); err != nil {
			return err
		}
	}

	return nil

}

func (op *OutboxPoller) StartPolling(ctx context.Context, confirmationsChannel <-chan amqp.Confirmation) error {
	ticker := time.NewTicker(op.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := op.processOutbox(ctx, confirmationsChannel); err != nil {
				log.Printf("Error during processing event records: %v", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
