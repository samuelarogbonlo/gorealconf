# DynConf

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/dynconf.svg)](https://pkg.go.dev/github.com/yourusername/dynconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/dynconf)](https://goreportcard.com/report/github.com/yourusername/dynconf)
[![codecov](https://codecov.io/gh/yourusername/dynconf/branch/main/graph/badge.svg)](https://codecov.io/gh/yourusername/dynconf)

DynConf is a type-safe dynamic configuration management library for Go applications.

## Features

- Type-safe configuration using Go generics
- Zero-downtime configuration updates
- Multiple configuration sources (file, etcd, Consul, Redis)
- Validation and automatic rollback
- Gradual rollout strategies
- Metrics and monitoring

## Requirements

- Go 1.21 or higher
- Compatible with etcd v3.5+
- Compatible with Consul v1.12+

## Installation

```bash
go get github.com/samuelarogbonlo/dynconf
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/samuelarogbonlo/dynconf"
)

type AppConfig struct {
    ServerPort int           `json:"server_port"`
    Timeout    time.Duration `json:"timeout"`
}

func main() {
    cfg := dynconf.New[AppConfig](
        dynconf.WithValidation[AppConfig](func(old, new AppConfig) error {
            if new.ServerPort < 1024 {
                return errors.New("server port must be >= 1024")
            }
            return nil
        }),
        dynconf.WithRollback[AppConfig](true),
    )

    changes, _ := cfg.Watch(context.Background())
    go func() {
        for newCfg := range changes {
            log.Printf("Config updated: %+v", newCfg)
        }
    }()

    // Your application logic
}
```

## Examples

See [examples](examples/) directory for:
- Basic usage
- Web server integration
- Database configuration
- Kubernetes deployment

## Documentation

- [Getting Started](docs/getting-started.md)
- [Configuration Sources](docs/configuration-sources.md)
- [Rollout Strategies](docs/rollout-strategies.md)
- [Metrics](docs/metrics.md)

## Contributing

See [Contributing Guide](docs/contributing.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.