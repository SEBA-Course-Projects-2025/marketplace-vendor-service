package repository

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain"
	"marketplace-vendor-service/vendor-service/internal/reviews/domain/models"
	"marketplace-vendor-service/vendor-service/internal/reviews/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
)

type GormReviewRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GormReviewRepository {
	return &GormReviewRepository{db: db}
}

func (grr *GormReviewRepository) FindAll(ctx context.Context, params dtos.ReviewQueryParams, vendorId uuid.UUID) ([]models.Review, error) {

	var reviews []models.Review

	db := grr.db.WithContext(ctx).Where("vendor_id = ?", vendorId)

	allowedSortBy := map[string]string{
		"rating": "rating",
		"date":   "created_at",
	}

	orderField := "created_at"

	if value, ok := allowedSortBy[params.SortBy]; ok {
		orderField = value
	}

	orderDir := "asc"

	if params.SortOrder == "desc" {
		orderDir = "desc"
	}

	db = db.Order(orderField + " " + orderDir)

	if err := db.Limit(params.Limit).Offset(params.Offset).Find(&reviews).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error getting paginated reviews data")
	}

	return reviews, nil

}

func (grr *GormReviewRepository) FindById(ctx context.Context, id uuid.UUID, vendorId uuid.UUID) (*models.Review, error) {

	var review models.Review

	if err := grr.db.WithContext(ctx).Preload("Replies").First(&review, "id = ? AND vendor_id = ?", id, vendorId).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error getting review data")
	}

	return &review, nil
}

func (grr *GormReviewRepository) Create(ctx context.Context, newReply *models.Reply) (*models.Reply, error) {

	if err := grr.db.WithContext(ctx).Create(newReply).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error creating new reply")
	}

	return newReply, nil
}

func (grr *GormReviewRepository) Update(ctx context.Context, updatedReply *models.Reply) error {

	res := grr.db.WithContext(ctx).Save(updatedReply)

	if res.Error != nil {
		return error_handler.ErrorHandler(res.Error, "Error updating reply")
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}

func (grr *GormReviewRepository) WithTx(tx *gorm.DB) domain.ReviewRepository {
	return &GormReviewRepository{
		db: tx,
	}
}

func (grr *GormReviewRepository) Transaction(fn func(txRepo domain.ReviewRepository) error) error {
	tx := grr.db.Begin()
	if tx.Error != nil {
		log.Printf("Transaction begin error: %v", tx.Error)
		return tx.Error
	}

	repo := grr.WithTx(tx)

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
