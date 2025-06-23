package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/reviews/application/services"
	"net/http"
)

// GetReviewHandler godoc
// @Summary      Get review by ID
// @Description  Returns a single review by its ID for the given vendor.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        reviewId path string true "Review ID (UUID)"
// @Success      200 {object} dtos.GetOneReviewDto
// @Failure      400 {object} map[string]interface{} "Invalid vendorId or reviewId"
// @Failure      404 {object} map[string]interface{} "Review not found"
// @Failure      500 {object} map[string]interface{}
// @Router       /reviews/{reviewId} [get]
func (h *ReviewHandler) GetReviewHandler(c *gin.Context) {

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

	review, err := services.GetReviewById(c.Request.Context(), h.ReviewRepo, reviewId, vendorId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)

}
