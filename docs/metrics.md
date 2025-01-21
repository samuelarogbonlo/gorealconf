# docs/metrics.md
# Metrics

DynConf provides Prometheus metrics for monitoring:

## Available Metrics

- `config_updates_total`: Counter of configuration updates
- `config_version`: Gauge of current configuration version
- `validation_errors_total`: Counter of validation errors
- `rollbacks_total`: Counter of configuration rollbacks
- `update_duration_seconds`: Histogram of update durations

## Usage

```go
metrics := dynconf.NewMetrics("myapp")
cfg.WithMetrics(metrics)
```
