package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/reviews/application/services"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"net/http"
)

// PutReplyHandler godoc
// @Summary      Update reply to review
// @Description  Updates a reply to a specific review for the given vendor.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        reviewId path string true "Review ID (UUID)"
// @Param        replyId path string true "Reply ID (UUID)"
// @Param        comment body dtos.CommentDto true "Reply update data"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "Invalid vendorId, reviewId, replyId, or comment data"
// @Failure      404 {object} map[string]interface{} "Reply not found"
// @Failure      500 {object} map[string]interface{}
// @Router       /reviews/{reviewId}/replies/{replyId} [put]
func (h *ReviewHandler) PutReplyHandler(c *gin.Context) {

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

	replyIdStr := c.Param("replyId")

	replyId, err := uuid.Parse(replyIdStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reply UUID"})
		return
	}

	var comment dtos.CommentDto

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment data"})
		return
	}

	err = services.PutReply(c.Request.Context(), h.ReviewRepo, comment, replyId, reviewId, vendorId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reply not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}
