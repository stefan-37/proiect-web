package models

import (
	"testing"
	"time"
)

func TestUserSubscriptionFactory_Valid(t *testing.T) {
	start := time.Now()
	expires := start.AddDate(0, 1, 0)

	us, err := UserSubscriptionFactory(
		UserSubscriptionWithUserID(1),
		UserSubscriptionWithSubscriptionID(2),
		UserSubscriptionWithStartedAt(start),
		UserSubscriptionWithExpiresAt(expires),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if us.UserID != 1 {
		t.Errorf("expected user_id 1, got %d", us.UserID)
	}
	if us.SubscriptionID != 2 {
		t.Errorf("expected subscription_id 2, got %d", us.SubscriptionID)
	}
	if !us.StartedAt.Equal(start) {
		t.Errorf("expected started_at %v, got %v", start, us.StartedAt)
	}
	if !us.ExpiresAt.Equal(expires) {
		t.Errorf("expected expires_at %v, got %v", expires, us.ExpiresAt)
	}
}

func TestUserSubscriptionFactory_ZeroUserID(t *testing.T) {
	_, err := UserSubscriptionFactory(
		UserSubscriptionWithSubscriptionID(2),
		UserSubscriptionWithStartedAt(time.Now()),
		UserSubscriptionWithExpiresAt(time.Now().AddDate(0, 1, 0)),
	)

	if err == nil {
		t.Fatal("expected error for zero user_id, got nil")
	}
}

func TestUserSubscriptionFactory_ZeroSubscriptionID(t *testing.T) {
	_, err := UserSubscriptionFactory(
		UserSubscriptionWithUserID(1),
		UserSubscriptionWithStartedAt(time.Now()),
		UserSubscriptionWithExpiresAt(time.Now().AddDate(0, 1, 0)),
	)

	if err == nil {
		t.Fatal("expected error for zero subscription_id, got nil")
	}
}

func TestUserSubscriptionFactory_DatesAreOptional(t *testing.T) {
	_, err := UserSubscriptionFactory(
		UserSubscriptionWithUserID(1),
		UserSubscriptionWithSubscriptionID(2),
	)

	if err != nil {
		t.Fatalf("expected no error without dates, got %v", err)
	}
}
