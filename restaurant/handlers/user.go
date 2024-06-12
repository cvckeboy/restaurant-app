package handlers

import (
	"github.com/cvckeboy/restaurant-app/restaurant/models"
	"github.com/cvckeboy/restaurant-app/restaurant/services"
	"github.com/cvckeboy/restaurant-app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
	logger  *utils.Logger
}

func NewUserHandler(service *services.UserService, logger *utils.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

func (h *UserHandler) Register(router *gin.Engine) {
	router.POST("/register", h.RegisterUser)
	router.POST("/login", h.LoginUser)
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req models.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.service.AuthenticateUser(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to authenticate user", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		h.logger.Error("Failed to generate JWT", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
