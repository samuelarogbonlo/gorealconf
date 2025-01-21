package main

import (
    "context"
    "errors"
    "log"
    "time"

    "github.com/samuelarogbonlo/dynconf/pkg/dynconf"
)

type DatabaseConfig struct {
    Host        string        `json:"host"`
    Port        int           `json:"port"`
    MaxConns    int           `json:"max_connections"`
    IdleTimeout time.Duration `json:"idle_timeout"`
}

func main() {
    ctx := context.Background()

    // Create file source
    fileSource, err := dynconf.NewFileSource[DatabaseConfig]("config.json")
    if err != nil {
        log.Fatal(err)
    }

    // Create etcd source
    etcdSource, err := dynconf.NewEtcdSource[DatabaseConfig](
        []string{"localhost:2379"},
        "/app/database",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create config with sources
    cfg := dynconf.New[DatabaseConfig](
        // Try the basic Option approach
        dynconf.Option[DatabaseConfig](func(c *dynconf.Config[DatabaseConfig]) {
            c.AddSource(fileSource)
            c.AddSource(etcdSource)
        }),
        dynconf.WithValidation[DatabaseConfig](validateDatabase),
        dynconf.WithRollback[DatabaseConfig](true),
    )

    changes, cleanup := cfg.Subscribe(ctx)
    defer cleanup()

    go func() {
        for newCfg := range changes {
            log.Printf("Database config updated: %+v", newCfg)
            if err := updateDatabase(newCfg); err != nil {
                log.Printf("Error updating database: %v", err)
            }
        }
    }()

    select {}
}

func validateDatabase(old, new DatabaseConfig) error {
    if new.MaxConns <= 0 {
        return errors.New("max connections must be positive")
    }
    if new.IdleTimeout < time.Second {
        return errors.New("idle timeout must be at least 1 second")
    }
    return nil
}

func updateDatabase(cfg DatabaseConfig) error {
    // Implementation to update database connection pool
    return nil
}