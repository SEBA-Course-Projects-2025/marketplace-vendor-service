package services

import (
	"context"
	"github.com/sirupsen/logrus"
	eventDomain "marketplace-vendor-service/vendor-service/internal/event/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
)

func CheckOrderProductsQuantity(ctx context.Context, newOrderDto dtos.OrderCreatedEventDto, eventRepo eventDomain.EventRepository) error {

	logrus.WithFields(logrus.Fields{
		"orderId": newOrderDto.OrderId,
		"items":   newOrderDto.Items,
	}).Info("Starting CheckOrderProductsQuantity application service")

	ctx, span := tracer.Tracer.Start(ctx, "CheckOrderProductsQuantity")
	defer span.End()

	return eventRepo.Transaction(func(txRepo eventDomain.EventRepository) error {

		outbox, err := dtos.OrderItemsToOutbox(newOrderDto, "vendor.check.product.quantity", "vendor.product.events")

		if err != nil {
			return error_handler.ErrorHandler(err, err.Error())
		}

		err = eventRepo.CreateOutboxRecord(ctx, outbox)

		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"orderId": newOrderDto.OrderId,
			"items":   newOrderDto.Items,
		}).Info("Successfully send event to check products quantities for the order")

		return nil

	})
}
