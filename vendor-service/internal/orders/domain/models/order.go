package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	Id         uuid.UUID   `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	VendorId   uuid.UUID   `json:"vendor_id" gorm:"column:vendor_id;type:uuid;not null"`
	CustomerId uuid.UUID   `json:"customer_id" gorm:"column:customer_id;type:uuid;not null"`
	TotalPrice float64     `json:"total_price" gorm:"column:total_price;type:numeric(12, 2);not null"`
	Status     string      `json:"status" gorm:"column:status;type:varchar(40);not null"`
	CreatedAt  time.Time   `gorm:"column:created_at;type:timestamp"`
	UpdatedAt  time.Time   `gorm:"column:updated_at;type:timestamp"`
	OrderItems []OrderItem `json:"items" gorm:"foreignKey:OrderId;references:Id"`
}

type OrderItem struct {
	Id          uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	ProductId   uuid.UUID `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	OrderId     uuid.UUID `json:"order_id" gorm:"column:order_id;type:uuid;not null"`
	ProductName string    `json:"product_name" gorm:"column:product_name;type:varchar(255);not null"`
	Quantity    int       `json:"quantity" gorm:"column:quantity;type:int;not null"`
	UnitPrice   float64   `json:"unit_price" gorm:"column:unit_price;type:numeric(12, 2);not null"`
	ImageUrl    string    `json:"image_url" gorm:"column:image_url;type:text;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp"`
}
