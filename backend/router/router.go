package router

import (
	"backend/handler"
	"github.com/gin-gonic/gin"
	"backend/middleware"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	// CORS: let the frontend origins call the API with credentials (the auth cookie).
	// Explicit origins only — "*" is not allowed together with credentials.
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		switch origin {
		case "http://localhost:5500", "http://127.0.0.1:5500":
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
			c.Header("Vary", "Origin")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	user := r.Group("/user")
	user.POST("/signup",handler.UserSignUp)
	user.POST("/login", handler.UserLogin)
	user.DELETE("/delete", middleware.AuthMiddleware("user"), handler.UserDelete)
	user.PUT("/update", middleware.AuthMiddleware("user"), handler.UserUpdate)
	user.GET("/get", middleware.AuthMiddleware("user"), handler.UserGet)
	user.POST("/logout", middleware.AuthMiddleware("user"), handler.UserLogout)
	user.POST("/subscribe", middleware.AuthMiddleware("user"), handler.UserSubscribe)
	user.GET("/subscription", middleware.AuthMiddleware("user"), handler.GetUserSubscriptions)
	user.POST("/checkin", middleware.AuthMiddleware("user"), handler.CheckIn)
	user.POST("/checkout", middleware.AuthMiddleware("user"), handler.CheckOut)
	user.POST("/book", middleware.AuthMiddleware("user"), handler.BookClass)

	admin := r.Group("/admin")
	//admin.POST("/signup", handler.AdminSignUp)
	admin.POST("/login", handler.AdminLogin)
	admin.DELETE("/delete", middleware.AuthMiddleware("admin"), handler.AdminDelete)
	admin.PUT("/update", middleware.AuthMiddleware("admin"), handler.AdminUpdate)
	admin.GET("/get", middleware.AuthMiddleware("admin"), handler.AdminGet)
	admin.POST("/class", middleware.AuthMiddleware("admin"), handler.AdminCreateClass)

	trainer := r.Group("/trainer")
	trainer.POST("/login", handler.TrainerLogin)
	//trainer.POST("/signup", handler.TrainerSignUp)
	trainer.DELETE("/delete", middleware.AuthMiddleware("trainer"), handler.TrainerDelete)
	trainer.PUT("/update", middleware.AuthMiddleware("trainer"), handler.TrainerUpdate)
	trainer.GET("/get", middleware.AuthMiddleware("trainer"), handler.TrainerGet)
	trainer.POST("/class", middleware.AuthMiddleware("trainer"), handler.TrainerCreateClass)

	r.GET("/subscriptions", handler.GetSubscriptions)
	r.GET("/classes", handler.GetClasses)
	r.GET("/personsCount", handler.GetPersonsCount)

	return r
}