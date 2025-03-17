package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Prints to the console a two base64 string to be stored in the .env.
// This function will print two secret key: one for the JWT_ACCESS_KEY and
// the other one for JWT_REFRESH_KEY in the .env.
// Used in the main function for a single print and then not called again.
func GenerateJWTKeys() {
	accessKey := make([]byte, 32)
	_, err := rand.Read(accessKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JWT_ACCESS_KEY=\"%s\"\n", base64.StdEncoding.EncodeToString(accessKey))

	refreshKey := make([]byte, 32)
	_, err = rand.Read(refreshKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JWT_REFRESH_KEY=\"%s\"\n", base64.StdEncoding.EncodeToString(refreshKey))

	fmt.Println("Please add these two lines to the .env file")
}
