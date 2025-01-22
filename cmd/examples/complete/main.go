package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samuelarogbonlo/gorealconf/pkg/gorealconf"
)

// Custom Duration type for JSON unmarshaling
type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return fmt.Errorf("invalid duration")
	}
}

func (d Duration) TimeDuration() time.Duration {
	return time.Duration(d)
}

type CompleteConfig struct {
	Server struct {
		Port         int      `json:"port"`
		ReadTimeout  Duration `json:"read_timeout"`
		WriteTimeout Duration `json:"write_timeout"`
	} `json:"server"`
	Database struct {
		Host        string   `json:"host"`
		Port        int      `json:"port"`
		MaxConns    int      `json:"max_connections"`
		IdleTimeout Duration `json:"idle_timeout"`
	} `json:"database"`
	Features struct {
		EnableNewUI    bool    `json:"enable_new_ui"`
		EnableBetaAPI  bool    `json:"enable_beta_api"`
		RolloutPercent float64 `json:"rollout_percent"`
	} `json:"features"`
}

func main() {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting application...")

	// Create file source
	fileSource, err := gorealconf.NewFileSource[CompleteConfig]("cmd/examples/complete/config.json")
	if err != nil {
		log.Fatal("Failed to create file source:", err)
	}
	log.Println("File source created successfully")

	// Create etcd source
	etcdSource, err := gorealconf.NewEtcdSource[CompleteConfig](
		[]string{"localhost:2379"},
		"/app/config",
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Etcd source created successfully")

	// Create metrics
	metrics := gorealconf.NewMetrics("myapp")
	log.Println("Metrics initialized")

	// Create configuration with all features
	cfg := gorealconf.New[CompleteConfig](
		gorealconf.WithSource(fileSource),
		gorealconf.WithSource(etcdSource),
		gorealconf.WithValidation(validateConfig),
		gorealconf.WithRollback[CompleteConfig](true),
		gorealconf.WithMetrics[CompleteConfig](metrics),
	)

	// Load initial configuration
	if err := cfg.Load(ctx); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	log.Printf("Initial configuration loaded: %+v", cfg.Get(ctx))

	// Simulate configuration changes
	go func() {
		time.Sleep(2 * time.Second)
		newConfig := CompleteConfig{
			Server: struct {
				Port         int      `json:"port"`
				ReadTimeout  Duration `json:"read_timeout"`
				WriteTimeout Duration `json:"write_timeout"`
			}{
				Port:         8081,
				ReadTimeout:  Duration(10 * time.Second),
				WriteTimeout: Duration(15 * time.Second),
			},
			Database: struct {
				Host        string   `json:"host"`
				Port        int      `json:"port"`
				MaxConns    int      `json:"max_connections"`
				IdleTimeout Duration `json:"idle_timeout"`
			}{
				Host:        "localhost",
				Port:        5432,
				MaxConns:    100,
				IdleTimeout: Duration(5 * time.Minute),
			},
			Features: struct {
				EnableNewUI    bool    `json:"enable_new_ui"`
				EnableBetaAPI  bool    `json:"enable_beta_api"`
				RolloutPercent float64 `json:"rollout_percent"`
			}{
				EnableNewUI:    false,
				EnableBetaAPI:  false,
				RolloutPercent: 10.0,
			},
		}
		log.Printf("Simulating config update...")
		if err := cfg.Update(ctx, newConfig); err != nil {
			log.Printf("Error updating config: %v", err)
		}
	}()

	// Initialize HTTP server
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Get(ctx).Server.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Server is running!")
		}),
	}

	// Create error channel
	errChan := make(chan error, 1)

	// Start HTTP server
	go func() {
		log.Printf("Starting server on port %d", cfg.Get(ctx).Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
			log.Printf("Server error: %v", err)
		}
	}()
	// Wait for interrupt signal or server shutdown
	select {
	case sig := <-sigChan:
		log.Printf("\nReceived signal %v, shutting down...", sig)

		// Shutdown the server
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		cancel()
		return // Exit cleanly instead of os.Exit

	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
		return
	}
}

func validateConfig(old, new CompleteConfig) error {
	// Validate server settings
	if new.Server.Port < 1024 || new.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d", new.Server.Port)
	}
	if new.Server.ReadTimeout.TimeDuration() < time.Second {
		return fmt.Errorf("read timeout too short")
	}
	if new.Server.WriteTimeout.TimeDuration() < time.Second {
		return fmt.Errorf("write timeout too short")
	}

	// Validate database settings
	if new.Database.MaxConns < 1 {
		return fmt.Errorf("max connections must be positive")
	}
	if new.Database.IdleTimeout.TimeDuration() < time.Second {
		return fmt.Errorf("idle timeout too short")
	}

	// Validate feature rollout
	if new.Features.RolloutPercent < 0 || new.Features.RolloutPercent > 100 {
		return fmt.Errorf("rollout percentage must be between 0 and 100")
	}

	return nil
}

func applyConfig(cfg CompleteConfig) error {
	// Apply configuration changes
	log.Printf("Applying configuration: Server Port=%d, Features: UI=%v, API=%v",
		cfg.Server.Port,
		cfg.Features.EnableNewUI,
		cfg.Features.EnableBetaAPI)
	return nil
}
