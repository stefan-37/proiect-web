package service

import (
	"gorm.io/gorm"
	"net/http"
	"backend/models"
	"github.com/gin-gonic/gin"
	"backend/repository"
	"strconv"
	"time"
)

func CreateClass(c *gin.Context, database *gorm.DB) {
	var body models.Class

	if c.BindJSON(&body) != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Failed to read body",
		})
		return
	}

	class, err := models.ClassFactory(
		models.ClassWithName(body.Name),
		models.ClassWithDescription(body.Description),
		models.ClassWithCapacity(body.Capacity),
		models.ClassWithAdminID(body.AdminID),
		models.ClassWithScheduledAt(body.ScheduledAt),
		models.ClassWithTrainerID(body.TrainerID),
		models.ClassWithPrice(body.Price),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid field(s)",
		})
		return
	}

	if repository.CreateClass(class, database) != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to create class",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Class created successfully",
	})

}

func AdminCreateClass(c *gin.Context, database *gorm.DB) {
	id, ok := c.Get("ID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin ID",
		})
		return
	}

	var body struct {
		ScheduledAt time.Time `json:"date"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Capacity    uint      `json:"capacity"`
		Price       float64   `json:"price"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Name == "" || body.Capacity == 0 || body.Price == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid field(s)",
		})
		return
	}

	class := &models.Class{
		Name:        body.Name,
		Description: body.Description,
		Capacity:    body.Capacity,
		ScheduledAt: body.ScheduledAt,
		TrainerID:   0,
		Users:       0,
		AdminID:     id.(uint),
		Price:       body.Price,
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

func DeleteClass(c *gin.Context, database *gorm.DB) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid class ID",
		})
		return
	}

	if repository.DeleteClassByID(uint(idInt), database) != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to delete class",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Class deleted successfully",
	})

}

func UpdateClass(c *gin.Context, database *gorm.DB) {
	
	class, _ := c.Get("class")
	classData, err := class.(models.Class)

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read class data",
		})
		return
	}

	if repository.UpdateClass(&classData, database) != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to update class",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Class updated successfully",
	})
}

func GetClasses(c *gin.Context, database *gorm.DB) {
	var classes []models.Class
	err := database.Find(&classes)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve classes",
		})
		return
	}
	c.JSON(http.StatusOK, classes)

}