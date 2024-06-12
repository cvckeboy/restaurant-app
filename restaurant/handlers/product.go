package handlers

import (
	"github.com/cvckeboy/restaurant-app/middleware"
	"github.com/cvckeboy/restaurant-app/restaurant/models"
	"github.com/cvckeboy/restaurant-app/restaurant/services"
	"github.com/cvckeboy/restaurant-app/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ProductHandler struct {
	service *services.ProductService
	logger  *utils.Logger
}

func NewProductHandler(service *services.ProductService, logger *utils.Logger) *ProductHandler {
	return &ProductHandler{service: service, logger: logger}
}

func (h *ProductHandler) Register(router *gin.Engine) {
	router.GET("/products", h.GetAllProducts)
	router.GET("/products/:id", h.GetProductByID)
	router.GET("/products/category/:category_name", h.GetProductsByCategory)
	router.GET("/products/sorted", h.GetProductsSortedByPrice)

	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.POST("/products", middleware.JWTAuthMiddleware(), h.CreateProduct)
	protected.PUT("/products/:id", middleware.AdminMiddleware(), h.UpdateProduct)
	protected.DELETE("/products/:id", middleware.AdminMiddleware(), h.DeleteProduct)
	protected.POST("/categories", middleware.AdminMiddleware(), h.CreateCategory)
	protected.POST("/images", middleware.AdminMiddleware(), h.CreateImage)
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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid product ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.UpdateProduct(c.Request.Context(), id, &req); err != nil {
		h.logger.Error("Failed to update product", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product updated"})
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

func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.service.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create category", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ProductHandler) CreateImage(c *gin.Context) {
	var req models.CreateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.service.CreateImage(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create image", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create image"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid product ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var product models.Product
	if err := h.service.GetProductByID(c.Request.Context(), id, &product); err != nil {
		h.logger.Error("Failed to get product by ID", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryName := c.Param("category_name")
	products, err := h.service.GetProductsByCategory(c.Request.Context(), categoryName)
	if err != nil {
		h.logger.Error("Failed to get products by category", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
func (h *ProductHandler) GetProductsSortedByPrice(c *gin.Context) {

}
