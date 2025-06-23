package handlers

import (
	"gorm.io/gorm"
	accountDomain "marketplace-vendor-service/vendor-service/internal/account/domain"
	orderDomain "marketplace-vendor-service/vendor-service/internal/orders/domain"
	reviewDomain "marketplace-vendor-service/vendor-service/internal/reviews/domain"
)

type Handler struct {
	AccountRepo accountDomain.AccountRepository

	ReviewRepo reviewDomain.ReviewRepository

	OrderRepo orderDomain.OrderRepository

	Db *gorm.DB
}
