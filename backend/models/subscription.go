package models

import (
	"fmt"
	"gorm.io/gorm"
)

type SubscriptionType uint

const (
	Basic SubscriptionType = iota
	Premium
)

type Subscription struct {
	gorm.Model
	Type    SubscriptionType `json:"type" gorm:"not null"`
	Price   float64          `json:"price" gorm:"not null"`
	AdminID uint             `json:"admin_id" gorm:"not null"`
}

type SubscriptionOption func(*Subscription)

func SubscriptionWithType(subscriptionType SubscriptionType) SubscriptionOption{
	return func(s *Subscription) {
		s.Type = subscriptionType
	}
}
func SubscriptionWithPrice(price float64) SubscriptionOption{
	return func(s *Subscription) {
		s.Price = price
	}
}
func SubscriptionWithAdminID(adminID uint) SubscriptionOption{
	return func(s *Subscription) {
		s.AdminID = adminID
	}
}

func (s *Subscription) SubscriptionBuild() error {
	if s.Price < 0 {
		return fmt.Errorf("invalid price")
	}
	if s.AdminID == 0 {
		return fmt.Errorf("invalid admin ID")
	}
	return nil
}

func SubscriptionFactory(options ...SubscriptionOption) (*Subscription, error) {
	subscription := &Subscription{}
	for _, option := range options {
		option(subscription)
	}
	if err := subscription.SubscriptionBuild(); err != nil {
		return nil, err
	}	
	return subscription, nil
}

