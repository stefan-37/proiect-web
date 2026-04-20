package handler

import(
	"backend/service"
	"github.com/gin-gonic/gin"
)

func GetSubscriptions(c *gin.Context) {
	service.GetSubscriptions(c, database)
}