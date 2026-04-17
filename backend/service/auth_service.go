package service

import (
	"backend/config"
	"backend/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserLogin(c *gin.Context, database *gorm.DB) {
	var body struct {
		email    string
		password string
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	}

	user, DBerr := repository.GetUserByEmail(body.email, database)

	if DBerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"subject": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Key))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("key", tokenString, 3600*24*7, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})

}
