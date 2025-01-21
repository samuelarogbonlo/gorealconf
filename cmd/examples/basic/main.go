package main

import (
    "context"
    "errors"
    "log"
    "time"

    "github.com/samuelarogbonlo/dynconf/pkg/dynconf"
)

type AppConfig struct {
    ServerPort int           `json:"server_port"`
    Timeout    time.Duration `json:"timeout"`
}

func main() {
    ctx := context.Background()

    cfg := dynconf.New[AppConfig](
        dynconf.WithValidation[AppConfig](func(old, new AppConfig) error {
            if new.ServerPort < 1024 {
                return errors.New("server port must be >= 1024")
            }
            return nil
        }),
        dynconf.WithRollback[AppConfig](true),
    )

    // Initialize with a default configuration
    initialConfig := AppConfig{
        ServerPort: 8080,
        Timeout:    5 * time.Second,
    }

    if err := cfg.Update(ctx, initialConfig); err != nil {
        log.Fatal(err)
    }

    // Subscribe to changes - get both the channel and cleanup function
    changes, cleanup := cfg.Subscribe(ctx)
    defer cleanup() // Ensure cleanup is called when main exits

    // Handle configuration updates in a separate goroutine
    go func() {
        for update := range changes {
            log.Printf("Config updated: %+v", update)
        }
    }()

    // Keep the application running
    select {}
}