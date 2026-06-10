package service

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"backend/models"
	"github.com/gin-gonic/gin"
	"backend/repository"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
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
		models.UserWithPhone(body.Phone),
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

	sendEmail(user.Email)

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
		Phone    string `json:"phone"`
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
	if body.Phone != "" {
		userData.Phone = body.Phone
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

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, nil)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user subscriptions",
		})
		return
	}

	c.JSON(http.StatusOK, subscriptions)

}

func BookClass(c *gin.Context, database *gorm.DB) {
	id, _ := c.Get("ID")

	var body struct {
		ClassID uint `json:"class_id"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	class, err := repository.GetClassByID(body.ClassID, database)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Class not found",
		})
		return
	}

	bookingSituation, err := models.BookingSituationFactory(
		models.BookingSituationWithUserID(id.(uint)),
		models.BookingSituationWithClassID(body.ClassID),
		models.BookingSituationWithAdminID(class.AdminID),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create booking situation",
		})
		return
	}

	if err := repository.BookClass(bookingSituation, database); err != nil {
		if errors.Is(err, repository.ErrClassFull) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Class is full",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to book class",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Class booked successfully",
	})
}

func UserLogout(c *gin.Context) {
	 	c.SetSameSite(http.SameSiteNoneMode)
        
        c.SetCookie("key", "", -1, "", "", true, true)

        c.JSON(http.StatusOK, gin.H{
                "message": "Logout successful",
        })
}

func UserResetPassword(c *gin.Context, database *gorm.DB) {
	id, _ := c.Get("ID")

	userData, err := repository.GetUserByID(id.(uint), database)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user data",
		})
		return
	}

	newPassword, err := randomPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate password",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	if err := sendPasswordEmail(userData.Email, newPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send password email",
		})
		return
	}

	userData.Password = string(hash)

	if repository.UpdateUser(&userData, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

func randomPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func GetReservations(c *gin.Context, database *gorm.DB) {
	id, _ := c.Get("ID")
	reservations, err := repository.GetReservationsByUserID(id.(uint), database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get reservations",
		})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

func DeleteReservation(c *gin.Context, database *gorm.DB) {
	id, _ := c.Get("ID")
	var body struct {
		ClassID uint `json:"class_id"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",})
		return
	}
	if err := repository.DeleteReservationByUserIDAndClassID(id.(uint), body.ClassID, database); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Reservation not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete reservation",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation deleted successfully",
	})
}

func ForgotPassword(c *gin.Context, database *gorm.DB) {
	var body struct {
		Email string `json:"email"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	userData, err := repository.GetUserByEmail(body.Email, database)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user data",
		})
		return
	}

	newPassword, err := randomPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate password",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	if err := sendPasswordEmail(userData.Email, newPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send password email",
		})
		return
	}

	userData.Password = string(hash)
	if repository.UpdateUser(&userData, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}
