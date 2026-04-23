package main

import (
	"backend/db"
	"backend/models"
	"backend/router"
	"backend/seed"
)

func main() {
	database := db.GetDB()
	database.AutoMigrate(&models.User{}, &models.Admin{}, &models.Trainer{}, &models.Subscription{}, &models.UserSubscription{}, &models.Class{})
	seed.LoadPlans("seed/plans.json", database)


	router := router.SetupRouter()
	router.Run(":8080")
}