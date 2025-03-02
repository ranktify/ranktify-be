package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DAO *dao.UserDAO
}

func NewUserHandler(dao *dao.UserDAO) *UserHandler {
	return &UserHandler{DAO: dao}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbUser, err := h.DAO.GetUser(user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// if dbUser is not nil, then we found a user that already has the same username or email
	if dbUser != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User already exist"})
		return
	}
	// TODO: missing atributes validation
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 11)
	user.Password = string(bytes)

	if err = h.DAO.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": fmt.Sprintf("Created user with id: %d", user.Id)})
}

func (h *UserHandler) ValidateUser(c *gin.Context) {
	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbUser, err := h.DAO.GetUser(user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if dbUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Succesfully authenticated"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := h.DAO.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.DAO.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users"})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No users found"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	err = h.DAO.UpdateUserByID(userID, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	err = h.DAO.DeleteUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to Delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "User Deleted Succesfully"})
}
