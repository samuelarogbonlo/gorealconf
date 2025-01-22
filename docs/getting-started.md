# docs/getting-started.md
# Getting Started with gorealconf

## Installation

```bash
go get github.com/samuelarogbonlo/gorealconf
```

## Basic Usage

1. Define your configuration type:

```go
type AppConfig struct {
    ServerPort int           `json:"server_port"`
    Timeout    time.Duration `json:"timeout"`
}
```

2. Create a configuration instance:

```go
cfg := gorealconf.New[AppConfig](
    gorealconf.WithValidation[AppConfig](validateConfig),
    gorealconf.WithRollback[AppConfig](true),
)
```

3. Watch for changes:

```go
changes, err := cfg.Watch(context.Background())
if err != nil {
    log.Fatal(err)
}

go func() {
    for newCfg := range changes {
        // Handle configuration updates
    }
}()
```