package config

import "time"

type Options struct {
	Format         string
	AutoReload     bool
	ReloadInterval time.Duration
}
