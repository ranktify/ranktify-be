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
	generateKey := func() string {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			panic(err)
		}
		return base64.StdEncoding.EncodeToString(key)
	}

	accessKey := generateKey()
	refreshKey := generateKey()
	fmt.Printf("JWT_ACCESS_KEY=\"%s\"\n", accessKey)
	fmt.Printf("JWT_REFRESH_KEY=\"%s\"\n", refreshKey)

	fmt.Println("Please add these two lines to the .env file")
}
