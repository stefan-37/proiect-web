package service

import (
	"gorm.io/gorm"
	"net/http"
	"backend/models"
	"github.com/gin-gonic/gin"
	"backend/repository"
	"golang.org/x/crypto/bcrypt"
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


func UserDelete(c *gin.Context, database *gorm.DB) {
	user, err := c.Get("user")

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user ID",
		})
		return
	}

	userData, err := user.(models.User)

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user data",
		})
		return
	}

	if repository.DeleteUserByID(userData.ID, database) != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to delete user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func UserUpdate(c *gin.Context, database *gorm.DB) {

	user, _ := c.Get("user")
	userData, ok := user.(models.User)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user data",
		})
		return
	}

	var body struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Name != nil {
		userData.Name = *body.Name
	}
	if body.Email != nil {
		userData.Email = *body.Email
	}
	if body.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		userData.Password = string(hash)
	}

	if repository.UpdateUser(&userData, database) != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})

}

func UserGet(c *gin.Context, database *gorm.DB) {

	user, _ := c.Get("user")
	userData, ok := user.(models.User)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user data",
		})
		return
	}

	c.JSON(http.StatusOK, userData)

}