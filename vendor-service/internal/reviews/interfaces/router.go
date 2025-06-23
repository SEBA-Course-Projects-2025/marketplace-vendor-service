package interfaces

import (
	"github.com/gin-gonic/gin"
	"marketplace-vendor-service/vendor-service/internal/reviews/interfaces/handlers"
)

func SetUpReviewsRouter(rg *gin.RouterGroup, h *handlers.ReviewHandler) {
	reviews := rg.Group("reviews")
	{
		reviews.GET("/", h.GetAllReviewsHandler)
		reviews.POST("/:reviewId/replies", h.PostReplyHandler)

		reviews.GET("/:reviewId", h.GetReviewHandler)
		reviews.PUT("/:reviewId/replies/:replyId", h.PutReplyHandler)

	}
}
