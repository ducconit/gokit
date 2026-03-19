package auth

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	SubjectID   string `json:"sid"`
	SubjectType string `json:"stp,omitempty"`
	SessionID   string `json:"ssi"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	SubjectID   string `json:"sid"`
	SubjectType string `json:"stp,omitempty"`
	SessionID   string `json:"ssi"`
	jwt.RegisteredClaims
}
