package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserSubscription struct {
	gorm.Model
	UserID         uint      `json:"user_id" gorm:"not null"`
	SubscriptionID uint      `json:"subscription_id" gorm:"not null"`
	StartedAt      time.Time `json:"started_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type UserSubscriptionOption func(*UserSubscription)

func UserSubscriptionWithUserID(userID uint) UserSubscriptionOption {
	return func(us *UserSubscription) {
		us.UserID = userID
	}
}

func UserSubscriptionWithSubscriptionID(subscriptionID uint) UserSubscriptionOption {
	return func(us *UserSubscription) {
		us.SubscriptionID = subscriptionID
	}
}

func UserSubscriptionWithStartedAt(startedAt time.Time) UserSubscriptionOption {
	return func(us *UserSubscription) {
		us.StartedAt = startedAt
	}
}

func UserSubscriptionWithExpiresAt(expiresAt time.Time) UserSubscriptionOption {
	return func(us *UserSubscription) {
		us.ExpiresAt = expiresAt
	}
}


func (us *UserSubscription) UserSubscriptionBuild() error {
	if us.UserID == 0 {
		return fmt.Errorf("invalid user ID")
	}
	if us.SubscriptionID == 0 {
		return fmt.Errorf("invalid subscription ID")
	}
	return nil
}

func UserSubscriptionFactory(options ...UserSubscriptionOption) (*UserSubscription, error) {
	userSubscription := &UserSubscription{}

	for _, option := range options {
		option(userSubscription)
	}

	if err := userSubscription.UserSubscriptionBuild(); err != nil {
		return nil, err
	}

	return userSubscription, nil
}
