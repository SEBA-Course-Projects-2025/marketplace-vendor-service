package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
)

func PutOrderStatus(ctx context.Context, orderRepo domain.OrderRepository, statusReq dtos.StatusRequestDto, id uuid.UUID, vendorId uuid.UUID) error {
	return orderRepo.Transaction(func(txRepo domain.OrderRepository) error {

		existingOrder, err := txRepo.FindById(ctx, id, vendorId)

		if err != nil {
			return err
		}

		existingOrder.Status = statusReq.Status
		existingOrder.VendorId = vendorId

		return txRepo.Update(ctx, existingOrder)
	})
}
