package main

import (
	"fmt"
	"os"

	"github.com/cuijxin/my-n9e/config"
	"github.com/cuijxin/my-n9e/models"
)

func main() {
	if err := config.Parse(); err != nil {
		fmt.Println("cannot parse configuration file:", err)
		os.Exit(1)
	}
	models.InitMySQL(config.Config.MySQL)
	if err := models.DB.AutoMigrate(
		new(models.User),
		new(models.Configs),
	); err != nil {
		fmt.Printf("database auto migrate failed: %v", err)
		os.Exit(1)
	}
	fmt.Println("database auto migrate success")
	os.Exit(0)
}
