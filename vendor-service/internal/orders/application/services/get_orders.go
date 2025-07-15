package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func GetOrders(ctx context.Context, orderRepo domain.OrderRepository, params dtos.OrderQueryParams, vendorId uuid.UUID) ([]dtos.GetOrdersResponse, error) {

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Starting GetOrders application service")

	ctx, span := tracer.Tracer.Start(ctx, "GetOrders")
	defer span.End()

	orders, err := orderRepo.FindAll(ctx, params, vendorId)

	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Successfully get orders by vendorId")

	return dtos.OrdersToDto(orders), nil

}
