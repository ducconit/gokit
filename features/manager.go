package features

import (
	"context"
	"fmt"
	"regexp"
	"sync"
)

// Manager is the main feature flag evaluator.
type Manager struct {
	store Store
	mu    sync.RWMutex
}

// NewManager creates a new feature flag manager.
func NewManager(store Store) *Manager {
	return &Manager{
		store: store,
	}
}

// IsEnabled checks if a feature flag is enabled for the given context.
func (m *Manager) IsEnabled(ctx context.Context, key string, userCtx *Context) bool {
	feature, err := m.store.Get(ctx, key)
	if err != nil || feature == nil {
		return false
	}

	if !feature.Enabled {
		return false
	}

	// If no rules, feature is globally enabled.
	if len(feature.Rules) == 0 {
		return true
	}

	// Evaluate rules. If ANY rule matches, feature is enabled for this context.
	for _, rule := range feature.Rules {
		if m.evaluateRule(rule, userCtx) {
			return true
		}
	}

	return false
}

func (m *Manager) evaluateRule(rule Rule, userCtx *Context) bool {
	if userCtx == nil {
		return false
	}

	var val any
	if rule.Attribute == "userID" {
		val = userCtx.UserID
	} else if userCtx.Attributes != nil {
		val = userCtx.Attributes[rule.Attribute]
	}

	if val == nil {
		return false
	}

	switch rule.Operator {
	case OpIn:
		for _, v := range rule.Values {
			if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", val) {
				return true
			}
		}
	case OpNotIn:
		match := false
		for _, v := range rule.Values {
			if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", val) {
				match = true
				break
			}
		}
		return !match
	case OpRegex:
		strVal := fmt.Sprintf("%v", val)
		for _, v := range rule.Values {
			re, err := regexp.Compile(fmt.Sprintf("%v", v))
			if err == nil && re.MatchString(strVal) {
				return true
			}
		}
	case OpPercentage:
		// Basic percentage rollouts using UserID hash would be more robust.
		// For simplicity, this expects a percentage value (0-100) and we can use a simpler approach.
		// A full implementation should use a deterministic hash of UserID.
		// Let's skip complex implementation for now or use a simple hash.
	}

	return false
}
