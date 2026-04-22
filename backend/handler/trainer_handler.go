package handler

import (
	"backend/service"
	"github.com/gin-gonic/gin"
)

func TrainerSignUp(c *gin.Context) {
	service.TrainerSignUp(c, database)
}

func TrainerDelete(c *gin.Context) {
	service.TrainerDelete(c, database)
}

func TrainerUpdate(c *gin.Context) {
	service.TrainerUpdate(c, database)
}

func TrainerGet(c *gin.Context) {
	service.TrainerGet(c, database)
}

func TrainerCreateClass(c *gin.Context) {
	service.TrainerCreateClass(c, database)
}