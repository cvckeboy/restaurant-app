package handlers

import (
	"github.com/cvckeboy/restaurant-app/restaurant/models"
	"github.com/cvckeboy/restaurant-app/restaurant/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type ProductHandler struct {
	service *services.ProductService
	logger  *slog.Logger
}

func NewProductHandler(service *services.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{service: service, logger: logger}
}

func (h *ProductHandler) Register(router *gin.Engine) {
	router.GET("/products", h.GetAllProducts)
	router.POST("/products", h.CreateProduct)
	router.PUT("/products/:id", h.UpdateProduct)
	router.DELETE("/products/:id", h.DeleteProduct)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get all products", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to create product", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.service.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("ERROR: %v", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr) // Используем uuid.Parse для преобразования строки в UUID
	if err != nil {
		h.logger.Error("failed to update product", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to update product", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ID = id // Присваиваем ID из параметра запроса

	if err := h.service.UpdateProduct(c.Request.Context(), &req); err != nil {
		h.logger.Error("failed to update product", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr) // Используем uuid.Parse для преобразования строки в UUID
	if err != nil {
		h.logger.Error("failed to delete product", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		h.logger.Error("deleting product", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
