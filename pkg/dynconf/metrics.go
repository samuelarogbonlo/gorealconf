package dynconf

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
    configUpdates   *prometheus.CounterVec
    configVersions  prometheus.Gauge
    validationErrs  prometheus.Counter
    rollbackCount   prometheus.Counter
    updateDuration  prometheus.Histogram
}

func NewMetrics(namespace string) *Metrics {
    return &Metrics{
        configUpdates: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Name:     "config_updates_total",
                Help:     "Total number of configuration updates",
            },
            []string{"source", "success"},
        ),
        configVersions: promauto.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:     "config_version",
                Help:     "Current configuration version",
            },
        ),
        validationErrs: promauto.NewCounter(
            prometheus.CounterOpts{
                Namespace: namespace,
                Name:     "validation_errors_total",
                Help:     "Total number of validation errors",
            },
        ),
        rollbackCount: promauto.NewCounter(
            prometheus.CounterOpts{
                Namespace: namespace,
                Name:     "rollbacks_total",
                Help:     "Total number of configuration rollbacks",
            },
        ),
        updateDuration: promauto.NewHistogram(
            prometheus.HistogramOpts{
                Namespace: namespace,
                Name:     "update_duration_seconds",
                Help:     "Duration of configuration updates",
                Buckets:  prometheus.DefBuckets,
            },
        ),
    }
}