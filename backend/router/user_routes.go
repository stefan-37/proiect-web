package router

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(r *gin.Engine) {
	user := r.Group("/user")
	user.POST("/signup", handler.UserSignUp)
	user.POST("/login", handler.UserLogin)
	user.DELETE("/delete", middleware.AuthMiddleware("user"), handler.UserDelete)
	user.PUT("/update", middleware.AuthMiddleware("user"), handler.UserUpdate)
	user.POST("/reset-password", middleware.AuthMiddleware("user"), handler.UserResetPassword)
	user.GET("/get", middleware.AuthMiddleware("user"), handler.UserGet)
	user.POST("/logout", middleware.AuthMiddleware("user"), handler.UserLogout)
	user.POST("/subscribe", middleware.AuthMiddleware("user"), handler.UserSubscribe)
	user.GET("/subscription", middleware.AuthMiddleware("user"), handler.GetUserSubscriptions)
	user.POST("/checkin", middleware.AuthMiddleware("user"), handler.CheckIn)
	user.POST("/checkout", middleware.AuthMiddleware("user"), handler.CheckOut)
	user.POST("/book", middleware.AuthMiddleware("user"), handler.BookClass)
	user.GET("/reservations", middleware.AuthMiddleware("user"), handler.GetReservations)
	user.DELETE("/reservation", middleware.AuthMiddleware("user"), handler.DeleteReservation)
	user.GET("/payments", middleware.AuthMiddleware("user"), handler.GetUserPayments)
}
