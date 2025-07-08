package main

import (
	"context"
	"fmt"
	"log"
	accountRepository "marketplace-vendor-service/vendor-service/internal/account/infrastructure/repository"
	accountHandlers "marketplace-vendor-service/vendor-service/internal/account/interfaces/handlers"
	eventInfrastructure "marketplace-vendor-service/vendor-service/internal/event/infrastructure"
	eventRepository "marketplace-vendor-service/vendor-service/internal/event/infrastructure/repository"
	orderRepository "marketplace-vendor-service/vendor-service/internal/orders/infrastructure/repository"
	orderHandlers "marketplace-vendor-service/vendor-service/internal/orders/interfaces/handlers"
	reviewRepository "marketplace-vendor-service/vendor-service/internal/reviews/infrastructure/repository"
	reviewHandlers "marketplace-vendor-service/vendor-service/internal/reviews/interfaces/handlers"
	"marketplace-vendor-service/vendor-service/internal/shared/amqp"
	handlers "marketplace-vendor-service/vendor-service/internal/shared/handler"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"

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

	newTracer := tracer.InitTracer()
	defer func() {
		if err := newTracer(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer: %v", err)
		}
	}()

	accountRepo := accountRepository.New(dbUsed)
	reviewRepo := reviewRepository.New(dbUsed)
	orderRepo := orderRepository.New(dbUsed)
	eventRepo := eventRepository.New(dbUsed)

	amqpConfig, err := amqp.ConnectAMQP()
	if err != nil {
		log.Fatalln(err)
	}

	outboxPoller := eventInfrastructure.NewOutboxPoller(eventRepo, amqpConfig.Channel, time.Second*2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := outboxPoller.StartPolling(ctx, amqpConfig.ConfirmationChannel); err != nil {
			log.Fatalln(err)
		}
	}()

	sharedHandler := &handlers.Handler{
		AccountRepo: accountRepo,
		ReviewRepo:  reviewRepo,
		OrderRepo:   orderRepo,
		EventRepo:   eventRepo,
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

	consumer := eventInfrastructure.NewConsumer(amqpConfig.Channel, orderHandler)

	queues := []string{
		"order.created.vendor",
		"vendor.product.quantity.checked",
		"order.canceled.vendor",
	}

	for _, queue := range queues {
		go func(q string) {
			if err := consumer.StartConsuming(ctx, queue); err != nil {
				log.Fatalf("Consumer error: %v", err)
			}
		}(queue)
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
