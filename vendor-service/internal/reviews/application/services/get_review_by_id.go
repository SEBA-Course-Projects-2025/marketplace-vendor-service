package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func GetReviewById(ctx context.Context, reviewRepo domain.ReviewRepository, id uuid.UUID, vendorId uuid.UUID) (dtos.GetOneReviewDto, error) {

	logrus.WithFields(logrus.Fields{
		"reviewId": id,
	}).Info("Starting GetReviewById application service")

	ctx, span := tracer.Tracer.Start(ctx, "GetReviewById")
	defer span.End()

	review, err := reviewRepo.FindById(ctx, id, vendorId)

	if err != nil {
		return dtos.GetOneReviewDto{}, err
	}

	logrus.WithFields(logrus.Fields{
		"reviewId": id,
	}).Info("Successfully get review by id")

	return dtos.ReviewToDto(review), nil
}
