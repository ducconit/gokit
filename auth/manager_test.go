package auth

import (
	"context"
	"testing"
	"time"
)

func TestIssueVerifyRefreshLogout(t *testing.T) {
	ctx := context.Background()

	store := NewMemoryStore()
	m, err := New(Config{
		Issuer:     "test",
		Audience:   []string{"test"},
		AccessTTL:  1 * time.Minute,
		RefreshTTL: 10 * time.Minute,
		HMACSecret: []byte("secret"),
	}, store)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	pair, err := m.Issue(ctx, "u1", "user", nil)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	if err := m.UpdateSessionMetadata(ctx, pair.SessionID, map[string]string{"device": "test"}); err != nil {
		t.Fatalf("UpdateSessionMetadata: %v", err)
	}

	claims, err := m.VerifyAccess(ctx, pair.AccessToken)
	if err != nil {
		t.Fatalf("VerifyAccess: %v", err)
	}
	if claims.UserID != "u1" {
		t.Fatalf("expected user u1, got %s", claims.UserID)
	}
	if claims.UserType != "user" {
		t.Fatalf("expected user type 'user', got %s", claims.UserType)
	}

	newPair, err := m.Refresh(ctx, pair.RefreshToken, RefreshOptions{Rotate: true})
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if newPair.RefreshToken == pair.RefreshToken {
		t.Fatalf("expected refresh token rotated")
	}

	// Test Refresh without rotation
	noRotatePair, err := m.Refresh(ctx, newPair.RefreshToken, RefreshOptions{Rotate: false})
	if err != nil {
		t.Fatalf("Refresh (no rotate): %v", err)
	}
	if noRotatePair.RefreshToken != newPair.RefreshToken {
		t.Fatalf("expected refresh token NOT rotated")
	}

	if err := m.Logout(ctx, pair.SessionID); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if _, err := m.VerifyAccess(ctx, newPair.AccessToken); err == nil {
		t.Fatalf("expected unauthorized after logout")
	}
}

func TestBanUser(t *testing.T) {
	ctx := context.Background()

	store := NewMemoryStore()
	m, err := New(Config{
		Issuer:     "test",
		Audience:   []string{"test"},
		AccessTTL:  1 * time.Minute,
		RefreshTTL: 10 * time.Minute,
		HMACSecret: []byte("secret"),
	}, store)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	pair, err := m.Issue(ctx, "u1", "user", nil)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	if err := m.BanUser(ctx, "u1", time.Now().Add(1*time.Hour)); err != nil {
		t.Fatalf("BanUser: %v", err)
	}

	if _, err := m.VerifyAccess(ctx, pair.AccessToken); err == nil {
		t.Fatalf("expected forbidden when banned")
	}
}
