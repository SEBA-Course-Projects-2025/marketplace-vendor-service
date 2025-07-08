package services

import (
	"context"
	eventDomain "marketplace-vendor-service/vendor-service/internal/event/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func CheckOrderProductsQuantity(ctx context.Context, newOrderDto dtos.OrderCreatedEventDto, eventRepo eventDomain.EventRepository) error {

	ctx, span := tracer.Tracer.Start(ctx, "CheckOrderProductsQuantity")
	defer span.End()

	return eventRepo.Transaction(func(txRepo eventDomain.EventRepository) error {

		outbox, err := dtos.OrderItemsToOutbox(newOrderDto, "vendor.check.product.quantity", "vendor.product.events")

		if err != nil {
			return err
		}

		err = eventRepo.CreateOutboxRecord(ctx, outbox)

		if err != nil {
			return err
		}

		return nil

	})
}
