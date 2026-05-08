package models

import "testing"

func TestTrainerFactory_ValidTrainer(t *testing.T) {
	trainer, err := TrainerFactory(
		TrainerWithName("Bob"),
		TrainerWithEmail("bob@gym.com"),
		TrainerWithPassword("pass123"),
		TrainerWithAdminID(1),
		TrainerWithDescription("Certified personal trainer"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if trainer.Name != "Bob" {
		t.Errorf("expected name 'Bob', got '%s'", trainer.Name)
	}
	if trainer.Email != "bob@gym.com" {
		t.Errorf("expected email 'bob@gym.com', got '%s'", trainer.Email)
	}
	if trainer.Description != "Certified personal trainer" {
		t.Errorf("expected description set, got '%s'", trainer.Description)
	}
}

func TestTrainerFactory_EmptyName(t *testing.T) {
	_, err := TrainerFactory(
		TrainerWithEmail("bob@gym.com"),
		TrainerWithPassword("pass123"),
	)

	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestTrainerFactory_EmptyEmail(t *testing.T) {
	_, err := TrainerFactory(
		TrainerWithName("Bob"),
		TrainerWithPassword("pass123"),
	)

	if err == nil {
		t.Fatal("expected error for empty email, got nil")
	}
}

func TestTrainerFactory_EmptyPassword(t *testing.T) {
	_, err := TrainerFactory(
		TrainerWithName("Bob"),
		TrainerWithEmail("bob@gym.com"),
	)

	if err == nil {
		t.Fatal("expected error for empty password, got nil")
	}
}

func TestTrainerFactory_DescriptionIsOptional(t *testing.T) {
	_, err := TrainerFactory(
		TrainerWithName("Bob"),
		TrainerWithEmail("bob@gym.com"),
		TrainerWithPassword("pass123"),
	)

	if err != nil {
		t.Fatalf("expected no error without description, got %v", err)
	}
}
