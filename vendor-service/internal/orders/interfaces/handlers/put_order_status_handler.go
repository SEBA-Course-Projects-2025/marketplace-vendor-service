package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/orders/application/services"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"net/http"
)

// PutOrderStatusHandler godoc
// @Summary      Update order status
// @Description  Updates the status of an order for the given vendor.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        orderId path string true "Order ID (UUID)"
// @Param        statusReq body dtos.StatusRequestDto true "Order status update data"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "Invalid vendorId, orderId, or request body"
// @Failure      404 {object} map[string]interface{} "Order not found"
// @Failure      500 {object} map[string]interface{}
// @Router       /orders/{orderId} [put]
func (h *OrderHandler) PutOrderStatusHandler(c *gin.Context) {

	v, _ := c.Get("vendorId")
	vendorId, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendorId"})
		return
	}

	orderIdStr := c.Param("orderId")

	orderId, err := uuid.Parse(orderIdStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order UUID"})
		return
	}

	var statusReq dtos.StatusRequestDto

	if err := c.ShouldBindJSON(&statusReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
		return
	}

	if err = services.PutOrderStatus(c.Request.Context(), h.OrderRepo, statusReq, orderId, vendorId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}
