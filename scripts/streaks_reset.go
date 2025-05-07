package main

import (
	"fmt"

	"github.com/ranktify/ranktify-be/config"
)

func main() {
	db := config.SetupConnection()

	fmt.Println(db)
}
