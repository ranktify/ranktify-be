package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/model"
	"github.com/ranktify/ranktify-be/internal/service"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusCode, content := h.Service.CreateUser(&user)

	c.JSON(statusCode, content)
}

func (h *UserHandler) ValidateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	statusCode, content := h.Service.ValidateUser(&user)

	c.JSON(statusCode, content)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	statusCode, content := h.Service.GetUserByID(userID)
	c.JSON(statusCode, content)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	statusCode, content := h.Service.GetAllUsers()

	c.JSON(statusCode, content)
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	statusCode, content := h.Service.UpdateUserByID(userID, &user)

	c.JSON(statusCode, content)
}

func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	statusCode, content := h.Service.DeleteUserByID(userID)
	c.JSON(statusCode, content)
}

func (h *UserHandler) SearchUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username can't be empty"})
		return
	}

	statusCode, users := h.Service.SearchUser(username)
	if statusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}
	if statusCode == http.StatusInternalServerError {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}
	c.JSON(statusCode, users)
}
