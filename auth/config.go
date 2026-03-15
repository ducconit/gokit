package auth

import "time"

type Config struct {
	Issuer     string
	Audience   []string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	HMACSecret []byte
}
