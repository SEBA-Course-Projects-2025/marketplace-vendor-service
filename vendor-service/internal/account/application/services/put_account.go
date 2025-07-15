package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func PutAccount(ctx context.Context, accountRepo domain.AccountRepository, accountReq dtos.PutRequestDto, vendorId uuid.UUID) error {

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Starting PutAccount application service")

	ctx, span := tracer.Tracer.Start(ctx, "PutAccount")
	defer span.End()

	return accountRepo.Transaction(func(txRepo domain.AccountRepository) error {
		existingAccount, err := txRepo.FindById(ctx, vendorId)

		if err != nil {
			return err
		}

		*existingAccount = dtos.UpdateAccountWithDto(accountReq, existingAccount)

		logrus.WithFields(logrus.Fields{
			"vendorId": vendorId,
		}).Info("Successfully fully modified account by vendorId")

		return txRepo.Update(ctx, existingAccount)
	})

}
