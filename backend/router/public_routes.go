package router

import (
	"backend/handler"
	"github.com/gin-gonic/gin"
)

func registerPublicRoutes(r *gin.Engine) {
	r.GET("/subscriptions", handler.GetSubscriptions)
	r.GET("/classes", handler.GetClasses)
	r.GET("/personsCount", handler.GetPersonsCount)
	r.POST("/forgot-password", handler.ForgotPassword)
}
