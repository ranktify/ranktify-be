package main

import (
	"context"
	"log"

	"github.com/ranktify/ranktify-be/config"
	"github.com/ranktify/ranktify-be/internal/dao"
)

func main() {
	db := config.SetupConnection()

	// reset the streaks with dao methods to be develop
	streaksDAO := dao.NewStreaksDAO(db)

	err := streaksDAO.ResetStreaksDaily(context.Background())
	if err != nil {
		log.Println("Couldn't reset streaks, error:", err.Error())
	}

	log.Println("Reset streaks successfully")
}
