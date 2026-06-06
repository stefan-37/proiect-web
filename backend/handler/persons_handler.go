package handler

import (
	"backend/service"
	"github.com/gin-gonic/gin"
)

func CheckIn(c *gin.Context) {
	service.CheckIn(c, database)
}

func CheckOut(c *gin.Context) {
	service.CheckOut(c, database)
}

func GetPersonsCount(c *gin.Context) {
	service.GetPersonsCount(c, database)
}