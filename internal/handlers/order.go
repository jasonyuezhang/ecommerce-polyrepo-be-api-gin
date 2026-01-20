package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ecommerce/be-api-gin/internal/models"
	grpcclient "github.com/ecommerce/be-api-gin/pkg/grpc"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	grpcClients *grpcclient.Clients
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(clients *grpcclient.Clients) *OrderHandler {
	return &OrderHandler{
		grpcClients: clients,
	}
}

// ListOrders returns a list of orders for the authenticated user
// GET /api/v1/orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, _ := c.Get("userID")

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	// Call user service via gRPC to get orders
	orders, total, err := h.grpcClients.ListOrders(c.Request.Context(), userID.(string), page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch orders",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       orders,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	})
}

// GetOrder returns a single order by ID
// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	// Call user service via gRPC
	order, err := h.grpcClients.GetOrder(c.Request.Context(), id, userID.(string))
	if err != nil {
		if err == grpcclient.ErrNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Order not found",
				Message: "No order exists with the given ID",
			})
			return
		}
		if err == grpcclient.ErrUnauthorized {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have permission to view this order",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch order",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateOrder creates a new order
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	userID, _ := c.Get("userID")

	// Validate inventory availability for all items
	for _, item := range req.Items {
		available, err := h.grpcClients.CheckInventory(c.Request.Context(), item.ProductID, item.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to check inventory",
				Message: err.Error(),
			})
			return
		}
		if !available {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Insufficient inventory",
				Message: "Product " + item.ProductID + " does not have enough stock",
			})
			return
		}
	}

	// Reserve inventory for all items
	reservationIDs := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		reservationID, err := h.grpcClients.ReserveInventory(c.Request.Context(), item.ProductID, item.Quantity)
		if err != nil {
			// Rollback previous reservations
			for _, rid := range reservationIDs {
				h.grpcClients.CancelReservation(c.Request.Context(), rid)
			}
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to reserve inventory",
				Message: err.Error(),
			})
			return
		}
		reservationIDs = append(reservationIDs, reservationID)
	}

	// Create the order
	order, err := h.grpcClients.CreateOrder(c.Request.Context(), userID.(string), &req, reservationIDs)
	if err != nil {
		// Rollback reservations on failure
		for _, rid := range reservationIDs {
			h.grpcClients.CancelReservation(c.Request.Context(), rid)
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create order",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// UpdateOrderStatus updates the status of an order
// PUT /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Call user service via gRPC
	order, err := h.grpcClients.UpdateOrderStatus(c.Request.Context(), id, userID.(string), req.Status)
	if err != nil {
		if err == grpcclient.ErrNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Order not found",
				Message: "No order exists with the given ID",
			})
			return
		}
		if err == grpcclient.ErrUnauthorized {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have permission to update this order",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update order status",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CancelOrder cancels an order
// DELETE /api/v1/orders/:id
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	// Get the order first to retrieve reservation IDs
	order, err := h.grpcClients.GetOrder(c.Request.Context(), id, userID.(string))
	if err != nil {
		if err == grpcclient.ErrNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Order not found",
				Message: "No order exists with the given ID",
			})
			return
		}
		if err == grpcclient.ErrUnauthorized {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have permission to cancel this order",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch order",
			Message: err.Error(),
		})
		return
	}

	// Check if order can be cancelled
	if order.Status != "pending" && order.Status != "confirmed" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Cannot cancel order",
			Message: "Order can only be cancelled when in pending or confirmed status",
		})
		return
	}

	// Cancel the order
	err = h.grpcClients.CancelOrder(c.Request.Context(), id, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to cancel order",
			Message: err.Error(),
		})
		return
	}

	// Release inventory reservations
	for _, reservationID := range order.ReservationIDs {
		h.grpcClients.CancelReservation(c.Request.Context(), reservationID)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Order cancelled successfully",
	})
}
