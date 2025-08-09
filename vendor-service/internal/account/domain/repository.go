package domain

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"marketplace-vendor-service/vendor-service/internal/account/domain/account_models"
)

type AccountRepository interface {
	FindById(ctx context.Context, vendorId uuid.UUID) (*account_models.VendorAccount, error)
	Update(ctx context.Context, updatedAccount *account_models.VendorAccount) error
	Patch(ctx context.Context, modifiedAccount *account_models.VendorAccount) (*account_models.VendorAccount, error)
	FindByEmail(ctx context.Context, email string) (*account_models.VendorAccount, error)
	Transaction(fn func(txRepo AccountRepository) error) error
	WithTx(tx *gorm.DB) AccountRepository
}
