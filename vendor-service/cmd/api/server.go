package main

import (
	"fmt"
	"log"
	accountRepository "marketplace-vendor-service/vendor-service/internal/account/infrastructure/repository"
	accountHandlers "marketplace-vendor-service/vendor-service/internal/account/interfaces/handlers"
	orderRepository "marketplace-vendor-service/vendor-service/internal/orders/infrastructure/repository"
	orderHandlers "marketplace-vendor-service/vendor-service/internal/orders/interfaces/handlers"
	reviewRepository "marketplace-vendor-service/vendor-service/internal/reviews/infrastructure/repository"
	reviewHandlers "marketplace-vendor-service/vendor-service/internal/reviews/interfaces/handlers"
	handlers "marketplace-vendor-service/vendor-service/internal/shared/handler"

	"marketplace-vendor-service/vendor-service/internal/shared/db"
	"marketplace-vendor-service/vendor-service/internal/shared/router"
	"os"
	"time"
)

// @title Vendor Service API
// @version 1.0
// @description API for managing vendor account.

// @schemes https
// @host marketplace-vendor-service.onrender.com
// @BasePath /api
func main() {

	dbUsed, err := db.ConnectDb()
	if err != nil {
		log.Fatalln(err)
	}

	accountRepo := accountRepository.New(dbUsed)
	reviewRepo := reviewRepository.New(dbUsed)
	orderRepo := orderRepository.New(dbUsed)

	sharedHandler := &handlers.Handler{
		AccountRepo: accountRepo,
		ReviewRepo:  reviewRepo,
		OrderRepo:   orderRepo,
		Db:          dbUsed,
	}

	accountHandler := &accountHandlers.AccountHandler{
		Handler: sharedHandler,
	}

	reviewHandler := &reviewHandlers.ReviewHandler{
		Handler: sharedHandler,
	}

	orderHandler := &orderHandlers.OrderHandler{
		Handler: sharedHandler,
	}

	mainRouter := router.SetUpRouter(accountHandler, reviewHandler, orderHandler)

	port := os.Getenv("API_PORT")

	if port == "" {
		port = "8081"
	}

	fmt.Println(time.Now())

	if err := mainRouter.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
