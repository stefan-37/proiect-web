package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Trainer struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null"`
	Email       string `json:"email" gorm:"unique not null"`
	Password    string `json:"password" gorm:"not null"`
	Description string `json:"description"`
	AdminID     uint   `json:"admin_id" gorm:"not null"`
}

type TrainerOption func(*Trainer)

func TrainerWithName(name string) TrainerOption {
	return func(t *Trainer) {
		t.Name = name
	}
}

func TrainerWithEmail(email string) TrainerOption {
	return func(t *Trainer) {
		t.Email = email
	}
}

func TrainerWithPassword(password string) TrainerOption {
	return func(t *Trainer) {
		t.Password = password
	}
}

func TrainerWithDescription(description string) TrainerOption {
	return func(t *Trainer) {
		t.Description = description
	}
}

func TrainerWithAdminID(adminID uint) TrainerOption {
	return func(t *Trainer) {
		t.AdminID = adminID
	}
}

func (t *Trainer) TrainerBuild() error {
	if t.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if t.Email == "" {
		return fmt.Errorf("invalid email")
	}
	if t.Password == "" {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func TrainerFactory(options ...TrainerOption) (*Trainer, error) {
	trainer := &Trainer{}

	for _, option := range options {
		option(trainer)
	}

	if err := trainer.TrainerBuild(); err != nil {
		return nil, err
	}

	return trainer, nil
}
