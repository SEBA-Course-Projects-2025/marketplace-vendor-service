package services

import (
	"context"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain/models"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
)

func PatchReply(ctx context.Context, reviewRepo domain.ReviewRepository, comment dtos.CommentDto, replyId uuid.UUID, reviewId uuid.UUID, vendorId uuid.UUID) (dtos.PostReplyDto, error) {

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
			return err
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

	return replyResponse, nil
}
