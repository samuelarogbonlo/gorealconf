
# Frequently Asked Questions (FAQ)

### Q: How do I handle configuration file changes?
A: gorealconf automatically watches for file changes when using FileSource:
```go
source := gorealconf.NewFileSource[Config]("config.json")
cfg := gorealconf.New[Config](gorealconf.WithSource(source))
```

### Q: How do I implement custom validation?
A: Use the WithValidation option:
```go
cfg := gorealconf.New[Config](
    gorealconf.WithValidation[Config](func(old, new Config) error {
        // Your validation logic here
        return nil
    }),
)
```

### Q: How do I handle rollback on failure?
A: Enable automatic rollback:
```go
cfg := gorealconf.New[Config](
    gorealconf.WithRollback[Config](true),
)
```

### Q: Can I use multiple configuration sources?
A: Yes, sources are checked in order:
```go
cfg := gorealconf.New[Config](
    gorealconf.WithSource(fileSource),
    gorealconf.WithSource(etcdSource),

```

### Q: How do I monitor configuration changes?
A: Use the metrics integration:
```go
metrics := gorealconf.NewMetrics("myapp")
cfg.WithMetrics(metrics)
```

### Q: How do I implement gradual rollouts?
A: Use rollout strategies:
```go
strategy := gorealconf.NewPercentageStrategy(10) // 10% of instances
rollout := gorealconf.NewRollout[Config](cfg).
    WithStrategy(strategy)
```
