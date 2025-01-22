package main

import (
	"context"
	"errors" // Added missing import
	"log"

	"github.com/samuelarogbonlo/gorealconf/pkg/gorealconf"
)

type FeatureConfig struct {
	Enabled     bool    `json:"enabled"`
	Percentage  float64 `json:"percentage"`
	MaxRequests int     `json:"max_requests"`
}

func main() {
	ctx := context.Background()

	cfg := gorealconf.New[FeatureConfig]()

	// Create validation function
	validateFunc := func(old, new FeatureConfig) error {
		return validateFeature(new)
	}

	// Configure with validation
	cfg = gorealconf.New[FeatureConfig](
		gorealconf.WithValidation[FeatureConfig](validateFunc),
		gorealconf.WithRollback[FeatureConfig](true),
	)

	// Subscribe to changes instead of Watch
	changes, cleanup := cfg.Subscribe(ctx)
	defer cleanup()

	// Create a rollout checker
	isEnabled := func(cfg FeatureConfig) bool {
		return cfg.Enabled && cfg.Percentage > 0
	}

	go func() {
		for newCfg := range changes {
			if isEnabled(newCfg) {
				log.Printf("Applying feature config: %+v", newCfg)
				if err := applyFeature(newCfg); err != nil {
					log.Printf("Error applying feature: %v", err)
				}
			}
		}
	}()

	select {}
}

func validateFeature(cfg FeatureConfig) error {
	if cfg.Percentage < 0 || cfg.Percentage > 100 {
		return errors.New("percentage must be between 0 and 100")
	}
	if cfg.MaxRequests < 0 {
		return errors.New("max requests must be non-negative")
	}
	return nil
}

func applyFeature(cfg FeatureConfig) error {
	// Implementation to apply feature configuration
	return nil
}
