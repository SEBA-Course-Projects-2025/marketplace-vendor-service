package services

import (
	"context"
	"gorm.io/gorm"
	eventDomain "marketplace-vendor-service/vendor-service/internal/event/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/metrics"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func CancelOrderStatus(ctx context.Context, orderRepo domain.OrderRepository, eventRepo eventDomain.EventRepository, db *gorm.DB, canceledOrderDto dtos.CanceledOrderEventDto) error {

	ctx, span := tracer.Tracer.Start(ctx, "CancelOrderStatus")
	defer span.End()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txOrderRepo := orderRepo.WithTx(tx)
		txEventRepo := eventRepo.WithTx(tx)

		existingOrder, err := txOrderRepo.FindById(ctx, canceledOrderDto.OrderId)

		if err != nil {
			return err
		}

		existingOrder.Status = canceledOrderDto.Status

		updatedOrder, err := txOrderRepo.Patch(ctx, existingOrder)
		if err != nil {
			return err
		}

		outbox, err := dtos.OrderStatusToOutbox(updatedOrder, "vendor.updated.order", "vendor.order.events")

		if err != nil {
			return err
		}

		metrics.OrderStatusUpdatedCounter.WithLabelValues(canceledOrderDto.Status).Inc()

		err = txEventRepo.CreateOutboxRecord(ctx, outbox)

		if err != nil {
			return err
		}

		canceledOrderProducts := dtos.CanceledOrderItemsToDto(updatedOrder)

		outbox, err = dtos.CanceledOrderProductsToOutbox(canceledOrderProducts, "vendor.cancel.product.order", "vendor.product.events")

		if err != nil {
			return err
		}

		err = txEventRepo.CreateOutboxRecord(ctx, outbox)

		if err != nil {
			return err
		}

		return nil

	})

}
