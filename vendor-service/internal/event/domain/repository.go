package domain

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/event/domain/models"
)

type EventRepository interface {
	CreateOutboxRecord(ctx context.Context, outbox *models.Outbox) error
	FetchUnprocessed(ctx context.Context) ([]models.Outbox, error)
	MarkProcessed(ctx context.Context, id uuid.UUID) error
	CheckProcessedMessage(ctx context.Context, id uuid.UUID) (bool, error)
	CreateProcessedMessage(ctx context.Context, id uuid.UUID) error
	Transaction(fn func(txRepo EventRepository) error) error
	WithTx(tx *gorm.DB) EventRepository
}
