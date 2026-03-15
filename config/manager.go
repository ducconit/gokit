package config

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Manager struct {
	v    *viper.Viper
	stop func()

	mu sync.RWMutex
}

func Load(ctx context.Context, src Source, opts Options) (*Manager, error) {
	if opts.ReloadInterval <= 0 {
		opts.ReloadInterval = 2 * time.Second
	}

	v := viper.New()

	m := &Manager{
		v: v,
	}

	if lfs, ok := src.(*LocalFileSource); ok {
		if err := m.loadFromLocalFile(ctx, lfs, opts); err != nil {
			return nil, err
		}
		return m, nil
	}

	if err := m.loadFromBytesSource(ctx, src, opts); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) Viper() *viper.Viper {
	return m.v
}

func (m *Manager) Unmarshal(dst any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.Unmarshal(dst)
}

func (m *Manager) Close() error {
	if m.stop != nil {
		m.stop()
	}
	return nil
}

func (m *Manager) loadFromLocalFile(ctx context.Context, src *LocalFileSource, opts Options) error {
	if src.Path == "" {
		return fmt.Errorf("config: missing local file path")
	}

	m.v.SetConfigFile(src.Path)
	if opts.Format != "" {
		m.v.SetConfigType(opts.Format)
	}
	if err := m.v.ReadInConfig(); err != nil {
		return err
	}

	if !opts.AutoReload {
		return nil
	}

	m.v.OnConfigChange(func(e fsnotify.Event) {
		_ = e
	})
	m.v.WatchConfig()
	return nil
}

func (m *Manager) loadFromBytesSource(ctx context.Context, src Source, opts Options) error {
	if opts.Format == "" {
		return fmt.Errorf("config: missing format (yaml/json/toml/ini/env)")
	}

	b, err := src.Load(ctx)
	if err != nil {
		return err
	}

	m.v.SetConfigType(opts.Format)
	if err := m.v.ReadConfig(bytes.NewReader(b)); err != nil {
		return err
	}

	if !opts.AutoReload {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	m.stop = cancel

	sum := sha256.Sum256(b)
	ticker := time.NewTicker(opts.ReloadInterval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				nb, nerr := src.Load(ctx)
				if nerr != nil {
					continue
				}

				nsum := sha256.Sum256(nb)
				if nsum == sum {
					continue
				}

				m.mu.Lock()
				_ = m.v.ReadConfig(bytes.NewReader(nb))
				m.mu.Unlock()

				sum = nsum
			}
		}
	}()

	return nil
}
