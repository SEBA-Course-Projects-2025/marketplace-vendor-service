package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func GetAllReviews(ctx context.Context, reviewRepo domain.ReviewRepository, params dtos.ReviewQueryParams, vendorId uuid.UUID) ([]dtos.GetReviewsResponse, error) {

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Starting GetAllReviews application service")

	ctx, span := tracer.Tracer.Start(ctx, "GetAllReviews")
	defer span.End()

	reviews, err := reviewRepo.FindAll(ctx, params, vendorId)

	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"vendorId": vendorId,
	}).Info("Successfully get reviews by vendorId")

	return dtos.GetReviewsToDto(reviews), nil

}
