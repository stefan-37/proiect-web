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
	id, err := c.Get("ID")

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user ID",
		})
		return
	}

	if repository.DeleteUserByID(id.(uint), database) != nil{
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

	id, _ := c.Get("ID")
	userData, err := repository.GetUserByID(id.(uint), database)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user data",
		})
		return
	}

	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Name != "" {
		userData.Name = body.Name
	}
	if body.Email != "" {
		userData.Email = body.Email
	}
	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
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

	id, _ := c.Get("ID")
	userData, err := repository.GetUserByID(id.(uint), database)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user data",
		})
		return
	}

	c.JSON(http.StatusOK, userData)

}

func GetUserSubscriptions(c *gin.Context, database *gorm.DB) {

	id, _ := c.Get("ID")

	subscriptions, err := repository.GetUserSubscriptionByUserID(id.(uint), database)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user subscriptions",
		})
		return
	}

	c.JSON(http.StatusOK, subscriptions)

}