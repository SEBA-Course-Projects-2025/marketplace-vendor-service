package services

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	eventDomain "marketplace-vendor-service/vendor-service/internal/event/domain"
	orderDomain "marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/metrics"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
)

func PostOrder(ctx context.Context, orderRepo orderDomain.OrderRepository, eventRepo eventDomain.EventRepository, db *gorm.DB, orderReq dtos.OrderCreatedEventResponseDto) error {

	logrus.WithFields(logrus.Fields{
		"items": orderReq.Items,
	}).Info("Starting PostOrder application service")

	ctx, span := tracer.Tracer.Start(ctx, "PostOrder")
	defer span.End()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txOrderRepo := orderRepo.WithTx(tx)
		txEventRepo := eventRepo.WithTx(tx)

		newOrder := dtos.CreateOrderWithDto(orderReq)

		if orderReq.Items == nil || len(orderReq.Items) == 0 {
			outbox, err := dtos.OrderStatusToOutbox(newOrder, "vendor.updated.order", "vendor.order.events")

			if err != nil {
				return error_handler.ErrorHandler(err, err.Error())
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

		metrics.OrdersAddedCounter.Inc()

		metrics.OrderStatusUpdatedCounter.WithLabelValues(orderReq.Status).Inc()

		outbox, err := dtos.OrderStatusToOutbox(order, "vendor.updated.order", "vendor.order.events")

		if err != nil {
			return error_handler.ErrorHandler(err, err.Error())
		}

		err = txEventRepo.CreateOutboxRecord(ctx, outbox)

		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"items": orderReq.Items,
		}).Info("Successfully added new order")

		return nil

	})

}
