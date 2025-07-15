package dtos

import (
	"encoding/json"
	"github.com/google/uuid"
	eventModels "marketplace-vendor-service/vendor-service/internal/event/domain/models"
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
	ProductId   uuid.UUID `json:"product_id"`
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

type OrderStatusEventDto struct {
	OrderId    uuid.UUID `json:"order_id"`
	CustomerId uuid.UUID `json:"customer_id"`
	Status     string    `json:"status"`
}

func OrderStatusToEventDto(order *models.Order) OrderStatusEventDto {
	return OrderStatusEventDto{
		OrderId:    order.Id,
		CustomerId: order.CustomerId,
		Status:     order.Status,
	}
}

type OrderUpdatedStatusEvent struct {
	EventId uuid.UUID           `json:"event_id"`
	Order   OrderStatusEventDto `json:"order"`
}

func OrderStatusToOutbox(order *models.Order, eventType, exchange string) (*eventModels.Outbox, error) {

	event := OrderUpdatedStatusEvent{
		EventId: uuid.New(),
		Order:   OrderStatusToEventDto(order),
	}

	payload, err := json.Marshal(event)

	if err != nil {
		return nil, err
	}

	return &eventModels.Outbox{
		Id:          uuid.New(),
		EventType:   eventType,
		Exchange:    exchange,
		Payload:     payload,
		CreatedAt:   time.Now(),
		Processed:   false,
		ProcessedAt: time.Time{},
	}, nil

}

type OrderCreatedEventDto struct {
	EventId    uuid.UUID      `json:"event_id"`
	OrderId    uuid.UUID      `json:"order_id"`
	CustomerId uuid.UUID      `json:"customer_id"`
	Items      []OrderItemDto `json:"items"`
	TotalPrice float64        `json:"total_price"`
}

type OrderItemDto struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

func OrderItemsToOutbox(newOrderDto OrderCreatedEventDto, eventType, exchange string) (*eventModels.Outbox, error) {

	payload, err := json.Marshal(newOrderDto)

	if err != nil {
		return nil, err
	}

	return &eventModels.Outbox{
		Id:          uuid.New(),
		EventType:   eventType,
		Exchange:    exchange,
		Payload:     payload,
		CreatedAt:   time.Now(),
		Processed:   false,
		ProcessedAt: time.Time{},
	}, nil

}

type OrderCreatedEventDtoFromProduct struct {
	EventId    uuid.UUID      `json:"event_id"`
	OrderId    uuid.UUID      `json:"order_id"`
	CustomerId uuid.UUID      `json:"customer_id"`
	Items      []OrderItemDto `json:"items"`
	TotalPrice float64        `json:"total_price"`
	Status     string         `json:"status"`
}

type PostOrderDtoReq struct {
	CustomerId uuid.UUID      `json:"customer_id"`
	Items      []OrderItemDto `json:"items"`
}

type OrderCreatedEventResponseDto struct {
	EventId    uuid.UUID                   `json:"event_id"`
	OrderId    uuid.UUID                   `json:"order_id"`
	CustomerId uuid.UUID                   `json:"customer_id"`
	Items      []OrderProductEventResponse `json:"items"`
	TotalPrice float64                     `json:"total_price"`
	Status     string                      `json:"status"`
}

type OrderProductEventResponse struct {
	VendorId    uuid.UUID `json:"vendor_id"`
	ProductId   uuid.UUID `json:"product_id"`
	StockId     uuid.UUID `json:"stock_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	ImageUrl    string    `json:"image_url"`
	UnitPrice   float64   `json:"unit_price"`
}

func CreateOrderWithDto(orderReq OrderCreatedEventResponseDto) *models.Order {

	var newOrderItems []models.OrderItem

	for _, item := range orderReq.Items {

		newOrderItem := models.OrderItem{
			Id:          uuid.New(),
			ProductId:   item.ProductId,
			OrderId:     orderReq.OrderId,
			StockId:     item.StockId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			ImageUrl:    item.ImageUrl,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		newOrderItems = append(newOrderItems, newOrderItem)

	}

	return &models.Order{
		Id:         orderReq.OrderId,
		VendorId:   uuid.Nil,
		CustomerId: orderReq.CustomerId,
		TotalPrice: orderReq.TotalPrice,
		Status:     orderReq.Status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		OrderItems: newOrderItems,
	}
}

type CanceledOrderEventDto struct {
	EventId uuid.UUID `json:"event_id"`
	OrderId uuid.UUID `json:"order_id"`
	Status  string    `json:"status"`
}

type CanceledOrderItemDto struct {
	EventId   uuid.UUID `json:"event_id"`
	ProductId uuid.UUID `json:"product_id"`
	StockId   uuid.UUID `json:"stock_id"`
	Quantity  int       `json:"quantity"`
}

func CanceledOrderItemsToDto(order *models.Order) []CanceledOrderItemDto {

	var canceledOrderProducts []CanceledOrderItemDto

	for _, item := range order.OrderItems {
		canceledItem := CanceledOrderItemDto{
			EventId:   uuid.New(),
			ProductId: item.ProductId,
			StockId:   item.StockId,
			Quantity:  item.Quantity,
		}

		canceledOrderProducts = append(canceledOrderProducts, canceledItem)
	}

	return canceledOrderProducts
}

func CanceledOrderProductsToOutbox(canceledOrderProducts []CanceledOrderItemDto, eventType, exchange string) (*eventModels.Outbox, error) {

	payload, err := json.Marshal(canceledOrderProducts)

	if err != nil {
		return nil, err
	}

	return &eventModels.Outbox{
		Id:          uuid.New(),
		EventType:   eventType,
		Exchange:    exchange,
		Payload:     payload,
		CreatedAt:   time.Now(),
		Processed:   false,
		ProcessedAt: time.Time{},
	}, nil

}

var AllowedStatuses = map[string]struct{}{
	"pending":   {},
	"confirmed": {},
	"shipped":   {},
	"delivered": {},
	"completed": {},
	"cancelled": {},
	"declined":  {},
}
