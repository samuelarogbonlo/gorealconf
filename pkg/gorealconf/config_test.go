package gorealconf

import (
	"context"
	"testing"
)

func TestConfig(t *testing.T) {
	type TestConfig struct {
		Value string
	}

	t.Run("basic operations", func(t *testing.T) {
		ctx := context.Background()
		cfg := New[TestConfig]()

		// Test initial state
		if got := cfg.Get(ctx); got.Value != "" {
			t.Errorf("expected empty initial value, got %v", got)
		}

		// Test update
		newCfg := TestConfig{Value: "test"}
		if err := cfg.Update(ctx, newCfg); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got := cfg.Get(ctx); got.Value != "test" {
			t.Errorf("expected value 'test', got %v", got)
		}
	})

	// Add more test cases...
}
