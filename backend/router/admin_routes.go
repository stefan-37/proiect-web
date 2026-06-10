package router

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func registerAdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	//admin.POST("/signup", handler.AdminSignUp)
	admin.POST("/login", handler.AdminLogin)
	admin.DELETE("/delete", middleware.AuthMiddleware("admin"), handler.AdminDelete)
	admin.PUT("/update", middleware.AuthMiddleware("admin"), handler.AdminUpdate)
	admin.GET("/get", middleware.AuthMiddleware("admin"), handler.AdminGet)
	admin.POST("/class", middleware.AuthMiddleware("admin"), handler.AdminCreateClass)
}
