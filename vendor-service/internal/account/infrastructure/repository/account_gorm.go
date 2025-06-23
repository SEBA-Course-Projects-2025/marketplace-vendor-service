package repository

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"marketplace-vendor-service/vendor-service/internal/account/domain"
	"marketplace-vendor-service/vendor-service/internal/account/domain/account_models"
	"marketplace-vendor-service/vendor-service/internal/shared/utils/error_handler"
)

type GormAccountRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GormAccountRepository {
	return &GormAccountRepository{db: db}
}

func (gar *GormAccountRepository) FindById(ctx context.Context, vendorId uuid.UUID) (*account_models.VendorAccount, error) {

	var account account_models.VendorAccount

	if err := gar.db.WithContext(ctx).First(&account, "id = ?", vendorId).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error getting vendor's account data")
	}

	return &account, nil
}

func (gar *GormAccountRepository) Update(ctx context.Context, updatedAccount *account_models.VendorAccount) error {

	if err := gar.db.WithContext(ctx).Save(updatedAccount).Error; err != nil {
		return error_handler.ErrorHandler(err, "Error updating account")
	}

	return nil

}

func (gar *GormAccountRepository) Patch(ctx context.Context, modifiedAccount *account_models.VendorAccount) (*account_models.VendorAccount, error) {

	if err := gar.db.WithContext(ctx).Save(modifiedAccount).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error modifying account")
	}

	return modifiedAccount, nil

}

func (gar *GormAccountRepository) FindByEmail(ctx context.Context, email string) (*account_models.VendorAccount, error) {

	var account account_models.VendorAccount

	if err := gar.db.WithContext(ctx).First(&account, "email = ?", email).Error; err != nil {
		return nil, error_handler.ErrorHandler(err, "Error getting vendor's account data")
	}

	return &account, nil

}

func (gar *GormAccountRepository) WithTx(tx *gorm.DB) domain.AccountRepository {
	return &GormAccountRepository{
		db: tx,
	}
}

func (gar *GormAccountRepository) Transaction(fn func(txRepo domain.AccountRepository) error) error {
	tx := gar.db.Begin()
	if tx.Error != nil {
		log.Printf("Transaction begin error: %v", tx.Error)
		return tx.Error
	}

	repo := gar.WithTx(tx)

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
