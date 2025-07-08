package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func GetAllReviews(ctx context.Context, reviewRepo domain.ReviewRepository, params dtos.ReviewQueryParams, vendorId uuid.UUID) ([]dtos.GetReviewsResponse, error) {

	ctx, span := tracer.Tracer.Start(ctx, "GetAllReviews")
	defer span.End()

	reviews, err := reviewRepo.FindAll(ctx, params, vendorId)

	if err != nil {
		return nil, err
	}

	return dtos.GetReviewsToDto(reviews), nil

}
