package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/samuelarogbonlo/dynconf/pkg/dynconf"
)

type ServerConfig struct {
    Port            int           `json:"port"`
    ReadTimeout     time.Duration `json:"read_timeout"`
    WriteTimeout    time.Duration `json:"write_timeout"`
    MaxHeaderBytes  int          `json:"max_header_bytes"`
}

func main() {
    // Create configuration
    cfg := dynconf.New[ServerConfig](
        dynconf.WithValidation[ServerConfig](validateServerConfig),
        dynconf.WithRollback[ServerConfig](true),
    )

    // Initial config
    if err := cfg.Update(context.Background(), ServerConfig{
        Port:           8080,
        ReadTimeout:    5 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }); err != nil {
        log.Fatal(err)
    }

    // Watch for config changes
    changes, _ := cfg.Watch(context.Background())
    go func() {
        for newCfg := range changes {
            log.Printf("Server config updated: %+v", newCfg)
            // In a real application, you might want to gracefully restart the server
        }
    }()

    // Start HTTP server
    server := &http.Server{
        Addr:           fmt.Sprintf(":%d", cfg.Get(context.Background()).Port),
        ReadTimeout:    cfg.Get(context.Background()).ReadTimeout,
        WriteTimeout:   cfg.Get(context.Background()).WriteTimeout,
        MaxHeaderBytes: cfg.Get(context.Background()).MaxHeaderBytes,
    }

    log.Printf("Server starting on port %d", cfg.Get(context.Background()).Port)
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

func validateServerConfig(old, new ServerConfig) error {
    if new.Port < 1024 || new.Port > 65535 {
        return fmt.Errorf("port must be between 1024 and 65535")
    }
    if new.ReadTimeout < time.Second {
        return fmt.Errorf("read timeout must be at least 1 second")
    }
    if new.WriteTimeout < time.Second {
        return fmt.Errorf("write timeout must be at least 1 second")
    }
    return nil
}