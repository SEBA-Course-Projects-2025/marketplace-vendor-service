package services

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	accountDomain "marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
)

func PostReply(ctx context.Context, reviewRepo domain.ReviewRepository, accountRepo accountDomain.AccountRepository, db *gorm.DB, commentReq dtos.CommentDto, vendorId uuid.UUID, reviewId uuid.UUID) (dtos.PostReplyDto, error) {

	ctx, span := tracer.Tracer.Start(ctx, "PostReply")
	defer span.End()

	var replyResponse dtos.PostReplyDto

	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txAccountRepo := accountRepo.WithTx(tx)
		txReviewRepo := reviewRepo.WithTx(tx)

		vendor, err := txAccountRepo.FindById(ctx, vendorId)

		if err != nil {
			return err
		}

		newReply := dtos.PostNewReplyWithDto(commentReq, vendor, reviewId)

		reply, err := txReviewRepo.Create(ctx, newReply)

		if err != nil {
			return err
		}

		replyResponse = dtos.PostReplyToReplyDto(reply)

		return nil

	}); err != nil {
		return dtos.PostReplyDto{}, err
	}

	return replyResponse, nil
}
