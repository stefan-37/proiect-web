package service

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"backend/models"
	"backend/repository"
	"net/http"
	"time"
	)

func UserSubscribe(c *gin.Context, database *gorm.DB) {

	user, ok := c.Get("ID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read user ID",
		})
		return
	}
	var subscriptionID uint
	if err := c.BindJSON(&subscriptionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read subscription ID",
		})
		return
		}

	if _, err := repository.GetSubscriptionByID(subscriptionID, database); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subscription not found",
		})
		return
	}

	userSubscription, err := models.UserSubscriptionFactory(
		models.UserSubscriptionWithUserID(user.(uint)),
		models.UserSubscriptionWithSubscriptionID(subscriptionID),
		models.UserSubscriptionWithStartedAt(time.Now()),
		models.UserSubscriptionWithExpiresAt(time.Now().AddDate(0, 1, 0)),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user subscription",
		})
		return
	}

	if repository.CreateUserSubscription(userSubscription, database) != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to create user subscription",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User subscribed successfully",
	})

}

func GetSubscriptions(c *gin.Context, database *gorm.DB) {

	subscriptions, err := repository.GetAllSubscriptions(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve subscriptions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscriptions": subscriptions,
	})

}

func CheckSubscriptions(database *gorm.DB) {

	database.Where("expires_at < ?",time.Now()).Delete(&models.UserSubscription{})

}