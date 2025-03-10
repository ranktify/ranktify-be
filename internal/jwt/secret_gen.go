package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Prints to the console a base64 string to be stored in the .env.
// This string is used as the JWT_SECRET_KEY in the .env
// Used in the main function for a single print and then not called again.
func GenerateNewJWTSecretKey() {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(bytes))
}
