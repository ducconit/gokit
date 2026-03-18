package features

import "context"

// Context is the context of the user/request to evaluate the feature flag.
type Context struct {
	UserID     string
	Attributes map[string]any
}

// Feature represents a feature flag definition.
type Feature struct {
	Key         string `json:"key"`
	Enabled     bool   `json:"enabled"`
	Rules       []Rule `json:"rules,omitempty"`
	Description string `json:"description,omitempty"`
}

// Rule represents a condition to enable the feature for specific users/contexts.
type Rule struct {
	Attribute string   `json:"attribute"` // e.g., "userID", "email", "role", "group"
	Operator  Operator `json:"operator"`  // e.g., "in", "not_in", "regex", "percentage"
	Values    []any    `json:"values"`
}

type Operator string

const (
	OpIn         Operator = "in"
	OpNotIn      Operator = "not_in"
	OpRegex      Operator = "regex"
	OpPercentage Operator = "percentage"
)

// Store is the interface for feature flag storage.
type Store interface {
	Get(ctx context.Context, key string) (*Feature, error)
	All(ctx context.Context) (map[string]*Feature, error)
	Set(ctx context.Context, feature *Feature) error
}
