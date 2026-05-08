package models

import "testing"

func TestUserFactory_ValidUser(t *testing.T) {
	user, err := UserFactory(
		UserWithName("Alice"),
		UserWithEmail("alice@example.com"),
		UserWithPassword("secret"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected name 'Alice', got '%s'", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got '%s'", user.Email)
	}
	if user.Password != "secret" {
		t.Errorf("expected password 'secret', got '%s'", user.Password)
	}
}

func TestUserFactory_EmptyName(t *testing.T) {
	_, err := UserFactory(
		UserWithEmail("alice@example.com"),
		UserWithPassword("secret"),
	)

	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestUserFactory_EmptyEmail(t *testing.T) {
	_, err := UserFactory(
		UserWithName("Alice"),
		UserWithPassword("secret"),
	)

	if err == nil {
		t.Fatal("expected error for empty email, got nil")
	}
}

func TestUserFactory_EmptyPassword(t *testing.T) {
	_, err := UserFactory(
		UserWithName("Alice"),
		UserWithEmail("alice@example.com"),
	)

	if err == nil {
		t.Fatal("expected error for empty password, got nil")
	}
}

func TestUserFactory_NoOptions(t *testing.T) {
	_, err := UserFactory()

	if err == nil {
		t.Fatal("expected error for no options, got nil")
	}
}
