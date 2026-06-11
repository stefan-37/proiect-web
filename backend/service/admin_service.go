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

	id, ok := c.Get("ID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin ID",
		})
		return
	}

	if repository.DeleteAdminByID(id.(uint), database) != nil {
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
	id, ok := c.Get("ID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin ID",
		})
		return
	}

	adminData, err := repository.GetAdminByID(id.(uint), database)

	if err != nil {
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
	id, ok := c.Get("ID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin ID",
		})
		return
	}

	adminData, err := repository.GetAdminByID(id.(uint), database)
	if err != nil {
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

func GetTrainers(c *gin.Context, database *gorm.DB) {
	var trainers []models.Trainer
	trainers, err := repository.GetAllTrainers(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get trainers",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"trainers": trainers,
	})
}

func AdminTrainerDelete(c *gin.Context, database *gorm.DB) {
	var body struct {
		ID uint `json:"id"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if err := repository.DeleteTrainerByID(body.ID, database); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete trainer",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Trainer deleted successfully",
	})
}

func GetProfit(c *gin.Context, database *gorm.DB) {
	userSubscriptions, err := repository.GetAllUserSubscriptions(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user subscriptions",
		})
		return
	}
	bookings, err := repository.GetAllBookings(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get bookings",
		})
		return
	}

	subscriptions, err := repository.GetAllSubscriptions(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get subscriptions",
		})
		return
	}

	classes, err := repository.GetAllClasses(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get classes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userSubscriptions": userSubscriptions,
		"bookings":          bookings,
		"subscriptions":     subscriptions,
		"classes":           classes,
	})
}

func AdminUpdateClass(c *gin.Context, database *gorm.DB) {
	var body struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       float64   `json:"price" gorm:"not null"`
		Capacity    uint      `json:"capacity" gorm:"not null"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	class, err := repository.GetClassByID(body.ID, database)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read class data",
		})
		return
	}

	class.Name = body.Name
	class.Description = body.Description
	class.Price = body.Price
	class.Capacity = body.Capacity

	if repository.UpdateClass(&class, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update class",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Class updated successfully",
	})
}

func AdminCreateTrainer(c *gin.Context, database *gorm.DB) {
	var body struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Class       string `json:"class"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	id, ok := c.Get("ID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read admin ID",
		})
		return
	}

	trainer, err := models.TrainerFactory(
		models.TrainerWithName(body.Name),
		models.TrainerWithEmail(body.Email),
		models.TrainerWithPassword(body.Password),
		models.TrainerWithClass(body.Class),
		models.TrainerWithAdminID(id.(uint)),
	)	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid field(s)",
		})
		return
	}

	if repository.CreateTrainer(trainer, database) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create trainer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trainer created successfully",
	})
}

func GetUsers(c *gin.Context, database *gorm.DB) {
	users, err := repository.GetAllUsers(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get users",	
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}