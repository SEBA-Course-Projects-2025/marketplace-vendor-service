package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
)

func GetReviewById(ctx context.Context, reviewRepo domain.ReviewRepository, id uuid.UUID, vendorId uuid.UUID) (dtos.GetOneReviewDto, error) {
	review, err := reviewRepo.FindById(ctx, id, vendorId)

	if err != nil {
		return dtos.GetOneReviewDto{}, err
	}

	return dtos.ReviewToDto(review), nil
}
