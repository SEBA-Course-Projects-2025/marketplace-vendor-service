package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func GetAccount(ctx context.Context, repo domain.AccountRepository, vendorId uuid.UUID) (dtos.AccountResponse, error) {

	ctx, span := tracer.Tracer.Start(ctx, "GetAccount")
	defer span.End()

	account, err := repo.FindById(ctx, vendorId)

	if err != nil {
		return dtos.AccountResponse{}, nil
	}

	return dtos.AccountToDto(account), nil
}
