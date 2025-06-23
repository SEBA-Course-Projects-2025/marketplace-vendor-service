package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
)

func PutAccount(ctx context.Context, accountRepo domain.AccountRepository, accountReq dtos.PutRequestDto, vendorId uuid.UUID) error {

	return accountRepo.Transaction(func(txRepo domain.AccountRepository) error {
		existingAccount, err := txRepo.FindById(ctx, vendorId)

		if err != nil {
			return err
		}

		*existingAccount = dtos.UpdateAccountWithDto(accountReq, existingAccount)

		return txRepo.Update(ctx, existingAccount)
	})

}
