package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique not null"`
	Password string `json:"password" gorm:"not null"`
}

type AdminOption func(*Admin)

func AdminWithName(name string) AdminOption {
	return func(a *Admin) {
		a.Name = name
	}
}

func AdminWithEmail(email string) AdminOption {
	return func(a *Admin) {
		a.Email = email
	}
}

func AdminWithPassword(password string) AdminOption {
	return func(a *Admin) {
		a.Password = password
	}
}

func (a *Admin) AdminBuild() error {
	if a.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if a.Email == "" {
		return fmt.Errorf("invalid email")
	}
	if a.Password == "" {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func AdminFactory(options ...AdminOption) (*Admin, error) {
	admin := &Admin{}

	for _, option := range options {
		option(admin)
	}

	if err := admin.AdminBuild(); err != nil {
		return nil, err
	}

	return admin, nil
}
