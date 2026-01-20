package models

import "time"

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int64       `json:"total_pages"`
}

// ProductsResponse represents a paginated products response
type ProductsResponse struct {
	Products   []*Product `json:"products"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	Total      int64      `json:"total"`
}

// Product represents a product
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category,omitempty"`
	ImageUrl    string    `json:"imageUrl,omitempty"`
	Images      []string  `json:"images,omitempty"`
	SellerID    string    `json:"seller_id,omitempty"`
	Stock       int32     `json:"stock,omitempty"`
	InStock     bool      `json:"inStock"`
	Available   bool      `json:"available,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	Name         string   `json:"name" binding:"required,min=1,max=200"`
	Description  string   `json:"description" binding:"max=5000"`
	Price        float64  `json:"price" binding:"required,gt=0"`
	Category     string   `json:"category" binding:"required"`
	Images       []string `json:"images"`
	InitialStock int32    `json:"initial_stock" binding:"gte=0"`
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Name        *string   `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	Description *string   `json:"description,omitempty" binding:"omitempty,max=5000"`
	Price       *float64  `json:"price,omitempty" binding:"omitempty,gt=0"`
	Category    *string   `json:"category,omitempty"`
	Images      *[]string `json:"images,omitempty"`
}

// Inventory represents inventory information
type Inventory struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Reserved  int32  `json:"reserved"`
	Available bool   `json:"available"`
}

// UpdateInventoryRequest represents a request to update inventory
type UpdateInventoryRequest struct {
	Quantity  int32  `json:"quantity" binding:"required"`
	Operation string `json:"operation" binding:"required,oneof=set add subtract"`
}

// Order represents an order
type Order struct {
	ID             string      `json:"id"`
	UserID         string      `json:"user_id"`
	Items          []OrderItem `json:"items"`
	Status         string      `json:"status"`
	TotalAmount    float64     `json:"total_amount"`
	ShippingAddr   Address     `json:"shipping_address"`
	ReservationIDs []string    `json:"reservation_ids,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// Address represents a shipping or billing address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	Items        []CreateOrderItem `json:"items" binding:"required,min=1,dive"`
	ShippingAddr Address           `json:"shipping_address" binding:"required"`
}

// CreateOrderItem represents an item in a create order request
type CreateOrderItem struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required,gt=0"`
}

// UpdateOrderStatusRequest represents a request to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed processing shipped delivered cancelled"`
}

// User represents a user
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
