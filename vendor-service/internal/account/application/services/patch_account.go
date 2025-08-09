package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func PatchAccount(ctx context.Context, accountRepo domain.AccountRepository, accountReq dtos.AccountPatchRequest, vendorId uuid.UUID) (dtos.AccountResponse, error) {

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Starting PatchAccount application service")

	ctx, span := tracer.Tracer.Start(ctx, "PatchAccount")
	defer span.End()

	var accountResponse dtos.AccountResponse

	if err := accountRepo.Transaction(func(txRepo domain.AccountRepository) error {

		existingAccount, err := txRepo.FindById(ctx, vendorId)

		if err != nil {
			return err
		}

		existingAccount = dtos.PatchDtoToAccount(existingAccount, accountReq)

		existingAccount.Id = vendorId

		existingAccount, err = txRepo.Patch(ctx, existingAccount)

		if err != nil {
			return err
		}

		accountResponse = dtos.AccountToDto(existingAccount)

		return nil

	}); err != nil {
		return dtos.AccountResponse{}, err
	}

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Successfully partially modified account by vendorId")

	return accountResponse, nil
}
