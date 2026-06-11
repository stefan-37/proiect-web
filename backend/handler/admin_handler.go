package handler

import (
	"backend/service"
	"github.com/gin-gonic/gin"
)

func AdminSignUp(c *gin.Context) {
	service.AdminSignUp(c, database)
}	

func AdminDelete(c *gin.Context) {
	service.AdminDelete(c, database)
}

func AdminUpdate(c *gin.Context) {
	service.AdminUpdate(c, database)
}

func AdminGet(c *gin.Context) {
	service.AdminGet(c, database)
}

func AdminCreateClass(c *gin.Context) {
	service.AdminCreateClass(c, database)
}

func GetTrainers(c *gin.Context) {
	service.GetTrainers(c, database)
}

func AdminTrainerDelete(c *gin.Context) {
	service.AdminTrainerDelete(c, database)
}

func GetProfit(c *gin.Context) {
	service.GetProfit(c, database)
}

func AdminUpdateClass(c *gin.Context) {
	service.AdminUpdateClass(c, database)
}

func AdminCreateTrainer(c *gin.Context) {
	service.AdminCreateTrainer(c, database)
}

func GetUsers(c *gin.Context) {
	service.GetUsers(c, database)
}

func AdminUpdateSubscription(c *gin.Context) {
	service.AdminUpdateSubscription(c, database)
}