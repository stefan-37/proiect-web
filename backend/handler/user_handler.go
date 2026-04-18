package handler

import(
	"backend/service"
	"github.com/gin-gonic/gin"
)

func UserSignUp(c *gin.Context) {
	service.UserSignUp(c, database)
}

func UserDelete(c *gin.Context) {
	service.UserDelete(c, database)
}

func UserUpdate(c *gin.Context) {
	service.UserUpdate(c, database)
}

func UserGet(c *gin.Context) {
	service.UserGet(c, database)
}