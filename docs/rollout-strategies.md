# docs/rollout-strategies.md
# Rollout Strategies

DynConf supports gradual configuration rollouts:

## Percentage-based Rollout

```go
strategy := dynconf.NewPercentageStrategy(10) // 10% of instances
rollout := dynconf.NewRollout[FeatureConfig](cfg).
    WithStrategy(strategy)
```

## Time-based Rollout

```go
strategy := dynconf.NewTimeBasedStrategy(24 * time.Hour) // Over 24 hours
rollout := dynconf.NewRollout[FeatureConfig](cfg



<!-- # Rollout Strategies

DynConf provides several strategies for gradually rolling out configuration changes across your infrastructure.

## Available Strategies

### 1. Percentage-Based Rollout
Applies changes to a specified percentage of instances.

```go
strategy := dynconf.NewPercentageStrategy(10) // 10% of instances
rollout := dynconf.NewRollout[FeatureConfig](cfg).
    WithStrategy(strategy).
    WithValidation(validateFeature)

if rollout.ShouldApply() {
    // Apply new configuration
}
```

### 2. Time-Based Rollout
Gradually rolls out changes over a specified duration.

```go
strategy := dynconf.NewTimeBasedStrategy(24 * time.Hour) // Roll out over 24 hours
rollout := dynconf.NewRollout[FeatureConfig](cfg).
    WithStrategy(strategy)
```

### 3. Region-Based Rollout
Rolls out changes to specific regions first.

```go
strategy := dynconf.NewRegionBasedStrategy(
    []string{"us-west-1", "us-west-2"},
    currentRegion,
)
rollout := dynconf.NewRollout[FeatureConfig](cfg).
    WithStrategy(strategy)
```

## Advanced Usage

### Combining Strategies
You can combine multiple strategies for more complex rollouts:

```go
rollout := dynconf.NewRollout[FeatureConfig](cfg).
    WithStrategy(dynconf.NewCompositeStrategy().
        Add(dynconf.NewRegionBasedStrategy(regions, currentRegion)).
        Add(dynconf.NewPercentageStrategy(20)),
    )
```

### Automatic Rollback
Configure automatic rollback based on error rates:

```go
rollout := dynconf.NewRollout[FeatureConfig](cfg).
    WithStrategy(dynconf.NewPercentageStrategy(10)).
    WithRollbackThreshold(0.01) // Rollback if error rate exceeds 1%
```

### Health Checks
Add health checks to validate the rollout:

```go
rollout := dynconf.NewRollout[FeatureConfig](cfg).
    WithHealthCheck(func(cfg FeatureConfig) error {
        // Verify the configuration is working as expected
        return checkFeatureHealth(cfg)
    })
```

## Monitoring Rollouts

### Metrics
Monitor rollout progress using built-in metrics:

```go
metrics := dynconf.NewMetrics("myapp_rollout")
rollout.WithMetrics(metrics)

// Available metrics:
// - rollout_progress_percentage
// - rollout_instances_updated
// - rollout_errors_total
// - rollout_duration_seconds
```

### Logging
Enable detailed rollout logging:

```go
rollout.OnProgress(func(event RolloutEvent) {
    log.Printf("Rollout progress: %d%% complete", event.ProgressPercentage)
    log.Printf("Instances updated: %d/%d", event.UpdatedCount, event.TotalCount)
})
```

## Best Practices

1. **Start Small**: Begin with a small percentage (e.g., 5-10%) to catch potential issues early.

2. **Monitor Closely**: Use metrics and logging to track rollout progress and health.

3. **Plan for Rollback**: Always have a rollback strategy and test it before starting the rollout.

4. **Gradual Increase**: Gradually increase the rollout percentage after validating each phase.

5. **Region Strategy**: Consider starting with less critical regions first.

## Example: Complete Rollout Configuration

```go
func setupRollout(cfg *dynconf.Config[FeatureConfig]) *dynconf.Rollout[FeatureConfig] {
    return dynconf.NewRollout[FeatureConfig](cfg).
        WithStrategy(dynconf.NewCompositeStrategy().
            Add(dynconf.NewRegionBasedStrategy([]string{"us-west-2"}, getRegion())).
            Add(dynconf.NewPercentageStrategy(10)),
        ).
        WithValidation(validateFeature).
        WithHealthCheck(checkFeatureHealth).
        WithRollbackThreshold(0.01).
        WithMetrics(dynconf.NewMetrics("feature_rollout")).
        OnProgress(logRolloutProgress)
}

func validateFeature(cfg FeatureConfig) error {
    // Validation logic
    return nil
}

func checkFeatureHealth(cfg FeatureConfig) error {
    // Health check logic
    return nil
}

func logRolloutProgress(event RolloutEvent) {
    log.Printf("Rollout progress: %+v", event)
}
``` -->