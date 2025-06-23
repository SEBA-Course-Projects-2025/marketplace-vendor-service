package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
)

func PatchAccount(ctx context.Context, accountRepo domain.AccountRepository, accountReq dtos.AccountPatchRequest, vendorId uuid.UUID) (dtos.AccountResponse, error) {

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

	return accountResponse, nil
}
