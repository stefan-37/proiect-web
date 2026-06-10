package repository

import (
	"backend/models"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrClassFull is returned by BookClass when the class has reached its capacity.
var ErrClassFull = errors.New("class is full")

func BookClass(bookingSituation *models.BookingSituation, database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		var class models.Class
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&class, bookingSituation.ClassID).Error; err != nil {
			return err
		}

		if class.Users >= class.Capacity {
			return ErrClassFull
		}

		if err := tx.Create(&bookingSituation).Error; err != nil {
			return err
		}

		return tx.Model(&class).Update("users", class.Users+1).Error
	})
}

func GetReservationsByUserID(userID uint, database *gorm.DB) ([]models.BookingSituation, error) {
	var reservations []models.BookingSituation
	err := database.Where("user_id = ?", userID).Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func DeleteReservationByUserIDAndClassID(userID uint, classID uint, database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		var class models.Class
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&class, classID).Error; err != nil {
			return err
		}

		var reservation models.BookingSituation
		if err := tx.Where("user_id = ? AND class_id = ?", userID, classID).First(&reservation).Error; err != nil {
			return err
		}

		if err := tx.Delete(&reservation).Error; err != nil {
			return err
		}

		if class.Users > 0 {
			return tx.Model(&class).Update("users", class.Users-1).Error
		}
		return nil
	})
}

