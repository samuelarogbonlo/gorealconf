// test/basic_test.go
package test

import (
    "context"
    "fmt"
    "os"
    "testing"
    "time"
    "encoding/json"

    "github.com/samuelarogbonlo/dynconf/pkg/dynconf"
)

type ServerConfig struct {
    Port    int           `json:"port"`
    Host    string        `json:"host"`
    Timeout time.Duration `json:"timeout"`
}

func (c *ServerConfig) UnmarshalJSON(data []byte) error {
    type Alias ServerConfig
    aux := &struct {
        Timeout string `json:"timeout"`
        *Alias
    }{
        Alias: (*Alias)(c),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    if aux.Timeout != "" {
        duration, err := time.ParseDuration(aux.Timeout)
        if err != nil {
            return err
        }
        c.Timeout = duration
    }
    return nil
}

func TestBasicConfiguration(t *testing.T) {
    ctx := context.Background()

    // Create configuration
    cfg := dynconf.New[ServerConfig](
        dynconf.WithValidation[ServerConfig](func(old, new ServerConfig) error {
            if new.Port < 1024 {
                return fmt.Errorf("port must be >= 1024")
            }
            return nil
        }),
        dynconf.WithRollback[ServerConfig](true),
    )

    // Test updating configuration
    testConfig := ServerConfig{
        Port:    8080,
        Host:    "localhost",
        Timeout: 5 * time.Second,
    }

    // Update config
    if err := cfg.Update(ctx, testConfig); err != nil {
        t.Errorf("Failed to update config: %v", err)
    }

    // Verify update
    current := cfg.Get(ctx)
    if current.Port != testConfig.Port {
        t.Errorf("Expected port %d, got %d", testConfig.Port, current.Port)
    }

    // Test validation
    invalidConfig := ServerConfig{
        Port: 80, // Should fail validation (< 1024)
        Host: "localhost",
    }

    if err := cfg.Update(ctx, invalidConfig); err == nil {
        t.Error("Expected validation error for invalid port, got nil")
    }
}

func TestFileSource(t *testing.T) {
    ctx := context.Background()

    // Create temporary config file
    configFile := createTempConfigFile(t)
    defer os.Remove(configFile)

    // Just use basic configuration
    cfg := dynconf.New[ServerConfig]()

    // Read and parse config file
    data, err := os.ReadFile(configFile)
    if err != nil {
        t.Fatalf("Failed to read config file: %v", err)
    }

    var config ServerConfig
    if err := json.Unmarshal(data, &config); err != nil {
        t.Fatalf("Failed to parse config: %v", err)
    }

    // Update config
    if err := cfg.Update(ctx, config); err != nil {
        t.Fatalf("Failed to update config: %v", err)
    }

    // Verify
    current := cfg.Get(ctx)
    if current.Port != 8080 {
        t.Errorf("Expected port 8080, got %d", current.Port)
    }
}

// Add the edge cases test here
func TestEdgeCases(t *testing.T) {
    t.Run("nil validator", func(t *testing.T) {
        ctx := context.Background()
        cfg := dynconf.New[ServerConfig]()

        config := ServerConfig{Port: 80} // Invalid port but no validator
        if err := cfg.Update(ctx, config); err != nil {
            t.Errorf("Expected no error with nil validator, got %v", err)
        }
    })

    t.Run("empty configuration", func(t *testing.T) {
        ctx := context.Background()
        cfg := dynconf.New[ServerConfig]()

        got := cfg.Get(ctx)
        if got.Port != 0 || got.Host != "" {
            t.Errorf("Expected zero value, got %+v", got)
        }
    })
}


// Helper function to create temporary config file
func createTempConfigFile(t *testing.T) string {
    content := `{
        "port": 8080,
        "host": "localhost",
        "timeout": "5s"
    }`

    tmpfile, err := os.CreateTemp("", "config-*.json")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }

    if _, err := tmpfile.Write([]byte(content)); err != nil {
        t.Fatalf("Failed to write to temp file: %v", err)
    }

    if err := tmpfile.Close(); err != nil {
        t.Fatalf("Failed to close temp file: %v", err)
    }

    return tmpfile.Name()
}