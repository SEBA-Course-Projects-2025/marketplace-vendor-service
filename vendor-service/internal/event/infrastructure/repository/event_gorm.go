package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"marketplace-vendor-service/vendor-service/internal/event/domain"
	"marketplace-vendor-service/vendor-service/internal/event/domain/models"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
	"time"
)

type GormEventRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GormEventRepository {
	return &GormEventRepository{db: db}
}

func (ger *GormEventRepository) CreateOutboxRecord(ctx context.Context, outbox *models.Outbox) error {

	if err := ger.db.WithContext(ctx).Create(outbox).Error; err != nil {
		return error_handler.ErrorHandler(err, "Error creating new event record")
	}

	return nil

}

func (ger *GormEventRepository) FetchUnprocessed(ctx context.Context) ([]models.Outbox, error) {

	var outboxRecords []models.Outbox

	if err := ger.db.WithContext(ctx).Where("processed = false").Find(&outboxRecords).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error getting event records data")
	}

	return outboxRecords, nil

}

func (ger *GormEventRepository) MarkProcessed(ctx context.Context, id uuid.UUID) error {

	res := ger.db.WithContext(ctx).Model(&models.Outbox{}).Where("id = ?", id).Updates(map[string]interface{}{"processed": true, "processed_at": time.Now()})

	if res.Error != nil {
		return error_handler.ErrorHandler(res.Error, "Error updating event record")
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}

func (ger *GormEventRepository) CheckProcessedMessage(ctx context.Context, id uuid.UUID) (bool, error) {

	var processedMessage models.ProcessedMessage

	if err := ger.db.WithContext(ctx).Where("message_id = ?", id).First(&processedMessage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, error_handler.ErrorHandler(err, "Error getting processed message record")
	}

	return true, nil

}

func (ger *GormEventRepository) CreateProcessedMessage(ctx context.Context, id uuid.UUID) error {

	message := models.ProcessedMessage{MessageId: id}

	if err := ger.db.WithContext(ctx).Create(message).Error; err != nil {
		return error_handler.ErrorHandler(err, "Error adding new processed message")
	}

	return nil

}

func (ger *GormEventRepository) WithTx(tx *gorm.DB) domain.EventRepository {
	return &GormEventRepository{
		db: tx,
	}
}

func (ger *GormEventRepository) Transaction(fn func(txRepo domain.EventRepository) error) error {

	tx := ger.db.Begin()
	if tx.Error != nil {
		log.Printf("Transaction begin error: %v", tx.Error)
		return tx.Error
	}

	repo := ger.WithTx(tx)

	if err := fn(repo); err != nil {
		log.Printf("Transaction function error: %v", err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit error: %v", err)
		return err
	}

	return nil

}
