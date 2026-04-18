package repository

import (
	"backend/models"
	"gorm.io/gorm"
)

func CreateClass (class *models.Class, database *gorm.DB) error {
	err := database.Create(&class)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func GetClassByID(id uint, database *gorm.DB) (models.Class, error) {
	var class models.Class
	err := database.Where("id = ?", id).First(&class)
	if err.Error != nil {
		return models.Class{}, err.Error
	}
	return class, nil
}

func UpdateClass(class *models.Class, database *gorm.DB) error {
	err := database.Save(&class)	
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func DeleteClassByID(id uint, database *gorm.DB) error {
	err := database.Where("id = ?", id).Delete(&models.Class{})	
	if err.Error != nil {	
		return err.Error
	}
	return nil
}
