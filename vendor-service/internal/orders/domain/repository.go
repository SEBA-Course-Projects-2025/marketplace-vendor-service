package domain

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/orders/domain/models"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
)

type OrderRepository interface {
	FindAll(ctx context.Context, params dtos.OrderQueryParams, vendorId uuid.UUID) ([]models.Order, error)
	FindById(ctx context.Context, id uuid.UUID, vendorId uuid.UUID) (*models.Order, error)
	Update(ctx context.Context, updatedOrder *models.Order) error
	Transaction(fn func(txRepo OrderRepository) error) error
	WithTx(tx *gorm.DB) OrderRepository
}
