package models

import (
	"github.com/google/uuid"
	"time"
)

type Review struct {
	Id           uuid.UUID `json:"id" gorm:"column:id;type:uuid;not null;primaryKey"`
	VendorId     uuid.UUID `json:"vendor_id" gorm:"column:vendor_id;type:uuid;not null"`
	ProductId    uuid.UUID `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	ReviewerId   uuid.UUID `json:"reviewer_id" gorm:"column:reviewer_id;type:uuid;not null"`
	ReviewerName string    `json:"reviewer_name" gorm:"column:reviewer_name;type:varchar(255);not null"`
	Rating       float32   `json:"rating" gorm:"column:rating;type:float8;not null"`
	Comment      string    `json:"comment" gorm:"column:comment;type:text;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp"`
	Replies      []Reply   `gorm:"foreignKey:ReviewId;references:Id"`
}

type Reply struct {
	Id          uuid.UUID `json:"id" gorm:"column:id;type:uuid;not null;primaryKey"`
	ReviewId    uuid.UUID `json:"review_id" gorm:"column:review_id;type:uuid;not null"`
	ReplierId   uuid.UUID `json:"replier_id" gorm:"column:replier_id;type:uuid;not null"`
	ReplierName string    `json:"replier_name" gorm:"column:replier_name;type:varchar(255);not null"`
	Comment     string    `json:"comment" gorm:"column:comment;type:text;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp"`
}
