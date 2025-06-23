package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/reviews/application/services"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"net/http"
)

// PostReplyHandler godoc
// @Summary      Add reply to review
// @Description  Adds a reply to a specific review for the given vendor.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        reviewId path string true "Review ID (UUID)"
// @Param        commentReq body dtos.CommentDto true "Reply data"
// @Success      201 {object} dtos.PostReplyDto
// @Failure      400 {object} map[string]interface{} "Invalid vendorId, reviewId, or reply data"
// @Failure      500 {object} map[string]interface{}
// @Router       /reviews/{reviewId}/replies [post]
func (h *ReviewHandler) PostReplyHandler(c *gin.Context) {

	v, _ := c.Get("vendorId")
	vendorId, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendorId"})
		return
	}

	reviewIdStr := c.Param("reviewId")

	reviewId, err := uuid.Parse(reviewIdStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review UUID"})
		return
	}

	var commentReq dtos.CommentDto

	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reply data"})
		return
	}

	newReply, err := services.PostReply(c.Request.Context(), h.ReviewRepo, h.AccountRepo, h.Db, commentReq, vendorId, reviewId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newReply)

}
