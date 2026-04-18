package service

import (
	"gorm.io/gorm"
	"net/http"
	"backend/models"
	"github.com/gin-gonic/gin"
	"backend/repository"
)

func AdminSignUp(c *gin.Context, database *gorm.DB) {
	var body models.Admin
	if c.BindJSON(&body) != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Failed to read body",
		})
		return
	}

	user, err := models.AdminFactory(
		models.AdminWithName(body.Name),
		models.AdminWithEmail(body.Email),
		models.AdminWithPassword(body.Password),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid field(s)",
		})
		return
	}

	if repository.CreateAdmin(user, database) != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to create admin",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin created successfully",
	})

}