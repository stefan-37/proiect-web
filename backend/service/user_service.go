package service

import (
	"gorm.io/gorm"
	"net/http"
	"backend/models"
	"github.com/gin-gonic/gin"
	"backend/repository"
)

func UserSignUp(c *gin.Context, database *gorm.DB) {
	var body models.User

	if c.BindJSON(&body) != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Failed to read body",
		})
		return
	}

	user, err := models.UserFactory(
		models.UserWithName(body.Name),
		models.UserWithEmail(body.Email),
		models.UserWithPassword(body.Password),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid field(s)",
		})
		return

	}

	if repository.CreateUser(user, database) != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"User created successfully",
	})
	
}