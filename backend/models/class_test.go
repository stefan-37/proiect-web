package models

import (
	"testing"
	"time"
)

func TestClassFactory_ValidClass(t *testing.T) {
	class, err := ClassFactory(
		ClassWithName("Yoga"),
		ClassWithTrainerID(1),
		ClassWithCapacity(20),
		ClassWithAdminID(1),
		ClassWithScheduledAt(time.Now().Add(24*time.Hour)),
		ClassWithDescription("Morning yoga session"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if class.Name != "Yoga" {
		t.Errorf("expected name 'Yoga', got '%s'", class.Name)
	}
	if class.TrainerID != 1 {
		t.Errorf("expected trainer_id 1, got %d", class.TrainerID)
	}
	if class.Capacity != 20 {
		t.Errorf("expected capacity 20, got %d", class.Capacity)
	}
	if class.Users != 0 {
		t.Errorf("expected users to default to 0, got %d", class.Users)
	}
}

func TestClassFactory_EmptyName(t *testing.T) {
	_, err := ClassFactory(
		ClassWithTrainerID(1),
		ClassWithCapacity(20),
	)

	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestClassFactory_ZeroTrainerID(t *testing.T) {
	_, err := ClassFactory(
		ClassWithName("Yoga"),
		ClassWithCapacity(20),
	)

	if err == nil {
		t.Fatal("expected error for zero trainer_id, got nil")
	}
}

func TestClassFactory_ZeroCapacity(t *testing.T) {
	_, err := ClassFactory(
		ClassWithName("Yoga"),
		ClassWithTrainerID(1),
	)

	if err == nil {
		t.Fatal("expected error for zero capacity, got nil")
	}
}

func TestClassFactory_UsersDefaultsToZero(t *testing.T) {
	class, err := ClassFactory(
		ClassWithName("Pilates"),
		ClassWithTrainerID(2),
		ClassWithCapacity(15),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if class.Users != 0 {
		t.Errorf("expected Users to be 0, got %d", class.Users)
	}
}
