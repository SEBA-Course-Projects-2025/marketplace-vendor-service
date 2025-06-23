package dtos

import (
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/account/domain/account_models"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain/models"
	"time"
)

type ReviewQueryParams struct {
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

type GetReviewsResponse struct {
	ReviewId     uuid.UUID `json:"review_id"`
	VendorId     uuid.UUID `json:"vendor_id"`
	ProductId    uuid.UUID `json:"product_id"`
	ReviewerId   uuid.UUID `json:"reviewer_id"`
	ReviewerName string    `json:"reviewer_name"`
	Rating       float32   `json:"rating"`
	Comment      string    `json:"comment"`
	Date         time.Time `json:"date"`
}

func GetReviewsToDto(reviews []models.Review) []GetReviewsResponse {

	var reviewsResponse []GetReviewsResponse

	for _, review := range reviews {
		reviewResponse := GetReviewsResponse{
			ReviewId:     review.Id,
			VendorId:     review.VendorId,
			ProductId:    review.ProductId,
			ReviewerId:   review.ReviewerId,
			ReviewerName: review.ReviewerName,
			Rating:       review.Rating,
			Comment:      review.Comment,
			Date:         review.CreatedAt,
		}
		reviewsResponse = append(reviewsResponse, reviewResponse)
	}

	return reviewsResponse
}

type ReplyDto struct {
	ReplyId   uuid.UUID `json:"reply_id"`
	ReplierId uuid.UUID `json:"replier_id"`
	Name      string    `json:"name"`
	Comment   string    `json:"comment"`
	Date      time.Time `json:"date"`
}

type GetOneReviewDto struct {
	ReviewId     uuid.UUID  `json:"review_id"`
	VendorId     uuid.UUID  `json:"vendor_id"`
	ProductId    uuid.UUID  `json:"product_id"`
	ReviewerId   uuid.UUID  `json:"reviewer_id"`
	ReviewerName string     `json:"reviewer_name"`
	Rating       float32    `json:"rating"`
	Comment      string     `json:"comment"`
	Date         time.Time  `json:"date"`
	Replies      []ReplyDto `json:"replies"`
}

func ReviewToDto(review *models.Review) GetOneReviewDto {

	var repliesDto []ReplyDto

	for _, reply := range review.Replies {
		replyDto := ReplyDto{
			ReplyId:   reply.Id,
			ReplierId: reply.ReplierId,
			Name:      reply.ReplierName,
			Comment:   reply.Comment,
			Date:      reply.CreatedAt,
		}
		repliesDto = append(repliesDto, replyDto)
	}

	return GetOneReviewDto{
		ReviewId:     review.Id,
		VendorId:     review.VendorId,
		ProductId:    review.ProductId,
		ReviewerId:   review.ReviewerId,
		ReviewerName: review.ReviewerName,
		Rating:       review.Rating,
		Comment:      review.Comment,
		Date:         review.CreatedAt,
		Replies:      repliesDto,
	}
}

type PostReplyDto struct {
	ReplyId   uuid.UUID `json:"reply_id"`
	ReplierId uuid.UUID `json:"replier_id"`
	ReviewId  uuid.UUID `json:"review_id"`
	Name      string    `json:"name"`
	Comment   string    `json:"comment"`
	Date      time.Time `json:"date"`
}

type PostReplyRequestDto struct {
	ReplierId uuid.UUID `json:"replier_id"`
	ReviewId  uuid.UUID `json:"review_id"`
	Name      string    `json:"name"`
	Comment   string    `json:"comment"`
}

func PostNewReplyWithDto(commentReq CommentDto, vendor *account_models.VendorAccount, reviewId uuid.UUID) *models.Reply {
	return &models.Reply{
		Id:          uuid.New(),
		ReviewId:    reviewId,
		ReplierId:   vendor.Id,
		ReplierName: vendor.Name,
		Comment:     commentReq.Comment,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func PostReplyToReplyDto(newReply *models.Reply) PostReplyDto {
	return PostReplyDto{
		ReplyId:   newReply.Id,
		ReplierId: newReply.ReplierId,
		ReviewId:  newReply.ReviewId,
		Name:      newReply.ReplierName,
		Comment:   newReply.Comment,
		Date:      newReply.CreatedAt,
	}
}

type CommentDto struct {
	Comment string `json:"comment"`
}
