package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ranktify/ranktify-be/internal/model"
)

func getAccessKey() []byte {
	return []byte(os.Getenv("JWT_ACCESS_KEY"))
}

func getRefreshKey() []byte {
	return []byte(os.Getenv("JWT_REFRESH_KEY"))
}

// validate a token using the provided secret key, assumes that the signing method was HS256.
// If valid returns a *jwt.Token and nil, otherwise nil and error.
func validateToken(tokenString string, tokenSecretKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tokenSecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

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
		"iss":        time.Now().Unix(),                       // Issued date
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
		"iss":     time.Now().Unix(),
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
	token, err := validateToken(tokenString, accessKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Validates the given refresh token and generates a new access token
// and new refresh token to be rotated in storage.
func RefreshTokens(refreshTokenString string, user model.User) (string, string, error) {
	refreshKey := getRefreshKey()
	_, err := validateToken(refreshTokenString, refreshKey)
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

// Parses a *jwt.Token. Returns a jwt.MapClaims to be used in middleware on success.
func ParseTokenClaims(jwtToken *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error extracting claims from token")
	}

	return claims, nil
}

// Validates and parses refresh token claims. Returns a model.JWTRefreshToken on success,
// with the following members filled:
// - UserID
// - JTI
// - ExpiresAt
// - RefreshToken
func ParseRefreshTokenClaims(tokenString string) (*model.JWTRefreshToken, error) {
	token, err := validateToken(tokenString, getRefreshKey())
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error extracting claims from token")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, fmt.Errorf("error extracting jti from token")
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("error extracting user_id from token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("error extracting exp from token")
	}
	// hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokenString), 11)
	// if err != nil {
	// 	return nil, err
	// }

	rt := &model.JWTRefreshToken{
		UserID:       uint64(userId),
		JTI:          jti,
		ExpiresAt:    time.Unix(int64(exp), 0),
		RefreshToken: tokenString,
	}

	return rt, nil
}
