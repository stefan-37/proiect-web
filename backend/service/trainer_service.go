package service

import (
	"gorm.io/gorm"
	"backend/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func TrainerSignUp(c *gin.Context, database *gorm.DB) {

	var body models.Trainer

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	trainer, err := models.TrainerFactory(
		models.TrainerWithName(body.Name),
		models.TrainerWithEmail(body.Email),
		models.TrainerWithPassword(body.Password),
		models.TrainerWithAdminID(body.AdminID),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field(s)"})
		return
	}

	if err := repository.CreateTrainer(trainer, database); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trainer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trainer created successfully",
	})

}

func TrainerDelete(c *gin.Context, database *gorm.DB) {
	id, ok := c.Get("ID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer ID",
		})
		return
	}

	if err := repository.DeleteTrainerByID(id.(uint), database); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete trainer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trainer deleted successfully",
	})

}

func TrainerUpdate(c *gin.Context, database *gorm.DB) {
	id, _ := c.Get("ID")
	trainerData, err := repository.GetTrainerByID(id.(uint), database)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer data",
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
		trainerData.Name = body.Name
	}
	if body.Email != "" {
		trainerData.Email = body.Email
	}
	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		trainerData.Password = string(hash)
	}

	if repository.UpdateTrainer(&trainerData, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update trainer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trainer updated successfully",
	})

}

func TrainerGet(c *gin.Context, database *gorm.DB) {
	id, ok := c.Get("ID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer ID",
		})
		return
	}

	trainerData, err := repository.GetTrainerByID(id.(uint), database)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer data",
		})
		return
	}

	c.JSON(http.StatusOK, trainerData)
}

func TrainerCreateClass(c *gin.Context, database *gorm.DB) {

	id, ok := c.Get("ID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer ID",
		})
		return
	}

	var body struct {
		ScheduledAt time.Time `json:"date"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		TrainerID   uint      `json:"trainer_id"`
		Capacity    uint      `json:"capacity"`
		Users       uint	  `json:"users"`
		AdminID     uint      `json:"admin_id"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	body.TrainerID = id.(uint)

	trainer, err := repository.GetTrainerByID(id.(uint), database)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer data",
		})
		return
	}

	class, err := models.ClassFactory(
		models.ClassWithName(body.Name),
		models.ClassWithDescription(body.Description),
		models.ClassWithCapacity(body.Capacity),
		models.ClassWithTrainerID(body.TrainerID),
		models.ClassWithScheduledAt(body.ScheduledAt),
		models.ClassWithAdminID(trainer.AdminID),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field(s)"})
		return
	}

	if repository.CreateClass(class, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create class",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Class created successfully",
	})

}