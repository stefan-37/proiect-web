package router

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func registerTrainerRoutes(r *gin.Engine) {
	trainer := r.Group("/trainer")
	trainer.POST("/login", handler.TrainerLogin)
	//trainer.POST("/signup", handler.TrainerSignUp)
	trainer.DELETE("/delete", middleware.AuthMiddleware("trainer"), handler.TrainerDelete)
	trainer.PUT("/update", middleware.AuthMiddleware("trainer"), handler.TrainerUpdate)
	trainer.GET("/get", middleware.AuthMiddleware("trainer"), handler.TrainerGet)
	trainer.POST("/class", middleware.AuthMiddleware("trainer"), handler.TrainerCreateClass)
}
