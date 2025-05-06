package main

import (
	"fmt"

	"github.com/ranktify/ranktify-be/config"
)

func main() {
	db := config.SetupConnection()

	// reset the streaks with dao methods to be develop
	fmt.Println(db)
}
