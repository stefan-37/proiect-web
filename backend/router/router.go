package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())

	registerUserRoutes(r)
	registerAdminRoutes(r)
	registerTrainerRoutes(r)
	registerPublicRoutes(r)

	return r
}
