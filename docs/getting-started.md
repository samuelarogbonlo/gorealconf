# docs/getting-started.md
# Getting Started with DynConf

## Installation

```bash
go get github.com/samuelarogbonlo/dynconf
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
cfg := dynconf.New[AppConfig](
    dynconf.WithValidation[AppConfig](validateConfig),
    dynconf.WithRollback[AppConfig](true),
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