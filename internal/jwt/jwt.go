package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ranktify/ranktify-be/internal/model"
)

func getSecretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func CreateToken(user model.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"user_id":    user.Id,
		"username":   user.Username,
		"email":      user.Email,
		"iss":        time.Now().Unix(),                    // Issued date
		"exp":        time.Now().Add(time.Hour * 2).Unix(), // Expiry date: two hours
	})

	secretKey := getSecretKey()
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

// Validates the token signature to see if it matches with
// JWT_SECRET_KEY and the claims. handle by the jwt.Parse function
func VerifyToken(tokenString string) (*jwt.Token, error) {
	secretKey := getSecretKey()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}
