package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain/models"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
)

func PatchReply(ctx context.Context, reviewRepo domain.ReviewRepository, comment dtos.CommentDto, replyId uuid.UUID, reviewId uuid.UUID, vendorId uuid.UUID) (dtos.PostReplyDto, error) {

	logrus.WithFields(logrus.Fields{
		"reviewId": reviewId,
		"replyId":  replyId,
		"vendorId": vendorId,
	}).Info("Starting PatchReply application service")

	ctx, span := tracer.Tracer.Start(ctx, "PatchReply")
	defer span.End()

	var replyResponse dtos.PostReplyDto

	if err := reviewRepo.Transaction(func(txRepo domain.ReviewRepository) error {

		existingReview, err := txRepo.FindById(ctx, reviewId, vendorId)

		if err != nil {
			return err
		}

		var existingReply *models.Reply
		for i := range existingReview.Replies {
			if existingReview.Replies[i].Id == replyId {
				existingReply = &existingReview.Replies[i]
				break
			}
		}

		if existingReply == nil {
			return error_handler.ErrorHandler(gorm.ErrRecordNotFound, gorm.ErrRecordNotFound.Error())
		}

		existingReply.Comment = comment.Comment

		updatedReply, err := txRepo.Patch(ctx, existingReply)

		if err != nil {
			return err
		}

		replyResponse = dtos.PostReplyToReplyDto(updatedReply)

		return nil

	}); err != nil {
		return dtos.PostReplyDto{}, err
	}

	logrus.WithFields(logrus.Fields{
		"reviewId": reviewId,
		"replyId":  replyId,
		"vendorId": vendorId,
	}).Info("Successfully patched reply by its reviewId, replyId and vendorId")

	return replyResponse, nil
}
