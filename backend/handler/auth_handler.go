package handler

import (
	"backend/db"
	"backend/service"

	"github.com/gin-gonic/gin"
)

var database = db.GetDB()

func AdminLogin(c *gin.Context) {
	service.AdminLogin(c, database)
}

func UserLogin(c *gin.Context) {
	service.UserLogin(c, database)
}

func TrainerLogin(c *gin.Context) {
	service.TrainerLogin(c, database)
}


