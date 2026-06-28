package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(corsMiddleware())

	registerUserRoutes(r)
	registerAdminRoutes(r)
	registerTrainerRoutes(r)
	registerPublicRoutes(r)

	return r
}
