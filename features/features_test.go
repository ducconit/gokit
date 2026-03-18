package features

import (
	"context"
	"testing"
)

func TestManager_IsEnabled(t *testing.T) {
	store := NewMemoryStore()
	manager := NewManager(store)
	ctx := context.Background()

	// 1. Global Enabled Flag
	store.Set(ctx, &Feature{
		Key:     "global-feature",
		Enabled: true,
	})

	if !manager.IsEnabled(ctx, "global-feature", nil) {
		t.Errorf("Global feature should be enabled")
	}

	// 2. Global Disabled Flag
	store.Set(ctx, &Feature{
		Key:     "disabled-feature",
		Enabled: false,
	})

	if manager.IsEnabled(ctx, "disabled-feature", nil) {
		t.Errorf("Disabled feature should be disabled")
	}

	// 3. Rule-based Feature (userID)
	store.Set(ctx, &Feature{
		Key:     "user-feature",
		Enabled: true,
		Rules: []Rule{
			{
				Attribute: "userID",
				Operator:  OpIn,
				Values:    []any{"user-1", "user-2"},
			},
		},
	})

	if !manager.IsEnabled(ctx, "user-feature", &Context{UserID: "user-1"}) {
		t.Errorf("User feature should be enabled for user-1")
	}

	if manager.IsEnabled(ctx, "user-feature", &Context{UserID: "user-3"}) {
		t.Errorf("User feature should be disabled for user-3")
	}

	// 4. Attribute-based Feature (email)
	store.Set(ctx, &Feature{
		Key:     "email-feature",
		Enabled: true,
		Rules: []Rule{
			{
				Attribute: "email",
				Operator:  OpRegex,
				Values:    []any{`.*@example\.com$`},
			},
		},
	})

	if !manager.IsEnabled(ctx, "email-feature", &Context{Attributes: map[string]any{"email": "test@example.com"}}) {
		t.Errorf("Email feature should be enabled for example.com emails")
	}

	if manager.IsEnabled(ctx, "email-feature", &Context{Attributes: map[string]any{"email": "test@gmail.com"}}) {
		t.Errorf("Email feature should be disabled for gmail.com emails")
	}
}
