
# Frequently Asked Questions (FAQ)

### Q: How do I handle configuration file changes?
A: DynConf automatically watches for file changes when using FileSource:
```go
source := dynconf.NewFileSource[Config]("config.json")
cfg := dynconf.New[Config](dynconf.WithSource(source))
```

### Q: How do I implement custom validation?
A: Use the WithValidation option:
```go
cfg := dynconf.New[Config](
    dynconf.WithValidation[Config](func(old, new Config) error {
        // Your validation logic here
        return nil
    }),
)
```

### Q: How do I handle rollback on failure?
A: Enable automatic rollback:
```go
cfg := dynconf.New[Config](
    dynconf.WithRollback[Config](true),
)
```

### Q: Can I use multiple configuration sources?
A: Yes, sources are checked in order:
```go
cfg := dynconf.New[Config](
    dynconf.WithSource(fileSource),
    dynconf.WithSource(etcdSource),

```

### Q: How do I monitor configuration changes?
A: Use the metrics integration:
```go
metrics := dynconf.NewMetrics("myapp")
cfg.WithMetrics(metrics)
```

### Q: How do I implement gradual rollouts?
A: Use rollout strategies:
```go
strategy := dynconf.NewPercentageStrategy(10) // 10% of instances
rollout := dynconf.NewRollout[Config](cfg).
    WithStrategy(strategy)
```
