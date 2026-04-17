package models

import (
	"fmt"
	"gorm.io/gorm"
)

type status uint

const (
	Booked status = iota
	Refused
)

type BookingSituation struct {
	gorm.Model
	UserId  uint   `json:"user_id" gorm:"not null"`
	ClassID uint   `json:"class_id" gorm:"not null"`
	AdminID uint   `json:"admin_id" gorm:"not null"`
	Status  status `json:"status" gorm:"not null"`
}

type BookingSituationOptions func(*BookingSituation)

func BookingSituationWithUserID(userID uint) BookingSituationOptions {
	return func(bs *BookingSituation) {
		bs.UserId = userID
	}
}

func BookingSituationWithClassID(classID uint) BookingSituationOptions {
	return func(bs *BookingSituation) {
		bs.ClassID = classID
	}
}

func BookingSituationWithAdminID(adminID uint) BookingSituationOptions{
	return func(bs *BookingSituation) {
		bs.AdminID = adminID
	}
}

func BookingSituationWithStatus(status status) BookingSituationOptions {
	return func(bs *BookingSituation) {
		bs.Status = status
	}	
}

func (bs *BookingSituation) BookingSituationBuild() error {
	if bs.UserId == 0 {
		return fmt.Errorf("invalid user id")
	}	
	if bs.ClassID == 0 {
		return fmt.Errorf("invalid class id")
	}
	if bs.AdminID == 0 {
		return fmt.Errorf("invalid admin id")
	}
	return nil
}

func BookingSituationFactory(options ...BookingSituationOptions) (*BookingSituation, error) {
	bs := &BookingSituation{}
	for _, option := range options {
		option(bs)
	}
	if err := bs.BookingSituationBuild(); err != nil {
		return nil, err
	}
	return bs, nil
}
