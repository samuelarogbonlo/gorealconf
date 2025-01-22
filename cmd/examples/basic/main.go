package main

import (
    "context"
    "errors"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/samuelarogbonlo/dynconf/pkg/dynconf"
)

type AppConfig struct {
    ServerPort int           `json:"server_port"`
    Timeout    time.Duration `json:"timeout"`
}

func main() {
    // Create a cancellable context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    cfg := dynconf.New[AppConfig](
        dynconf.WithValidation[AppConfig](func(old, new AppConfig) error {
            if new.ServerPort < 1024 {
                return errors.New("server port must be >= 1024")
            }
            return nil
        }),
        dynconf.WithRollback[AppConfig](true),
    )

    // Initialize configuration
    initialConfig := AppConfig{
        ServerPort: 8080,
        Timeout:    5 * time.Second,
    }

    if err := cfg.Update(ctx, initialConfig); err != nil {
        log.Fatal(err)
    }

    // Setup signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Watch for changes
    changes, err := cfg.Watch(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Handle configuration updates
    go func() {
        for {
            select {
            case newCfg := <-changes:
                log.Printf("Config updated: %+v", newCfg)
            case <-ctx.Done():
                return
            }
        }
    }()

    // Wait for shutdown signal
    sig := <-sigChan
    log.Printf("Received signal %v, shutting down", sig)
}