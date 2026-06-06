package repository

import (
	"backend/models"
	"gorm.io/gorm"
	"time"
)

func CheckInPerson(person *models.Person, database *gorm.DB) error {
	err := database.Create(person)
	if err != nil {
		return err.Error
	}
	return nil
}

func CheckOutPerson(personID uint, database *gorm.DB) error {
	var person models.Person
	err := database.Where("person_id = ? AND checked_out IS NULL", personID).
                Order("checked_in DESC").First(&person)
	if err.Error != nil {
		return err.Error
	}
    checkedOut := time.Now()
    person.CheckedOut = &checkedOut
    return database.Save(&person).Error
}

func GetPersonsCount(database *gorm.DB) (int64, error) {
	var count int64
	err := database.Model(&models.Person{}).Where("checked_out IS NULL").Count(&count)
	if err.Error != nil {
		return 0, err.Error
	}
	return count, nil
}