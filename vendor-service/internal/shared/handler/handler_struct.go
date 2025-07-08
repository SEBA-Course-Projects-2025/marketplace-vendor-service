package handlers

import (
	"gorm.io/gorm"
	accountDomain "marketplace-vendor-service/vendor-service/internal/account/domain"
	eventDomain "marketplace-vendor-service/vendor-service/internal/event/domain"
	orderDomain "marketplace-vendor-service/vendor-service/internal/orders/domain"
	reviewDomain "marketplace-vendor-service/vendor-service/internal/reviews/domain"
)

type Handler struct {
	AccountRepo accountDomain.AccountRepository

	ReviewRepo reviewDomain.ReviewRepository

	OrderRepo orderDomain.OrderRepository

	EventRepo eventDomain.EventRepository

	Db *gorm.DB
}
