package service

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/jwt"
	"github.com/ranktify/ranktify-be/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserDAO   *dao.UserDAO
	TokensDAO *dao.TokensDAO
}

// replaces gin.H, hence decoupling web framework from service layer
type content map[string]any

func NewUserService(userDAO *dao.UserDAO, tokensDAO *dao.TokensDAO) *UserService {
	return &UserService{
		UserDAO:   userDAO,
		TokensDAO: tokensDAO,
	}
}

func (s *UserService) CreateUser(user *model.User) (int, content) {
	dbUser, err := s.UserDAO.GetUser(user.Email, user.Username)
	if err != nil {
		return http.StatusInternalServerError, content{"error": err.Error()}
	}
	// if dbUser is not nil, then we found a user that already has the same username or email
	if dbUser != nil {
		return http.StatusConflict, content{"error": "User already exist"}
	}
	// TODO: missing atributes validation
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 11)
	if err != nil {
		return http.StatusInternalServerError, content{"error": err.Error()}
	}
	user.Password = string(bytes)

	if err = s.UserDAO.CreateUser(user); err != nil {
		return http.StatusInternalServerError, content{"error": err.Error()}

	}

	accessToken, refreshToken := jwt.CreateTokens(*user)
	// TODO: Store the rt in storage or rotate
	return http.StatusCreated, content{
		"success":       fmt.Sprintf("Created user with id: %d", user.Id),
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
}

func (s *UserService) ValidateUser(user *model.User) (int, content) {
	dbUser, err := s.UserDAO.GetUser(user.Email, user.Username)
	if err != nil {
		return http.StatusInternalServerError, content{"error": err.Error()}

	}

	if dbUser == nil {
		return http.StatusNotFound, content{"error": "User does not exist"}

	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return http.StatusUnauthorized, content{"error": "Password is incorrect"}

	}

	accessToken, refreshToken := jwt.CreateTokens(*dbUser)
	return http.StatusOK, content{
		"success":       fmt.Sprintf("Logged in user with id: %d", dbUser.Id),
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
}

func (s *UserService) GetUserByID(userID uint64) (int, content) {
	user, err := s.UserDAO.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, content{"error": fmt.Sprintf("User with id %d not found", userID)}
		}

		return http.StatusInternalServerError, content{"error": "Failed to retrieve user"}
	}

	return http.StatusOK, content{"user": user}
}

func (s *UserService) GetAllUsers() (int, content) {
	users, err := s.UserDAO.GetAllUsers()
	if err != nil {
		return http.StatusInternalServerError, content{"error": "Failed to retrieve users"}
	}
	if len(users) == 0 {
		return http.StatusNotFound, content{"message": "No users found"}
	}

	return http.StatusOK, content{"users": users}
}

func (s *UserService) UpdateUserByID(userID uint64, user *model.User) (int, content) {
	err := s.UserDAO.UpdateUserByID(userID, user)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, content{"error": fmt.Sprintf("User with id %d not found", userID)}
		}

		return http.StatusInternalServerError, content{"error": "Failed to update user"}
	}

	return http.StatusOK, content{"message": "User updated successfully"}
}

func (s *UserService) DeleteUserByID(userID uint64) (int, content) {
	err := s.UserDAO.DeleteUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, content{"error": fmt.Sprintf("User with id %d not found", userID)}
		}

		return http.StatusInternalServerError, content{"error": "Failed to delete user"}
	}

	return http.StatusOK, content{"message": "User deleted successfully"}
}
