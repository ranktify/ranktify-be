package jwt

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ranktify/ranktify-be/internal/model"
)

type AccessTokenClaims struct {
	UserID   uint64
	Username string
	Email    string
}

func getAccessKey() []byte {
	accesskey := os.Getenv("JWT_ACCESS_KEY")
	if accesskey == "" {
		log.Fatalf("No access key found")
	}
	return []byte(accesskey)
}

func getRefreshKey() []byte {
	refreshKey := os.Getenv("JWT_REFRESH_KEY")
	if refreshKey == "" {
		log.Fatalf("No refresh key found")
	}

	return []byte(refreshKey)
}

func getIssuerString() string {
	return os.Getenv("JWT_ISSUER")
}

type customClaims struct {
	jwt.RegisteredClaims
	UserID   float64 `json:"user_id"`
	Username string  `json:"username,omitempty"` // missing in the rt but not in the at
	Email    string  `json:"email,omitempty"`
}

// validate a token using the provided secret key, assumes that the signing method was HS256.
// If valid returns a *jwt.Token and nil, otherwise nil and error.
func validateToken(tokenString string, tokenSecretKey []byte, issuer string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return tokenSecretKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return token, nil
}

// Creates the access token string, claims user_id, username, email. Expires in 15 minutes.
func createAccessToken(user model.User) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"user_id":    user.Id,
		"username":   user.Username,
		"email":      user.Email,
		"iss":        getIssuerString(),
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(15 * time.Minute).Unix(), // Expiry date: 15 minutes
		"jti":        uuid.New().String(),
	})

	accessTokenString, err := accessToken.SignedString(getAccessKey())
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}

// Creates the refresh token string, claims the user.Id, and expires in a month.
func createRefreshToken(user model.User) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"iss":     getIssuerString(),
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(24 * time.Hour * 30).Unix(), // Expires in a month
		"jti":     uuid.New().String(),                        // Unique identifier for the token
	})

	refreshTokenString, err := refreshToken.SignedString(getRefreshKey())
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}

// CreateTokens returns both access and refresh tokens
// claiming atributes from the user provided.
func CreateTokens(user model.User) (string, string) {
	accessTokenString, err := createAccessToken(user)
	if err != nil {
		panic(err) // TODO: change this and handle accordingly
	}
	refreshTokenString, err := createRefreshToken(user)
	if err != nil {
		panic(err)
	}
	return accessTokenString, refreshTokenString
}

// validates the access token using the access key env.
func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	accessKey := getAccessKey()
	iss := getIssuerString()
	token, err := validateToken(tokenString, accessKey, iss)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Validates the given refresh token and generates a new access token
// and new refresh token to be rotated in storage.
func RefreshTokens(refreshTokenString string, user model.User) (string, string, error) {
	refreshKey := getRefreshKey()
	iss := getIssuerString()
	_, err := validateToken(refreshTokenString, refreshKey, iss)
	if err != nil {
		return "", "", fmt.Errorf("refresh token failed to validate")
	}
	// Create new access and refresh tokens.
	accessTokenString, err := createAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("error creating jwt access token")
	}
	newRefreshTokenString, err := createRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("error creating jwt refresh token")
	}
	return accessTokenString, newRefreshTokenString, nil
}

// Validates and parses refresh token claims. Returns a model.JWTRefreshToken on success,
// with the following members filled:
// - UserID
// - JTI
// - ExpiresAt
// - RefreshToken
func ParseRefreshTokenClaims(tokenString string) (*model.JWTRefreshToken, error) {
	rk := getRefreshKey()
	iss := getIssuerString()
	token, err := validateToken(tokenString, rk, iss)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return nil, fmt.Errorf("error extracting claims from token")
	}
	rt := &model.JWTRefreshToken{
		UserID:       uint64(claims.UserID),
		JTI:          claims.ID,
		ExpiresAt:    claims.ExpiresAt.Time,
		RefreshToken: tokenString,
	}

	return rt, nil
}

func GetClaimsFromAccessToken(token *jwt.Token) (*AccessTokenClaims, error) {
	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return nil, fmt.Errorf("error extracting claims from token 1")
	}

	return &AccessTokenClaims{
		UserID:   uint64(claims.UserID),
		Username: claims.Username,
		Email:    claims.Email,
	}, nil
}
