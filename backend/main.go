package main

import (
	"backend/db"
	"backend/models"
	//"fmt"
)

func main() {
	database := db.GetDB()
	database.AutoMigrate(&models.User{}, &models.Admin{}, &models.Trainer{})
}