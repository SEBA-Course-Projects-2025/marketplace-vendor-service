package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
)

func PatchOrderStatus(ctx context.Context, orderRepo domain.OrderRepository, statusReq dtos.StatusRequestDto, id uuid.UUID, vendorId uuid.UUID) (dtos.OneOrderResponse, error) {

	var orderResponse dtos.OneOrderResponse

	if err := orderRepo.Transaction(func(txRepo domain.OrderRepository) error {

		existingOrder, err := txRepo.FindById(ctx, id, vendorId)

		if err != nil {
			return err
		}

		existingOrder.Status = statusReq.Status
		existingOrder.VendorId = vendorId

		updatedOrder, err := txRepo.Patch(ctx, existingOrder)
		if err != nil {
			return err
		}

		orderResponse = dtos.OrderToDto(updatedOrder)

		return nil

	}); err != nil {
		return dtos.OneOrderResponse{}, err
	}
	return orderResponse, nil
}
