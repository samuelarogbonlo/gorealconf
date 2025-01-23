# gorealconf

[![Go Reference](https://pkg.go.dev/badge/github.com/samuelarogbonlo/gorealconf.svg)](https://pkg.go.dev/github.com/samuelarogbonlo/gorealconf)
[![Build Status](https://github.com/samuelarogbonlo/gorealconf/workflows/CI/badge.svg)](https://github.com/samuelarogbonlo/gorealconf/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/samuelarogbonlo/gorealconf)](https://goreportcard.com/report/github.com/samuelarogbonlo/gorealconf)
[![codecov](https://codecov.io/gh/samuelarogbonlo/gorealconf/branch/main/graph/badge.svg)](https://codecov.io/gh/samuelarogbonlo/gorealconf)

`gorealconf` is a powerful, type-safe dynamic configuration management library designed for modern Go applications. With support for real-time updates, multiple configuration sources, and gradual rollouts, it ensures zero-downtime deployments while maintaining operational stability.

## Table of Contents
- [Documentation](#documentation)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Key Features](#key-features-with-examples)
- [Requirements](#requirements)
- [Advanced Features](#advanced-features)
- [Examples](#examples)
- [Security](#security)
- [Support](#support)
- [Stability Status](#stability-status)
- [Contributing](#contributing)
- [License](#license)

## Documentation
- [Getting Started](docs/getting-started.md)
- [Configuration Sources](docs/configuration-sources.md)
- [Rollout Strategies](docs/rollout-strategies.md)
- [Metrics](docs/metrics.md)
- [FAQ](docs/faq.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Roadmap](docs/ROADMAP.md)

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

### Key Features
- **Type-safe Configuration**: Ensure your configurations are validated at compile time.
- **Multiple Configuration Sources**: Combine Redis, etcd, and files for flexibility.
- **Automatic Validation & Rollback**: Catch and rollback bad configurations automatically.
- **Gradual Rollouts**: Deploy new configurations incrementally.
- **Metrics Integration**: Export configuration change metrics via Prometheus.

## Examples

Explore the examples provided to understand how to use `gorealconf` in different scenarios:

- **[Basic Usage](examples/basic/main.go)**: Learn how to set up a simple configuration.
- **[Multi-source Configuration](examples/multistore/main.go)**: Use multiple configuration backends with priority handling.
- **[Webserver Integration](examples/webserver/main.go)**: Configure and manage dynamic settings for a webserver.
- **[Database Configuration](examples/database/main.go)**: Dynamically update and manage database configurations.
- **[Gradual Rollouts](examples/rollout/main.go)**: Deploy configuration updates incrementally and safely.
- **[Complete Setup](examples/complete/main.go)**: A fully functional example combining multiple sources, validation, and rollouts.

Each example includes inline comments and demonstrates real-world usage scenarios. Visit the [examples/](examples/) directory for the full code.


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
