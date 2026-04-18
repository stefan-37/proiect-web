package service

import (
	"gorm.io/gorm"
	"backend/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/repository"
	"golang.org/x/crypto/bcrypt"
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
	trainer, err := c.Get("trainer")

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer ID",
		})
		return
	}

	trainerData, err := trainer.(models.Trainer)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer data",
		})
		return
	}

	if err := repository.DeleteTrainerByID(trainerData.ID, database); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete trainer",
		})
		return
	}
}

func TrainerUpdate(c *gin.Context, database *gorm.DB) {
	trainer, _ := c.Get("trainer")
	trainerData, ok := trainer.(models.Trainer)

	if !ok {
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
	trainer, err := c.Get("trainer")

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer ID",
		})
		return
	}

	trainerData, err := trainer.(models.Trainer)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read trainer data",
		})
		return
	}

	c.JSON(http.StatusOK, trainerData)
}