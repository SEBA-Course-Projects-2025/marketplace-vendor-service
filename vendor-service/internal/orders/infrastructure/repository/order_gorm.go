package repository

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"marketplace-vendor-service/vendor-service/internal/orders/domain"
	"marketplace-vendor-service/vendor-service/internal/orders/domain/models"
	"marketplace-vendor-service/vendor-service/internal/orders/dtos"
	"marketplace-vendor-service/vendor-service/internal/shared/tracer"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (gor *GormOrderRepository) FindAll(ctx context.Context, params dtos.OrderQueryParams, vendorId uuid.UUID) ([]models.Order, error) {

	var orders []models.Order

	db := gor.db.WithContext(ctx).Where("vendor_id = ?", vendorId).Preload("OrderItems")

	allowedSortBy := map[string]string{
		"totalPrice": "total_price",
		"date":       "created_at",
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

	if err := db.Limit(params.Limit).Offset(params.Offset).Find(&orders).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error getting orders data")
	}

	return orders, nil

}

func (gor *GormOrderRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Order, error) {

	ctx, span := tracer.Tracer.Start(ctx, "FindById")
	defer span.End()

	var order models.Order

	if err := gor.db.WithContext(ctx).Preload("OrderItems").First(&order, "id = ?", id).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error getting order data")
	}

	return &order, nil
}

func (gor *GormOrderRepository) Patch(ctx context.Context, updatedOrder *models.Order) (*models.Order, error) {

	ctx, span := tracer.Tracer.Start(ctx, "Patch")
	defer span.End()

	res := gor.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", updatedOrder.Id).Updates(updatedOrder)

	if res.Error != nil {
		return nil, error_handler.ErrorHandler(res.Error, "Error updating order status")
	}

	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return updatedOrder, nil
}

func (gor *GormOrderRepository) Create(ctx context.Context, newOrder *models.Order, vendorId uuid.UUID) (*models.Order, error) {

	ctx, span := tracer.Tracer.Start(ctx, "Create")
	defer span.End()

	newOrder.VendorId = vendorId

	if err := gor.db.WithContext(ctx).Create(newOrder).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error creating order")
	}

	return newOrder, nil
}

func (gor *GormOrderRepository) WithTx(tx *gorm.DB) domain.OrderRepository {
	return &GormOrderRepository{
		db: tx,
	}
}

func (gor *GormOrderRepository) Transaction(fn func(txRepo domain.OrderRepository) error) error {

	tx := gor.db.Begin()
	if tx.Error != nil {
		log.Printf("Transaction begin error: %v", tx.Error)
		return tx.Error
	}

	repo := gor.WithTx(tx)

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
