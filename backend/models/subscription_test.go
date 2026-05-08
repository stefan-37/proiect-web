package models

import "testing"

func TestSubscriptionFactory_ValidBasic(t *testing.T) {
	sub, err := SubscriptionFactory(
		SubscriptionWithType(Basic),
		SubscriptionWithPrice(29.99),
		SubscriptionWithAdminID(1),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sub.Type != Basic {
		t.Errorf("expected type Basic, got %v", sub.Type)
	}
	if sub.Price != 29.99 {
		t.Errorf("expected price 29.99, got %v", sub.Price)
	}
	if sub.AdminID != 1 {
		t.Errorf("expected admin_id 1, got %d", sub.AdminID)
	}
}

func TestSubscriptionFactory_ValidPremium(t *testing.T) {
	sub, err := SubscriptionFactory(
		SubscriptionWithType(Premium),
		SubscriptionWithPrice(59.99),
		SubscriptionWithAdminID(1),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sub.Type != Premium {
		t.Errorf("expected type Premium, got %v", sub.Type)
	}
}

func TestSubscriptionFactory_ZeroPriceIsValid(t *testing.T) {
	_, err := SubscriptionFactory(
		SubscriptionWithType(Basic),
		SubscriptionWithPrice(0),
		SubscriptionWithAdminID(1),
	)

	if err != nil {
		t.Fatalf("expected zero price to be valid, got %v", err)
	}
}

func TestSubscriptionFactory_NegativePrice(t *testing.T) {
	_, err := SubscriptionFactory(
		SubscriptionWithType(Basic),
		SubscriptionWithPrice(-1.0),
		SubscriptionWithAdminID(1),
	)

	if err == nil {
		t.Fatal("expected error for negative price, got nil")
	}
}

func TestSubscriptionFactory_ZeroAdminID(t *testing.T) {
	_, err := SubscriptionFactory(
		SubscriptionWithType(Basic),
		SubscriptionWithPrice(29.99),
	)

	if err == nil {
		t.Fatal("expected error for zero admin_id, got nil")
	}
}
