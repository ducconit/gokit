package pagination

import (
	"net/url"
	"testing"
)

func TestBindQuerySimpleDefaults(t *testing.T) {
	var p Simple
	if err := BindQuery(url.Values{}, &p); err != nil {
		t.Fatalf("BindQuery: %v", err)
	}
	// p.Normalize() is already called inside BindQuery
	if p.Page != 1 {
		t.Fatalf("expected page=1, got %d", p.Page)
	}
	if p.Limit != 20 {
		t.Fatalf("expected limit=20, got %d", p.Limit)
	}
}

func TestBindQuerySimpleValues(t *testing.T) {
	v := url.Values{}
	v.Set("page", "3")
	v.Set("limit", "10")

	var p Simple
	if err := BindQuery(v, &p); err != nil {
		t.Fatalf("BindQuery: %v", err)
	}
	// p.Normalize() is already called inside BindQuery
	if p.Offset() != 20 {
		t.Fatalf("expected offset=20, got %d", p.Offset())
	}
}

func TestGenericNormalize(t *testing.T) {
	type CustomReq struct {
		Simple
		MaxLimit int `query:"limit" default:"50" max:"100"`
	}

	req := CustomReq{}
	// Test default from tag
	Normalize(&req)
	if req.Simple.Page != 1 {
		t.Fatalf("expected embedded page=1, got %d", req.Simple.Page)
	}
	if req.Simple.Limit != 20 {
		t.Fatalf("expected embedded limit=20, got %d", req.Simple.Limit)
	}
	if req.MaxLimit != 50 {
		t.Fatalf("expected MaxLimit=50, got %d", req.MaxLimit)
	}

	// Test max constraint
	req.MaxLimit = 200
	Normalize(&req)
	if req.MaxLimit != 100 {
		t.Fatalf("expected MaxLimit=100 (capped), got %d", req.MaxLimit)
	}
}
