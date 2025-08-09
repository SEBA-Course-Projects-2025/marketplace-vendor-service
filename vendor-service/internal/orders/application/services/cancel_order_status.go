package services

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	eventDomain "marketplace-vendor-service/vendor-service/internal/event/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/metrics"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
)

func CancelOrderStatus(ctx context.Context, orderRepo domain.OrderRepository, eventRepo eventDomain.EventRepository, db *gorm.DB, canceledOrderDto dtos.CanceledOrderEventDto) error {

	logrus.WithFields(logrus.Fields{
		"orderId": canceledOrderDto.OrderId,
	}).Info("Starting CancelOrderStatus application service")

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
			return error_handler.ErrorHandler(err, err.Error())
		}

		metrics.OrderStatusUpdatedCounter.WithLabelValues(canceledOrderDto.Status).Inc()

		err = txEventRepo.CreateOutboxRecord(ctx, outbox)

		if err != nil {
			return err
		}

		canceledOrderProducts := dtos.CanceledOrderItemsToDto(updatedOrder)

		outbox, err = dtos.CanceledOrderProductsToOutbox(canceledOrderProducts, "vendor.cancel.product.order", "vendor.product.events")

		if err != nil {
			return error_handler.ErrorHandler(err, err.Error())
		}

		err = txEventRepo.CreateOutboxRecord(ctx, outbox)

		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"orderId": canceledOrderDto.OrderId,
		}).Info("Successfully canceled order")

		return nil

	})

}
