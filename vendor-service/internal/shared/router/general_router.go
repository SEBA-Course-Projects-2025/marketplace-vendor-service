package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	_ "marketplace-vendor-service/docs"
	accountInterfaces "marketplace-vendor-service/vendor-service/internal/account/interfaces"
	accountHandlers "marketplace-vendor-service/vendor-service/internal/account/interfaces/handlers"
	orderInterfaces "marketplace-vendor-service/vendor-service/internal/orders/interfaces"
	orderHandlers "marketplace-vendor-service/vendor-service/internal/orders/interfaces/handlers"
	reviewInterfaces "marketplace-vendor-service/vendor-service/internal/reviews/interfaces"
	reviewHandlers "marketplace-vendor-service/vendor-service/internal/reviews/interfaces/handlers"
	"marketplace-vendor-service/vendor-service/internal/shared/middlewares"
	"os"
)

func SetUpRouter(accountHandler *accountHandlers.AccountHandler, reviewHandler *reviewHandlers.ReviewHandler, orderHandler *orderHandlers.OrderHandler) *gin.Engine {

	r := gin.New()

	p := ginprometheus.NewPrometheus("vendor_service")
	p.MetricsPath = ""
	p.Use(r)

	r.GET("/metrics", gin.BasicAuth(gin.Accounts{
		os.Getenv("METRICS_ACCESS_USERNAME"): os.Getenv("METRICS_ACCESS_PASSWORD"),
	}), gin.WrapH(promhttp.Handler()))

	r.Use(otelgin.Middleware("vendor_service"))

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-length"},
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("auth/login", accountHandler.LoginAccountHandler)

	api := r.Group("/api", middlewares.AuthMiddleware())

	accountInterfaces.SetUpAccountsRouter(api, accountHandler)
	reviewInterfaces.SetUpReviewsRouter(api, reviewHandler)
	orderInterfaces.SetUpOrdersRouter(api, orderHandler)

	return r
}
