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
	user.GET("/get", middleware.UserAuthMiddleware, handler.UserGet)
	user.POST("/subscribe", middleware.UserAuthMiddleware, handler.UserSubscribe)
	user.GET("/subscriptions", middleware.UserAuthMiddleware, handler.GetSubscriptions)
	user.GET("/subscription", middleware.UserAuthMiddleware, handler.GetUserSubscriptions)

	admin := r.Group("/admin")
	admin.POST("/signup", handler.AdminSignUp)
	admin.POST("/login", handler.AdminLogin)
	admin.DELETE("/delete", middleware.AdminAuthMiddleware, handler.AdminDelete)
	admin.PUT("/update", middleware.AdminAuthMiddleware, handler.AdminUpdate)
	admin.GET("/get", middleware.AdminAuthMiddleware, handler.AdminGet)
	admin.GET("/subscriptions", middleware.AdminAuthMiddleware, handler.GetSubscriptions)


	trainer := r.Group("/trainer")
	trainer.POST("/login", handler.TrainerLogin)
	trainer.POST("/signup", handler.TrainerSignUp)
	trainer.DELETE("/delete", middleware.TrainerAuthMiddleware, handler.TrainerDelete)
	trainer.PUT("/update", middleware.TrainerAuthMiddleware, handler.TrainerUpdate)
	trainer.GET("/get", middleware.TrainerAuthMiddleware, handler.TrainerGet)
	trainer.GET("/subscriptions", middleware.TrainerAuthMiddleware, handler.GetSubscriptions)

	return r
}