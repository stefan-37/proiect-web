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