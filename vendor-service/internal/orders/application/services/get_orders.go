package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
)

func GetOrders(ctx context.Context, orderRepo domain.OrderRepository, params dtos.OrderQueryParams, vendorId uuid.UUID) ([]dtos.GetOrdersResponse, error) {

	orders, err := orderRepo.FindAll(ctx, params, vendorId)

	if err != nil {
		return nil, err
	}

	return dtos.OrdersToDto(orders), nil

}
