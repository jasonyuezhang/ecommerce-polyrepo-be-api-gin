package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ecommerce/be-api-gin/internal/models"
	grpcclient "github.com/ecommerce/be-api-gin/pkg/grpc"
)

// ProductHandler handles product-related requests
type ProductHandler struct {
	grpcClients *grpcclient.Clients
}

// NewProductHandler creates a new product handler
func NewProductHandler(clients *grpcclient.Clients) *ProductHandler {
	return &ProductHandler{
		grpcClients: clients,
	}
}

// ListProducts returns a list of all products
// GET /api/v1/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	search := c.Query("search")

	// Call listing service via gRPC
	products, total, err := h.grpcClients.ListProducts(c.Request.Context(), page, limit, category, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch products",
			Message: err.Error(),
		})
		return
	}

	// Set InStock field for frontend compatibility
	for i := range products {
		products[i].InStock = products[i].Available
		// Set ImageUrl from first image if available
		if len(products[i].Images) > 0 {
			products[i].ImageUrl = products[i].Images[0]
		}
	}

	c.JSON(http.StatusOK, models.ProductsResponse{
		Products: products,
		Page:     page,
		Limit:    limit,
		Total:    total,
	})
}

// GetProduct returns a single product by ID
// GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	// Call listing service via gRPC
	product, err := h.grpcClients.GetProduct(c.Request.Context(), id)
	if err != nil {
		if err == grpcclient.ErrNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Product not found",
				Message: "No product exists with the given ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch product",
			Message: err.Error(),
		})
		return
	}

	// Get inventory info
	inventory, err := h.grpcClients.GetInventory(c.Request.Context(), id)
	if err == nil {
		product.Stock = inventory.Quantity
		product.Available = inventory.Available
	}

	// Set InStock field for frontend compatibility
	product.InStock = product.Available
	// Set ImageUrl from first image if available
	if len(product.Images) > 0 {
		product.ImageUrl = product.Images[0]
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product
// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, _ := c.Get("userID")

	// Call listing service via gRPC
	product, err := h.grpcClients.CreateProduct(c.Request.Context(), &req, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create product",
			Message: err.Error(),
		})
		return
	}

	// Initialize inventory for the product
	if err := h.grpcClients.InitializeInventory(c.Request.Context(), product.ID, req.InitialStock); err != nil {
		// Log error but don't fail the request
		// Inventory can be updated later
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product
// PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("userID")

	// Call listing service via gRPC
	product, err := h.grpcClients.UpdateProduct(c.Request.Context(), id, &req, userID.(string))
	if err != nil {
		if err == grpcclient.ErrNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Product not found",
				Message: "No product exists with the given ID",
			})
			return
		}
		if err == grpcclient.ErrUnauthorized {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have permission to update this product",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update product",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product
// DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID, _ := c.Get("userID")

	// Call listing service via gRPC
	err := h.grpcClients.DeleteProduct(c.Request.Context(), id, userID.(string))
	if err != nil {
		if err == grpcclient.ErrNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Product not found",
				Message: "No product exists with the given ID",
			})
			return
		}
		if err == grpcclient.ErrUnauthorized {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have permission to delete this product",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete product",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// UpdateInventory updates product inventory
// PUT /api/v1/products/:id/inventory
func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Call inventory service via gRPC
	inventory, err := h.grpcClients.UpdateInventory(c.Request.Context(), id, req.Quantity, req.Operation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update inventory",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, inventory)
}
