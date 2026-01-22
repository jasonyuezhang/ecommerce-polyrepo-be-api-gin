package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ecommerce/be-api-gin/internal/config"
	"github.com/ecommerce/be-api-gin/internal/models"
)

// Common errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
)

// Clients holds all gRPC client connections
type Clients struct {
	userConn      *grpc.ClientConn
	listingConn   *grpc.ClientConn
	inventoryConn *grpc.ClientConn
	config        *config.Config
}

// NewClients creates and initializes all gRPC client connections
// Connections are established asynchronously and will automatically retry if services are unavailable
func NewClients(cfg *config.Config) (*Clients, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Removed WithBlock() to allow async connection establishment and automatic retries
	}

	// Connect to User Service (non-blocking)
	userConn, err := grpc.Dial(cfg.UserServiceAddr, opts...)
	if err != nil {
		return nil, err
	}

	// Connect to Listing Service (non-blocking)
	listingConn, err := grpc.Dial(cfg.ListingServiceAddr, opts...)
	if err != nil {
		userConn.Close()
		return nil, err
	}

	// Connect to Inventory Service (non-blocking)
	inventoryConn, err := grpc.Dial(cfg.InventoryServiceAddr, opts...)
	if err != nil {
		userConn.Close()
		listingConn.Close()
		return nil, err
	}

	log.Printf("gRPC clients initialized for: user=%s, listing=%s, inventory=%s",
		cfg.UserServiceAddr, cfg.ListingServiceAddr, cfg.InventoryServiceAddr)

	return &Clients{
		userConn:      userConn,
		listingConn:   listingConn,
		inventoryConn: inventoryConn,
		config:        cfg,
	}, nil
}

// Close closes all gRPC connections
func (c *Clients) Close() {
	if c.userConn != nil {
		c.userConn.Close()
	}
	if c.listingConn != nil {
		c.listingConn.Close()
	}
	if c.inventoryConn != nil {
		c.inventoryConn.Close()
	}
}

// HealthCheck checks the health of all connected services
// A connection is considered healthy if it's in READY or IDLE state
// IDLE means the connection hasn't been used yet but can establish connection when needed
func (c *Clients) HealthCheck(ctx context.Context) map[string]bool {
	isHealthy := func(conn *grpc.ClientConn) bool {
		if conn == nil {
			return false
		}
		state := conn.GetState().String()
		// READY: Connection is established and ready
		// IDLE: Connection hasn't been used yet but will connect when needed
		return state == "READY" || state == "IDLE"
	}

	return map[string]bool{
		"user-service":      isHealthy(c.userConn),
		"listing-service":   isHealthy(c.listingConn),
		"inventory-service": isHealthy(c.inventoryConn),
	}
}

// handleGRPCError converts gRPC errors to application errors
func handleGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.PermissionDenied, codes.Unauthenticated:
		return ErrUnauthorized
	default:
		return ErrInternal
	}
}

// --- Listing Service Methods ---

// ListProducts fetches products from the listing service
func (c *Clients) ListProducts(ctx context.Context, page, limit int, category, search string) ([]*models.Product, int64, error) {
	// TODO: Implement actual gRPC call when proto files are available
	// For now, return mock data for development
	products := []*models.Product{
		{
			ID:          "prod-001",
			Name:        "Sample Product",
			Description: "A sample product for testing",
			Price:       29.99,
			Category:    "electronics",
			Available:   true,
		},
	}
	return products, 1, nil
}

// GetProduct fetches a single product from the listing service
func (c *Clients) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	// TODO: Implement actual gRPC call
	if id == "not-found" {
		return nil, ErrNotFound
	}
	return &models.Product{
		ID:          id,
		Name:        "Sample Product",
		Description: "A sample product for testing",
		Price:       29.99,
		Category:    "electronics",
		Available:   true,
	}, nil
}

// CreateProduct creates a new product via the listing service
func (c *Clients) CreateProduct(ctx context.Context, req *models.CreateProductRequest, userID string) (*models.Product, error) {
	// TODO: Implement actual gRPC call
	return &models.Product{
		ID:          "prod-new",
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Images:      req.Images,
		SellerID:    userID,
		Available:   true,
	}, nil
}

// UpdateProduct updates an existing product
func (c *Clients) UpdateProduct(ctx context.Context, id string, req *models.UpdateProductRequest, userID string) (*models.Product, error) {
	// TODO: Implement actual gRPC call
	return &models.Product{
		ID:       id,
		SellerID: userID,
	}, nil
}

// DeleteProduct deletes a product
func (c *Clients) DeleteProduct(ctx context.Context, id, userID string) error {
	// TODO: Implement actual gRPC call
	return nil
}

// --- Inventory Service Methods ---

// GetInventory gets inventory for a product
func (c *Clients) GetInventory(ctx context.Context, productID string) (*models.Inventory, error) {
	// TODO: Implement actual gRPC call
	return &models.Inventory{
		ProductID: productID,
		Quantity:  100,
		Reserved:  5,
		Available: true,
	}, nil
}

// InitializeInventory sets up initial inventory for a new product
func (c *Clients) InitializeInventory(ctx context.Context, productID string, quantity int32) error {
	// TODO: Implement actual gRPC call
	return nil
}

// UpdateInventory updates inventory quantity
func (c *Clients) UpdateInventory(ctx context.Context, productID string, quantity int32, operation string) (*models.Inventory, error) {
	// TODO: Implement actual gRPC call
	return &models.Inventory{
		ProductID: productID,
		Quantity:  quantity,
		Available: quantity > 0,
	}, nil
}

// CheckInventory checks if requested quantity is available
func (c *Clients) CheckInventory(ctx context.Context, productID string, quantity int32) (bool, error) {
	// TODO: Implement actual gRPC call
	return true, nil
}

// ReserveInventory reserves inventory for an order
func (c *Clients) ReserveInventory(ctx context.Context, productID string, quantity int32) (string, error) {
	// TODO: Implement actual gRPC call
	return "reservation-" + productID, nil
}

// CancelReservation cancels an inventory reservation
func (c *Clients) CancelReservation(ctx context.Context, reservationID string) error {
	// TODO: Implement actual gRPC call
	return nil
}

// --- User/Order Service Methods ---

// ListOrders fetches orders for a user
func (c *Clients) ListOrders(ctx context.Context, userID string, page, limit int, status string) ([]*models.Order, int64, error) {
	// TODO: Implement actual gRPC call
	return []*models.Order{}, 0, nil
}

// GetOrder fetches a single order
func (c *Clients) GetOrder(ctx context.Context, orderID, userID string) (*models.Order, error) {
	// TODO: Implement actual gRPC call
	if orderID == "not-found" {
		return nil, ErrNotFound
	}
	return &models.Order{
		ID:     orderID,
		UserID: userID,
		Status: "pending",
	}, nil
}

// CreateOrder creates a new order
func (c *Clients) CreateOrder(ctx context.Context, userID string, req *models.CreateOrderRequest, reservationIDs []string) (*models.Order, error) {
	// TODO: Implement actual gRPC call
	var items []models.OrderItem
	var total float64
	for _, item := range req.Items {
		orderItem := models.OrderItem{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  29.99, // Would come from product lookup
			TotalPrice: float64(item.Quantity) * 29.99,
		}
		items = append(items, orderItem)
		total += orderItem.TotalPrice
	}

	return &models.Order{
		ID:             "order-new",
		UserID:         userID,
		Items:          items,
		Status:         "pending",
		TotalAmount:    total,
		ShippingAddr:   req.ShippingAddr,
		ReservationIDs: reservationIDs,
	}, nil
}

// UpdateOrderStatus updates the status of an order
func (c *Clients) UpdateOrderStatus(ctx context.Context, orderID, userID, status string) (*models.Order, error) {
	// TODO: Implement actual gRPC call
	return &models.Order{
		ID:     orderID,
		UserID: userID,
		Status: status,
	}, nil
}

// CancelOrder cancels an order
func (c *Clients) CancelOrder(ctx context.Context, orderID, userID string) error {
	// TODO: Implement actual gRPC call
	return nil
}
