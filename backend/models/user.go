package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique not null"`
	Password string `json:"password" gorm:"not null"`
}

type UserOption func(*User)

func UserWithName(name string) UserOption {
	return func(u *User) {
		u.Name = name
	}
}

func UserWithEmail(email string) UserOption {
	return func(u *User) {
		u.Email = email
	}
}

func UserWithPassword(password string) UserOption {
	return func(u *User) {
		u.Password = password
	}
}

func (u *User) UserBuild() error {
	if u.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if u.Email == "" {
		return fmt.Errorf("invalid email")
	}
	if u.Password == "" {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func UserFactory(options ...UserOption) (*User, error) {
	user := &User{}

	for _, option := range options {
		option(user)
	}

	if err := user.UserBuild(); err != nil {
		return nil, err
	}

	return user, nil
}
