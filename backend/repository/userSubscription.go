package repository

import (
	"backend/models"
	"gorm.io/gorm"
)

func CreateUserSubscription(userSubscription *models.UserSubscription, database *gorm.DB) error {
	err := database.Create(&userSubscription)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func GetUserSubscriptionByUserID(id uint, database *gorm.DB) (models.UserSubscription, error) {
	var userSubscription models.UserSubscription
	err := database.Where("user_id = ?", id).Last(&userSubscription)
	if err.Error != nil {
		return models.UserSubscription{}, err.Error
	}
	return userSubscription, nil
}

func UpdateUserSubscription(userSubscription models.UserSubscription, database *gorm.DB) error {
	err := database.Save(&userSubscription)
	if err.Error != nil {
		return err.Error
	}
	return nil	
}

func DeleteUserSubscriptionByID(id uint, database *gorm.DB) error {
	err := database.Where("id = ?", id).Delete(&models.UserSubscription{})
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func GetPaymentsByUserID(id uint, database *gorm.DB) ([]models.UserSubscription, error) {
	var payments []models.UserSubscription
	err := database.Where("user_id = ?", id).Find(&payments)
	if err.Error != nil {
		return nil, err.Error
	}
	return payments, nil
}

