package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"
	"marketplace-vendor-service/vendor-service/internal/orders/application/services"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/orders/interfaces/handlers"
)

type Consumer struct {
	AMQPChannel  *amqp.Channel
	OrderHandler *handlers.OrderHandler
}

func NewConsumer(channel *amqp.Channel, orderHandler *handlers.OrderHandler) *Consumer {
	return &Consumer{
		AMQPChannel:  channel,
		OrderHandler: orderHandler,
	}
}
func (c *Consumer) StartConsuming(ctx context.Context, queueName string) error {

	msgs, err := c.AMQPChannel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				msgId, err := uuid.Parse(msg.MessageId)
				if err != nil {
					log.Printf("Invalid message id: %v", err)
					_ = msg.Nack(false, false)
					continue
				}

				processed, err := c.OrderHandler.EventRepo.CheckProcessedMessage(ctx, msgId)
				if err != nil {
					log.Printf("Error checking processed message: %v", err)
					_ = msg.Nack(false, false)
					continue
				}
				if processed {
					_ = msg.Ack(false)
					continue
				}

				if err = c.processMessage(msg, ctx); err != nil {
					log.Printf("Error processing new message: %v", err)
					_ = msg.Nack(false, false)
					continue
				}

				if err = c.OrderHandler.EventRepo.CreateProcessedMessage(ctx, msgId); err != nil {
					log.Printf("Error addign processed message: %v", err)
					_ = msg.Nack(false, false)
					continue
				}

				_ = msg.Ack(false)
			}
		}
	}()

	return nil

}

func (c *Consumer) processMessage(msg amqp.Delivery, ctx context.Context) error {

	eventType := msg.Type

	if eventType == "" {
		return fmt.Errorf("event type not found")
	}

	switch eventType {

	case "order.created.vendor":

		var newOrderDto dtos.OrderCreatedEventDto

		if err := json.Unmarshal(msg.Body, &newOrderDto); err != nil {
			return err
		}

		if err := services.CheckOrderProductsQuantity(ctx, newOrderDto, c.OrderHandler.EventRepo); err != nil {
			return err
		}

	case "vendor.product.quantity.checked":

		var createdOrderDto dtos.OrderCreatedEventResponseDto

		if err := json.Unmarshal(msg.Body, &createdOrderDto); err != nil {
			return err
		}

		if err := services.PostOrder(ctx, c.OrderHandler.OrderRepo, c.OrderHandler.EventRepo, c.OrderHandler.Db, createdOrderDto); err != nil {
			return err
		}

	case "order.canceled.vendor":

		var canceledOrderDto dtos.CanceledOrderEventDto

		if err := json.Unmarshal(msg.Body, &canceledOrderDto); err != nil {
			return err
		}

		if err := services.CancelOrderStatus(ctx, c.OrderHandler.OrderRepo, c.OrderHandler.EventRepo, c.OrderHandler.Db, canceledOrderDto); err != nil {
			return err
		}

	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}

	return nil

}
