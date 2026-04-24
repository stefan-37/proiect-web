package main

import (
	"backend/db"
	"backend/models"
	"backend/router"
	"backend/seed"
	"backend/service"
	"log"
	"time"
)

func main() {
	database := db.GetDB()
	database.AutoMigrate(&models.User{}, &models.Admin{}, &models.Trainer{}, &models.Subscription{}, &models.UserSubscription{}, &models.Class{})
	seed.LoadPlans("seed/plans.json", database)

	go func() {
		for {
			service.CheckSubscriptions(database)
			log.Printf("Checked subscriptions at %s", time.Now().Format(time.RFC3339))
			time.Sleep(24 * time.Hour)
		}
	}()

	router := router.SetupRouter()
	router.Run(":8080")
}
