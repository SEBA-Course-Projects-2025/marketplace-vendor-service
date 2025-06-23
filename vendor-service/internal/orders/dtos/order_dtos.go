package dtos

import (
	"github.com/google/uuid"
	"marketplace-vendor-service/vendor-service/internal/orders/domain/models"
	"time"
)

type OrderQueryParams struct {
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

type GetOrdersResponse struct {
	OrderId    uuid.UUID `json:"order_id"`
	CustomerId uuid.UUID `json:"customer_id"`
	VendorId   uuid.UUID `json:"vendor_id"`
	Items      []string  `json:"items"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	Date       time.Time `json:"date"`
}

func OrdersToDto(orders []models.Order) []GetOrdersResponse {

	var ordersResponse []GetOrdersResponse

	for _, order := range orders {

		var items []string
		var totalPrice float64

		for _, item := range order.OrderItems {
			items = append(items, item.ProductName)
			totalPrice += item.UnitPrice * float64(item.Quantity)
		}

		orderResponse := GetOrdersResponse{
			OrderId:    order.Id,
			CustomerId: order.CustomerId,
			VendorId:   order.VendorId,
			Items:      items,
			TotalPrice: totalPrice,
			Status:     order.Status,
			Date:       order.CreatedAt,
		}

		ordersResponse = append(ordersResponse, orderResponse)
	}

	return ordersResponse

}

type OrderItemResponse struct {
	ProductId   uuid.UUID `json:"productId"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	ImageUrl    string    `json:"image_url"`
	UnitPrice   float64   `json:"unit_price"`
}

type OneOrderResponse struct {
	OrderId    uuid.UUID           `json:"order_id"`
	CustomerId uuid.UUID           `json:"customer_id"`
	VendorId   uuid.UUID           `json:"vendor_id"`
	Items      []OrderItemResponse `json:"items"`
	TotalPrice float64             `json:"total_price"`
	Status     string              `json:"status"`
	Date       time.Time           `json:"date"`
}

func OrderToDto(order *models.Order) OneOrderResponse {

	var items []OrderItemResponse
	var totalPrice float64

	for _, item := range order.OrderItems {
		items = append(items, OrderItemResponse{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			ImageUrl:    item.ImageUrl,
			UnitPrice:   item.UnitPrice,
		})
		totalPrice += item.UnitPrice * float64(item.Quantity)
	}

	return OneOrderResponse{
		OrderId:    order.Id,
		CustomerId: order.CustomerId,
		VendorId:   order.VendorId,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     order.Status,
		Date:       order.CreatedAt,
	}

}

type StatusRequestDto struct {
	Status string `json:"status"`
}
