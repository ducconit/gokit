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
	if claims.SubjectID != "u1" {
		t.Fatalf("expected subject u1, got %s", claims.SubjectID)
	}
	if claims.SubjectType != "user" {
		t.Fatalf("expected subject type 'user', got %s", claims.SubjectType)
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

func TestLogoutAll(t *testing.T) {
	ctx := context.Background()

	store := NewMemoryStore()
	m, _ := New(Config{
		Issuer:     "test",
		Audience:   []string{"test"},
		AccessTTL:  1 * time.Minute,
		RefreshTTL: 10 * time.Minute,
		HMACSecret: []byte("secret"),
	}, store)

	pair1, _ := m.Issue(ctx, "u1", "user", nil)
	pair2, _ := m.Issue(ctx, "u1", "user", nil)

	if err := m.LogoutAll(ctx, "u1"); err != nil {
		t.Fatalf("LogoutAll: %v", err)
	}

	if _, err := m.VerifyAccess(ctx, pair1.AccessToken); err == nil {
		t.Fatalf("expected unauthorized for pair1 after logout all")
	}
	if _, err := m.VerifyAccess(ctx, pair2.AccessToken); err == nil {
		t.Fatalf("expected unauthorized for pair2 after logout all")
	}
}

func TestBanSubject(t *testing.T) {
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

	if err := m.BanSubject(ctx, "u1", time.Now().Add(1*time.Hour)); err != nil {
		t.Fatalf("BanSubject: %v", err)
	}

	if _, err := m.VerifyAccess(ctx, pair.AccessToken); err == nil {
		t.Fatalf("expected forbidden when banned")
	}
}
