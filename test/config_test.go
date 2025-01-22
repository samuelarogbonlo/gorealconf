// test/config_test.go
package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/samuelarogbonlo/gorealconf/pkg/gorealconf"
)

func TestConfigOperations(t *testing.T) {
	tests := []struct {
		name    string
		config  ServerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ServerConfig{
				Port:    8080,
				Host:    "localhost",
				Timeout: time.Duration(5 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: ServerConfig{
				Port:    80,
				Host:    "localhost",
				Timeout: time.Duration(5 * time.Second),
			},
			wantErr: true,
		},
		{
			name: "zero timeout",
			config: ServerConfig{
				Port:    8080,
				Host:    "localhost",
				Timeout: time.Duration(0),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := gorealconf.New[ServerConfig](
				gorealconf.WithValidation[ServerConfig](func(old, new ServerConfig) error {
					if new.Port < 1024 {
						return fmt.Errorf("port must be >= 1024")
					}
					return nil
				}),
				gorealconf.WithRollback[ServerConfig](true),
			)

			err := cfg.Update(ctx, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got := cfg.Get(ctx)
				if got.Port != tt.config.Port {
					t.Errorf("Get() got = %v, want %v", got.Port, tt.config.Port)
				}
			}
		})
	}
}

func TestConcurrentOperations(t *testing.T) {
	ctx := context.Background()
	cfg := gorealconf.New[ServerConfig]()

	// Run multiple goroutines updating config
	const numGoroutines = 10
	errCh := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(port int) {
			config := ServerConfig{
				Port:    8080 + port,
				Host:    "localhost",
				Timeout: time.Duration(port) * time.Second,
			}
			errCh <- cfg.Update(ctx, config)
		}(i)
	}

	// Check for errors
	for i := 0; i < numGoroutines; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("Concurrent update failed: %v", err)
		}
	}

	close(errCh)
}
