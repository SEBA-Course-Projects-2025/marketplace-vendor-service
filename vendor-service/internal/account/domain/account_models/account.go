package account_models

import (
	"github.com/google/uuid"
	"time"
)

type VendorAccount struct {
	Id           uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Email        string    `json:"email" gorm:"column:email;type:varchar(255);not null unique"`
	PasswordHash string    `json:"password_hash" gorm:"column:password_hash;type:varchar(255);not null"`
	Name         string    `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Description  string    `json:"description" gorm:"column:description;type:text;not null"`
	Logo         string    `json:"logo" gorm:"column:logo;type:varchar(255);not null"`
	Address      string    `json:"address" gorm:"column:address;type:varchar(255);not null"`
	Website      string    `json:"website" gorm:"column:website;type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp"`
}
