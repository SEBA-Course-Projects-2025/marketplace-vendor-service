package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
)

func GetAccount(ctx context.Context, repo domain.AccountRepository, vendorId uuid.UUID) (dtos.AccountResponse, error) {

	account, err := repo.FindById(ctx, vendorId)

	if err != nil {
		return dtos.AccountResponse{}, nil
	}

	return dtos.AccountToDto(account), nil
}
