package models

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	PersonID   uint      `json:"person_id" gorm:"not null"`
	CheckedIn  time.Time `json:"checked_in"`
	CheckedOut *time.Time `json:"checked_out"`
}

func NewPerson(personID uint, checkedIn time.Time, checkedOut *time.Time) *Person {
	return &Person{
		PersonID:   personID,
		CheckedIn:  checkedIn,
		CheckedOut: checkedOut,
	}
}
