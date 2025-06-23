package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
)

func GetOrderById(ctx context.Context, orderRepo domain.OrderRepository, id uuid.UUID, vendorId uuid.UUID) (dtos.OneOrderResponse, error) {

	order, err := orderRepo.FindById(ctx, id, vendorId)

	if err != nil {
		return dtos.OneOrderResponse{}, err
	}

	return dtos.OrderToDto(order), nil
}
