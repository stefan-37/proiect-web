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
	dberr := database.AutoMigrate(&models.User{}, &models.Admin{}, &models.Trainer{}, &models.Subscription{}, &models.UserSubscription{}, &models.Class{}, &models.BookingSituation{}, &models.Person{})
	if dberr != nil {
		log.Fatal(dberr)
	}

	err := seed.LoadPlans("seed/plans.json", database)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			service.CheckSubscriptions(database)
			log.Printf("Checked subscriptions at %s", time.Now().Format(time.RFC3339))
			time.Sleep(24 * time.Hour)
		}
	}()

	go func() {
		for {
			seed.LoadAdmins("seed/admins.json", database)
			seed.LoadTrainers("seed/trainers.json", database)
			seed.LoadClasses("seed/classes.json", database)
			time.Sleep(5 * time.Minute)
		}
	}()

	router := router.SetupRouter()
	if router.Run(":8080") != nil {
		log.Fatal("Unable to start server")
	}
}
