package repository

import (
	"backend/models"
	"gorm.io/gorm"
)

func CreateSubscription(subscription *models.Subscription, database *gorm.DB) error {
	err := database.Create(&subscription)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func GetSubscriptionByID(id uint, database *gorm.DB) (models.Subscription, error) {
	var subscription models.Subscription
	err := database.Where("id = ?", id).First(&subscription)
	if err.Error != nil {
		return models.Subscription{}, err.Error
	}
	return subscription, nil
}

func UpdateSubscription(subscription *models.Subscription, database *gorm.DB) error {
	err := database.Save(&subscription)
	if err.Error != nil {
		return err.Error
	}	
	return nil
}

func DeleteSubscriptionByID(id uint, database *gorm.DB) error {
	err := database.Where("id = ?", id).Delete(&models.Subscription{})	
	if err.Error != nil {
		return err.Error
	}	
	return nil
}

func GetAllSubscriptions(database *gorm.DB) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := database.Find(&subscriptions)
	if err.Error != nil {
		return nil, err.Error
	}
	return subscriptions, nil
}