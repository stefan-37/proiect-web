package service

import (
	"backend/models"
	"backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AdminSignUp(c *gin.Context, database *gorm.DB) {
	var body models.Admin
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	admin, err := models.AdminFactory(
		models.AdminWithName(body.Name),
		models.AdminWithEmail(body.Email),
		models.AdminWithPassword(body.Password),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid field(s)",
		})
		return
	}

	if repository.CreateAdmin(admin, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create admin",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin created successfully",
	})

}

func AdminDelete(c *gin.Context, database *gorm.DB) {

	admin, err := c.Get("admin")

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin ID",
		})
		return
	}

	adminData, err := admin.(models.Admin)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin data",
		})
		return
	}

	if repository.DeleteAdminByID(adminData.ID, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete admin",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin deleted successfully",
	})

}

func AdminUpdate(c *gin.Context, database *gorm.DB) {
	admin, err := c.Get("admin")

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin ID",
		})
		return
	}

	adminData, err := admin.(models.Admin)

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin data",
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
		adminData.Name = body.Name
	}
	if body.Email != "" {
		adminData.Email = body.Email
	}
	if body.Password != "" {
		hash, hasherr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if hasherr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		adminData.Password = string(hash)
	}

	if repository.UpdateAdmin(&adminData, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update admin",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin updated successfully",
	})

}

func AdminGet(c *gin.Context, database *gorm.DB) {
	admin, _ := c.Get("admin")
	adminData, ok := admin.(models.Admin)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    adminData.ID,
		"name":  adminData.Name,
		"email": adminData.Email,
	})

}
