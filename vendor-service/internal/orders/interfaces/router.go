package interfaces

import (
	"github.com/gin-gonic/gin"
	"marketplace-vendor-service/vendor-service/internal/orders/interfaces/handlers"
)

func SetUpOrdersRouter(rg *gin.RouterGroup, h *handlers.OrderHandler) {
	orders := rg.Group("orders")
	{
		orders.GET("/", h.GetOrdersHandler)

		orders.GET("/:orderId", h.GetOrderByIdHandler)
		orders.PATCH("/:orderId", h.PatchOrderStatusHandler)
	}
}
