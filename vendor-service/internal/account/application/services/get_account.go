package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func GetAccount(ctx context.Context, repo domain.AccountRepository, vendorId uuid.UUID) (dtos.AccountResponse, error) {

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Starting GetAccount application service")

	ctx, span := tracer.Tracer.Start(ctx, "GetAccount")
	defer span.End()

	account, err := repo.FindById(ctx, vendorId)

	if err != nil {
		return dtos.AccountResponse{}, nil
	}

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Successfully get account by vendorId")

	return dtos.AccountToDto(account), nil
}
