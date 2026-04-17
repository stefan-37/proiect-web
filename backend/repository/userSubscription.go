package repository

import (
	"backend/models"
	"gorm.io/gorm"
)

func CreateUserSubscription(userSubscription models.UserSubscription, database *gorm.DB) error {
	err := database.Create(&userSubscription)
	if err != nil {
		return err.Error
	}
	return nil
}

func GetUserSubscriptionByID(id uint, database *gorm.DB) (models.UserSubscription, error) {
	var userSubscription models.UserSubscription
	err := database.Where("id = ?", id).First(&userSubscription)
	if err != nil {
		return models.UserSubscription{}, err.Error
	}	
	return userSubscription, nil
}

func UpdateUserSubscription(userSubscription models.UserSubscription, database *gorm.DB) error {
	err := database.Save(&userSubscription)
	if err != nil {
		return err.Error
	}
	return nil	
}

func DeleteUserSubscriptionByID(id uint, database *gorm.DB) error {
	err := database.Where("id = ?", id).Delete(&models.UserSubscription{})
	if err != nil {
		return err.Error
	}
	return nil
}

