package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "marketplace-vendor-service/docs"
	accountInterfaces "marketplace-vendor-service/vendor-service/internal/account/interfaces"
	accountHandlers "marketplace-vendor-service/vendor-service/internal/account/interfaces/handlers"
	orderInterfaces "marketplace-vendor-service/vendor-service/internal/orders/interfaces"
	orderHandlers "marketplace-vendor-service/vendor-service/internal/orders/interfaces/handlers"
	reviewInterfaces "marketplace-vendor-service/vendor-service/internal/reviews/interfaces"
	reviewHandlers "marketplace-vendor-service/vendor-service/internal/reviews/interfaces/handlers"
	"marketplace-vendor-service/vendor-service/internal/shared/middlewares"
)

func SetUpRouter(accountHandler *accountHandlers.AccountHandler, reviewHandler *reviewHandlers.ReviewHandler, orderHandler *orderHandlers.OrderHandler) *gin.Engine {

	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("auth/login", accountHandler.LoginAccountHandler)

	api := r.Group("/api", middlewares.AuthMiddleware())

	accountInterfaces.SetUpAccountsRouter(api, accountHandler)
	reviewInterfaces.SetUpReviewsRouter(api, reviewHandler)
	orderInterfaces.SetUpOrdersRouter(api, orderHandler)

	return r
}
