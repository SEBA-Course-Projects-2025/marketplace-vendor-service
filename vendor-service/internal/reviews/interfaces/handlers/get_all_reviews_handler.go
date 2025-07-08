package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/reviews/application/services"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"net/http"
	"strconv"
)

// GetAllReviewsHandler godoc
// @Summary      List reviews
// @Description  Returns a paginated list of reviews for the given vendor.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        page query int false "Page number (default 1)"
// @Param        size query int false "Page size (default 15)"
// @Param        offset query int false "Offset (overrides page/size if set)"
// @Param        limit query int false "Limit (overrides size if set)"
// @Param        sortBy query string false "Sort by field (default date)"
// @Param        sortOrder query string false "Sort order: asc or desc (default asc)"
// @Success      200 {array} dtos.GetReviewsResponse
// @Failure      400 {object} map[string]interface{} "Invalid query parameters or vendorId"
// @Failure      500 {object} map[string]interface{}
// @Router       /reviews [get]
func (h *ReviewHandler) GetAllReviewsHandler(c *gin.Context) {

	ctx, span := tracer.Tracer.Start(c.Request.Context(), "GetAllReviewsHandler")
	defer span.End()

	v, _ := c.Get("vendorId")
	vendorId, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendorId"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "15"))
	if err != nil || size < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "-1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "-1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	if offset < 0 || limit <= 0 {
		offset = (page - 1) * size
		limit = size
	}

	sortBy := c.DefaultQuery("sortBy", "date")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	queryParams := dtos.ReviewQueryParams{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	reviews, err := services.GetAllReviews(ctx, h.ReviewRepo, queryParams, vendorId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)

}
