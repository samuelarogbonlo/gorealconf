# gorealconf

[![Go Reference](https://pkg.go.dev/badge/github.com/samuelarogbonlo/gorealconf.svg)](https://pkg.go.dev/github.com/samuelarogbonlo/gorealconf)
[![Build Status](https://github.com/samuelarogbonlo/gorealconf/workflows/CI/badge.svg)](https://github.com/samuelarogbonlo/gorealconf/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/samuelarogbonlo/gorealconf)](https://goreportcard.com/report/github.com/samuelarogbonlo/gorealconf)
[![codecov](https://codecov.io/gh/samuelarogbonlo/gorealconf/branch/main/graph/badge.svg)](https://codecov.io/gh/samuelarogbonlo/gorealconf)

gorealconf is a type-safe dynamic configuration management library for Go applications.

## Installation

Latest version:
```bash
go get github.com/samuelarogbonlo/gorealconf
```

Specific version:
```bash
go get github.com/samuelarogbonlo/gorealconf@v0.1.0
```

## Quick Start

1. Define your configuration:
```go
type AppConfig struct {
    ServerPort int           `json:"server_port"`
    Timeout    time.Duration `json:"timeout"`
}
```

2. Create and use configuration:
```go
cfg := gorealconf.New[AppConfig]()

// Add validation
cfg = gorealconf.New[AppConfig](
    gorealconf.WithValidation[AppConfig](func(old, new AppConfig) error {
        if new.ServerPort < 1024 {
            return errors.New("port must be >= 1024")
        }
        return nil
    }),
)

// Watch for changes
changes, _ := cfg.Watch(context.Background())
go func() {
    for newCfg := range changes {
        log.Printf("Config updated: %+v", newCfg)
    }
}()
```

## Key Features with Examples

### Type-safe Configuration
```go
type Config struct {
    Port    int           `json:"port"`
    Timeout time.Duration `json:"timeout"`
}
cfg := gorealconf.New[Config]()
```

### Multiple Configuration Sources
```go
cfg := gorealconf.New[Config](
    gorealconf.WithSource(fileSource),
    gorealconf.WithSource(etcdSource),
)
```

### Automatic Validation & Rollback
```go
cfg := gorealconf.New[Config](
    gorealconf.WithValidation[Config](validateConfig),
    gorealconf.WithRollback[Config](true),
)
```

### Gradual Rollouts
```go
strategy := gorealconf.NewPercentageStrategy(10)
rollout := gorealconf.NewRollout[Config](cfg).
    WithStrategy(strategy)
```

### Metrics Integration
```go
metrics := gorealconf.NewMetrics("myapp")
cfg.WithMetrics(metrics)
```

## Requirements
- Go 1.21 or higher
- Compatible with:
  - etcd v3.5+
  - Consul v1.12+
  - Redis v6+
- Prometheus for metrics (optional)

## Advanced Features
- Graceful shutdown support
- Health check integration
- Automatic rollback on validation failure
- Multiple source priority handling
- Custom duration parsing for JSON configs

## Examples
See [Examples](examples/) directory for:
- Basic usage (`examples/basic/`)
- Multi-source configuration (`examples/multistore/`)
- Gradual rollouts (`examples/rollout/`)
- Complete server example (`examples/complete/`)

## Documentation
- [Getting Started](docs/getting-started.md)
- [Configuration Sources](docs/configuration-sources.md)
- [Rollout Strategies](docs/rollout-strategies.md)
- [Metrics](docs/metrics.md)
- [FAQ](docs/faq.md)
- [Troubleshooting](docs/troubleshooting.md)

## Security
For security concerns, please email sbayo971@gmail.com

## Support
- Report bugs via GitHub Issues
- Join discussions in GitHub Discussions
- Get community help via GitHub Discussions

## Stability Status
Production-Ready: v0.1.0

## Contributing
See [Contributing Guide](docs/contributing.md) for guidelines.

## License
MIT License - see [LICENSE](LICENSE) for details.
```
