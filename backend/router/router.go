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
	user.DELETE("/delete", middleware.AuthMiddleware("user"), handler.UserDelete)
	user.PUT("/update", middleware.AuthMiddleware("user"), handler.UserUpdate)
	user.GET("/get", middleware.AuthMiddleware("user"), handler.UserGet)
	user.POST("/subscribe", middleware.AuthMiddleware("user"), handler.UserSubscribe)
	user.GET("/subscription", middleware.AuthMiddleware("user"), handler.GetUserSubscriptions)

	admin := r.Group("/admin")
	admin.POST("/signup", handler.AdminSignUp)
	admin.POST("/login", handler.AdminLogin)
	admin.DELETE("/delete", middleware.AuthMiddleware("admin"), handler.AdminDelete)
	admin.PUT("/update", middleware.AuthMiddleware("admin"), handler.AdminUpdate)
	admin.GET("/get", middleware.AuthMiddleware("admin"), handler.AdminGet)

	trainer := r.Group("/trainer")
	trainer.POST("/login", handler.TrainerLogin)
	trainer.POST("/signup", handler.TrainerSignUp)
	trainer.DELETE("/delete", middleware.AuthMiddleware("trainer"), handler.TrainerDelete)
	trainer.PUT("/update", middleware.AuthMiddleware("trainer"), handler.TrainerUpdate)
	trainer.GET("/get", middleware.AuthMiddleware("trainer"), handler.TrainerGet)
	trainer.POST("/class", middleware.AuthMiddleware("trainer"), handler.TrainerCreateClass)

	r.GET("/subscriptions", handler.GetSubscriptions)
	r.GET("/classes", handler.GetClasses)

	return r
}