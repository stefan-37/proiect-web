package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	ScheduledAt time.Time `json:"date" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	TrainerID   uint      `json:"trainer_id" gorm:"not null"`
	Capacity    uint      `json:"capacity" gorm:"not null"`
	Users       uint      `json:"users" gorm:"not null"`
	AdminID     uint      `json:"admin_id"`
}

type ClassOption func(*Class)

func ClassWithName(name string) ClassOption {
	return func(c *Class) {
		c.Name = name
	}
}

func ClassWithScheduledAt(scheduledAt time.Time) ClassOption {
	return func(c *Class) {
		c.ScheduledAt = scheduledAt
	}
}

func ClassWithDescription(description string) ClassOption {
	return func(c *Class) {
		c.Description = description
	}
}

func ClassWithTrainerID(trainerID uint) ClassOption {
	return func(c *Class) {
		c.TrainerID = trainerID
	}
}

func ClassWithCapacity(capacity uint) ClassOption {
	return func(c *Class) {
		c.Capacity = capacity
	}
}

func ClassWithAdminID(adminID uint) ClassOption {
	return func(c *Class) {
		c.AdminID = adminID
	}
}

func ClassWithPrice(price float64) ClassOption {
	return func(c *Class) {
		c.Price = price
	}
}

func (c *Class) ClassBuild() error {
	if c.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if c.TrainerID == 0 {
		return fmt.Errorf("invalid trainer_id")
	}
	if c.Capacity == 0 {
		return fmt.Errorf("invalid capacity")
	}
	if c.Price == 0 {
		return fmt.Errorf("invalid price")
	}
	return nil
}

func ClassFactory(options ...ClassOption) (*Class, error) {
	class := &Class{Users: 0}

	for _, option := range options {
		option(class)
	}

	if err := class.ClassBuild(); err != nil {
		return nil, err
	}
	return class, nil
}
