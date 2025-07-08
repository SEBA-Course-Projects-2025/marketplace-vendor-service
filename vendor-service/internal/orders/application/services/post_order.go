package services

import (
	"context"
	"gorm.io/gorm"
	eventDomain "marketplace-vendor-service/vendor-service/internal/event/domain"
	orderDomain "marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func PostOrder(ctx context.Context, orderRepo orderDomain.OrderRepository, eventRepo eventDomain.EventRepository, db *gorm.DB, orderReq dtos.OrderCreatedEventResponseDto) error {

	ctx, span := tracer.Tracer.Start(ctx, "PostOrder")
	defer span.End()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txOrderRepo := orderRepo.WithTx(tx)
		txEventRepo := eventRepo.WithTx(tx)

		newOrder := dtos.CreateOrderWithDto(orderReq)

		if orderReq.Items == nil || len(orderReq.Items) == 0 {
			outbox, err := dtos.OrderStatusToOutbox(newOrder, "vendor.updated.order", "vendor.order.events")

			if err != nil {
				return err
			}

			if err = txEventRepo.CreateOutboxRecord(ctx, outbox); err != nil {
				return err
			}

			return nil

		}

		order, err := txOrderRepo.Create(ctx, newOrder, orderReq.Items[0].VendorId)

		if err != nil {
			return err
		}

		outbox, err := dtos.OrderStatusToOutbox(order, "vendor.updated.order", "vendor.order.events")

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
