package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
)

func GetAllReviews(ctx context.Context, reviewRepo domain.ReviewRepository, params dtos.ReviewQueryParams, vendorId uuid.UUID) ([]dtos.GetReviewsResponse, error) {

	reviews, err := reviewRepo.FindAll(ctx, params, vendorId)

	if err != nil {
		return nil, err
	}

	return dtos.GetReviewsToDto(reviews), nil

}
