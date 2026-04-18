package router

import (
	"backend/handler"
	"github.com/gin-gonic/gin"
	"backend/middleware"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	user := r.Group("/user")
	user.POST("/signup",handler.UserSignUp)
	user.POST("/login", handler.UserLogin)
	user.DELETE("/delete", middleware.UserAuthMiddleware, handler.UserDelete)
	user.PUT("/update", middleware.UserAuthMiddleware, handler.UserUpdate)


	admin := r.Group("/admin")
	admin.POST("/signup", handler.AdminSignUp)
	admin.POST("/login", handler.AdminLogin)

	trainer := r.Group("/trainer")
	trainer.POST("/login", handler.TrainerLogin)
	//trainer.POST("/signup", handler.TrainerSignUp)

	return r
}