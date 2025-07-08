package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/orders/application/services"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"net/http"
)

// PatchOrderStatusHandler godoc
// @Summary      Partially update order status
// @Description  Partially updates the status of an order for the given vendor.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer access token"
// @Param        orderId path string true "Order ID (UUID)"
// @Param        statusReq body dtos.StatusRequestDto true "Order status update data"
// @Success      200 {object} dtos.OneOrderResponse
// @Failure      400 {object} map[string]interface{} "Invalid vendorId, orderId, or request body"
// @Failure      404 {object} map[string]interface{} "Order not found"
// @Failure      500 {object} map[string]interface{}
// @Router       /orders/{orderId} [patch]
func (h *OrderHandler) PatchOrderStatusHandler(c *gin.Context) {

	ctx, span := tracer.Tracer.Start(c.Request.Context(), "PatchOrderStatusHandler")
	defer span.End()

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

	order, err := services.PatchOrderStatus(ctx, h.OrderRepo, h.EventRepo, h.Db, statusReq, orderId, vendorId)
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
