package handler

import (
	"backend/service"
	"github.com/gin-gonic/gin"
)

func AdminSignUp(c *gin.Context) {
	service.AdminSignUp(c, database)
}	
