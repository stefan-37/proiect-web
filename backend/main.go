package main

import (
	"backend/db"
	"backend/models"
	"backend/router"
)

func main() {
	database := db.GetDB()
	database.AutoMigrate(&models.User{}, &models.Admin{}, &models.Trainer{})


	router := router.SetupRouter()
	router.Run(":8080")
}