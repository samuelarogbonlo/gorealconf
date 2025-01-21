package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/samuelarogbonlo/dynconf/pkg/dynconf"
)

// Define a comprehensive configuration structure
type CompleteConfig struct {
    // Server settings
    Server struct {
        Port         int           `json:"port"`
        ReadTimeout  time.Duration `json:"read_timeout"`
        WriteTimeout time.Duration `json:"write_timeout"`
    } `json:"server"`

    // Database settings
    Database struct {
        Host         string        `json:"host"`
        Port         int          `json:"port"`
        MaxConns     int          `json:"max_connections"`
        IdleTimeout  time.Duration `json:"idle_timeout"`
    } `json:"database"`

    // Feature flags
    Features struct {
        EnableNewUI     bool    `json:"enable_new_ui"`
        EnableBetaAPI   bool    `json:"enable_beta_api"`
        RolloutPercent  float64 `json:"rollout_percent"`
    } `json:"features"`
}

func main() {
    ctx := context.Background()

    // Create file source
    fileSource, err := dynconf.NewFileSource[CompleteConfig]("config.json")
    if err != nil {
        log.Fatal(err)
    }

    // Create etcd source
    etcdSource, err := dynconf.NewEtcdSource[CompleteConfig](
        []string{"localhost:2379"},
        "/app/config",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create metrics
    metrics := dynconf.NewMetrics("myapp")

    // Create configuration with all features
    cfg := dynconf.New[CompleteConfig](
        dynconf.WithSource(fileSource),
        dynconf.WithSource(etcdSource),
        dynconf.WithValidation(validateConfig),
        dynconf.WithRollback[CompleteConfig](true),
        dynconf.WithMetrics[CompleteConfig](metrics),  // Specify the type parameter here
    )

    // Load initial configuration
    if err := cfg.Load(ctx); err != nil {
        log.Fatal("Failed to load configuration:", err)
    }

    // Subscribe to configuration updates
    updates, unsubscribe := cfg.Subscribe(ctx)
    defer unsubscribe()

    // Handle configuration updates
    go func() {
        for newCfg := range updates {
            log.Printf("Applying new configuration: %+v", newCfg)
            if err := applyConfig(newCfg); err != nil {
                log.Printf("Error applying config: %v", err)
            }
        }
    }()

    // Initialize HTTP server
    server := &http.Server{
        Addr: fmt.Sprintf(":%d", cfg.Get(ctx).Server.Port),
    }

    // Start HTTP server
    go func() {
        log.Printf("Starting server on port %d", cfg.Get(ctx).Server.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Keep application running
    select {}
}

// Configuration validation
func validateConfig(old, new CompleteConfig) error {
    // Validate server settings
    if new.Server.Port < 1024 || new.Server.Port > 65535 {
        return fmt.Errorf("invalid port: %d", new.Server.Port)
    }
    if new.Server.ReadTimeout < time.Second {
        return fmt.Errorf("read timeout too short")
    }
    if new.Server.WriteTimeout < time.Second {
        return fmt.Errorf("write timeout too short")
    }

    // Validate database settings
    if new.Database.MaxConns < 1 {
        return fmt.Errorf("max connections must be positive")
    }
    if new.Database.IdleTimeout < time.Second {
        return fmt.Errorf("idle timeout too short")
    }

    // Validate feature rollout
    if new.Features.RolloutPercent < 0 || new.Features.RolloutPercent > 100 {
        return fmt.Errorf("rollout percentage must be between 0 and 100")
    }

    return nil
}

// Apply configuration changes
func applyConfig(cfg CompleteConfig) error {
    // Apply server changes
    // Apply database changes
    // Update feature flags
    return nil
}