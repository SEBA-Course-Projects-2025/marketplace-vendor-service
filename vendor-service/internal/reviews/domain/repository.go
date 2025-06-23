package domain

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain/models"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
)

type ReviewRepository interface {
	FindAll(ctx context.Context, params dtos.ReviewQueryParams, vendorId uuid.UUID) ([]models.Review, error)
	FindById(ctx context.Context, id uuid.UUID, vendorId uuid.UUID) (*models.Review, error)
	Create(ctx context.Context, newReply *models.Reply) (*models.Reply, error)
	Update(ctx context.Context, updatedReply *models.Reply) error
	Transaction(fn func(txRepo ReviewRepository) error) error
	WithTx(tx *gorm.DB) ReviewRepository
}
