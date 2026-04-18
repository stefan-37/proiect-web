package repository

import (
	"backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateTrainer(trainer *models.Trainer, database *gorm.DB) error {
	hash, hasherr := bcrypt.GenerateFromPassword([]byte(trainer.Password),bcrypt.DefaultCost)
	if hasherr != nil {
		return hasherr
	}
	trainer.Password = string(hash)

	err := database.Create(&trainer)
	if err.Error != nil{
		return err.Error
	}
	return nil
}

func GetTrainerByEmail(email string, database *gorm.DB) (models.Trainer, error) {
	var trainer models.Trainer
	err := database.Where("email = ?", email).First(&trainer)
	if err.Error != nil{
		return models.Trainer{},err.Error
	}
	return trainer, nil
}

func GetTrainerByID(id uint, database *gorm.DB) (models.Trainer, error){
	var trainer models.Trainer
	err := database.Where("id = ?", id).First(&trainer)
	if err.Error != nil{
		return models.Trainer{},err.Error
	}
	return trainer, nil
}

func  UpdateTrainer(trainer *models.Trainer, database *gorm.DB) error{
	err := database.Save(&trainer)
	if err.Error != nil{
		return err.Error
	}
	return nil
}

func DeleteTrainerByID(id uint, database *gorm.DB) error{
	err := database.Where("id = ?", id).Delete(&models.Trainer{})
	if err.Error != nil{
		return err.Error
	}
	return nil
}	