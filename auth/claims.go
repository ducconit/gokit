package auth

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	UserID    string `json:"uid"`
	UserType  string `json:"utp,omitempty"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID    string `json:"uid"`
	UserType  string `json:"utp,omitempty"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}
