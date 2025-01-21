package main

import (
    "context"
    "errors" // Added this import
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

    // You might need to use Subscribe() instead of Watch() depending on the package version
    changes, err := cfg.Subscribe(ctx) // Changed Watch to Subscribe
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        for newCfg := range changes {
            log.Printf("Config updated: %+v", newCfg)
        }
    }()

    // Application logic...
    select {}
}