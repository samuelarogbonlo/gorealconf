package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/samuelarogbonlo/dynconf/pkg/dynconf"
)

type DBConfig struct {
    Host         string        `json:"host"`
    Port         int           `json:"port"`
    User         string        `json:"user"`
    Password     string        `json:"password"`
    Database     string        `json:"database"`
    MaxOpenConns int           `json:"max_open_conns"`
    MaxIdleConns int           `json:"max_idle_conns"`
    ConnMaxLife  time.Duration `json:"conn_max_life"`
}

func main() {
    cfg := dynconf.New[DBConfig](
        dynconf.WithValidation[DBConfig](validateDBConfig),
        dynconf.WithRollback[DBConfig](true),
    )

    // Watch for config changes
    changes, _ := cfg.Watch(context.Background())
    go func() {
        for newCfg := range changes {
            log.Printf("Database config updated: %+v", newCfg)
            updateDBConnection(newCfg)
        }
    }()

    // Initial configuration
    if err := cfg.Update(context.Background(), DBConfig{
        Host:         "localhost",
        Port:         5432,
        MaxOpenConns: 25,
        MaxIdleConns: 5,
        ConnMaxLife:  5 * time.Minute,
    }); err != nil {
        log.Fatal(err)
    }

    select {} // Keep running
}

func updateDBConnection(cfg DBConfig) {
    // Implementation to update database connection pool
}

func validateDBConfig(old, new DBConfig) error {
    if new.MaxOpenConns < new.MaxIdleConns {
        return fmt.Errorf("max_open_conns must be >= max_idle_conns")
    }
    if new.ConnMaxLife < time.Second {
        return fmt.Errorf("conn_max_life must be at least 1 second")
    }
    return nil
}