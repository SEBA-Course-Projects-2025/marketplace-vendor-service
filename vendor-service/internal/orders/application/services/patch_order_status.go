package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	eventDomain "marketplace-vendor-service/vendor-service/internal/event/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
)

func PatchOrderStatus(ctx context.Context, orderRepo domain.OrderRepository, eventRepo eventDomain.EventRepository, db *gorm.DB, statusReq dtos.StatusRequestDto, id uuid.UUID, vendorId uuid.UUID) (dtos.OneOrderResponse, error) {

	logrus.WithFields(logrus.Fields{
		"orderId":  id,
		"vendorId": vendorId,
	}).Info("Starting PatchOrderStatus application service")

	ctx, span := tracer.Tracer.Start(ctx, "PatchOrderStatus")
	defer span.End()

	var orderResponse dtos.OneOrderResponse

	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txOrderRepo := orderRepo.WithTx(tx)
		txEventRepo := eventRepo.WithTx(tx)

		existingOrder, err := txOrderRepo.FindById(ctx, id)

		if err != nil {
			return err
		}

		existingOrder.Status = statusReq.Status
		existingOrder.VendorId = vendorId

		updatedOrder, err := txOrderRepo.Patch(ctx, existingOrder)
		if err != nil {
			return err
		}

		outbox, err := dtos.OrderStatusToOutbox(updatedOrder, "vendor.updated.order", "vendor.order.events")

		if err != nil {
			return error_handler.ErrorHandler(err, err.Error())
		}

		err = txEventRepo.CreateOutboxRecord(ctx, outbox)

		if err != nil {
			return err
		}

		orderResponse = dtos.OrderToDto(updatedOrder)

		return nil

	}); err != nil {
		return dtos.OneOrderResponse{}, error_handler.ErrorHandler(err, err.Error())
	}

	logrus.WithFields(logrus.Fields{
		"orderId":  id,
		"vendorId": vendorId,
	}).Info("Successfully partially modified order status by orderId and vendorId")

	return orderResponse, nil
}
