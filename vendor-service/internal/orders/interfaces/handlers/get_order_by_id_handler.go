package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/orders/application/services"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"net/http"
)

// GetOrderByIdHandler godoc
// @Summary      Get order by ID
// @Description  Returns a single order by its ID for the given vendor.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        orderId path string true "Order ID (UUID)"
// @Success      200 {object} dtos.OneOrderResponse
// @Failure      400 {object} map[string]interface{} "Invalid vendorId or orderId"
// @Failure      404 {object} map[string]interface{} "Order not found"
// @Failure      500 {object} map[string]interface{}
// @Router       /orders/{orderId} [get]
func (h *OrderHandler) GetOrderByIdHandler(c *gin.Context) {

	ctx, span := tracer.Tracer.Start(c.Request.Context(), "GetOrderByIdHandler")
	defer span.End()

	orderIdStr := c.Param("orderId")

	orderId, err := uuid.Parse(orderIdStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order UUID"})
		return
	}

	order, err := services.GetOrderById(ctx, h.OrderRepo, orderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)

}
