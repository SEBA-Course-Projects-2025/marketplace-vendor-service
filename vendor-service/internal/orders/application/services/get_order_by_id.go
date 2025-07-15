package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func GetOrderById(ctx context.Context, orderRepo domain.OrderRepository, id uuid.UUID) (dtos.OneOrderResponse, error) {

	logrus.WithFields(logrus.Fields{
		"orderId": id,
	}).Info("Starting GetOrderById application service")

	ctx, span := tracer.Tracer.Start(ctx, "GetOrderById")
	defer span.End()

	order, err := orderRepo.FindById(ctx, id)

	if err != nil {
		return dtos.OneOrderResponse{}, err
	}

	logrus.WithFields(logrus.Fields{
		"orderId": id,
	}).Info("Successfully get order by id")

	return dtos.OrderToDto(order), nil
}
