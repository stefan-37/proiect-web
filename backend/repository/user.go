package repository

import (
	"backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(user *models.User, database *gorm.DB) error {
	hash, hasherr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hasherr != nil {
		return hasherr
	}
	user.Password = string(hash)

	err := database.Create(&user)
	if err != nil {
		return err.Error
	}
	return nil
}

func GetUserByEmail(email string, database *gorm.DB) (models.User, error) {
	var user models.User
	err := database.Where("email = ?", email).First(&user)

	if err != nil {
		return models.User{}, err.Error
	}
	return user, nil
}

func GetUserByID(id uint, database *gorm.DB) (models.User, error) {
	var user models.User
	err := database.Where("id = ?", id).First(&user)
	if err != nil {
		return models.User{}, err.Error
	}
	return user, nil
}

func UpdateUser(user *models.User, database *gorm.DB) error {
	err := database.Save(&user)
	if err != nil {
		return err.Error
	}
	return nil
}

func DeleteUserByEmail(email string, database *gorm.DB) error {
	err := database.Where("email = ?", email).Delete(&models.User{})
	if err != nil {
		return err.Error
	}
	return nil
}
