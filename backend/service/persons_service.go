package service

import (
	"gorm.io/gorm"

	"backend/models"
	"backend/repository"
	"github.com/gin-gonic/gin"
	"time"
)

func CheckIn(c *gin.Context, database *gorm.DB) {
	id, _ := c.Get("ID")
	body := models.NewPerson(id.(uint),time.Now(),nil)
	if repository.CheckInPerson(body, database) != nil {
		c.JSON(500, gin.H{
			"error": "Failed to check in",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Checked in successfully",
	})
}

func CheckOut(c *gin.Context, database *gorm.DB) {
	id, _ := c.Get("ID")
	timeNow := time.Now()
	person := models.NewPerson(id.(uint),time.Time{},&timeNow)
	if repository.CheckOutPerson(person.PersonID, database) != nil {
		c.JSON(500, gin.H{
			"error": "Failed to check out",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Checked out successfully",
	})
}

func GetPersonsCount(c *gin.Context, database *gorm.DB) {
	count, err := repository.GetPersonsCount(database)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to get persons count",
		})
		return
	}
	c.JSON(200, gin.H{
		"count": count,
	})
}