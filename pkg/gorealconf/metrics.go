// metrics.go
package gorealconf

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	configUpdates  *prometheus.CounterVec
	configVersions prometheus.Gauge
	updateDuration prometheus.Histogram
	validationErrs prometheus.Counter
	rollbackCount  prometheus.Counter
	loadErrors     prometheus.Counter
	watchErrors    prometheus.Counter
	updateErrors   prometheus.Counter
}

func NewMetrics(name string) *Metrics {
	return &Metrics{
		configUpdates: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: name + "_config_updates_total",
				Help: "Total number of configuration updates",
			},
			[]string{"source", "valid"},
		),
		configVersions: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: name + "_config_version",
				Help: "Current configuration version",
			},
		),
		updateDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: name + "_update_duration_seconds",
				Help: "Duration of configuration updates",
			},
		),
		validationErrs: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: name + "_validation_errors_total",
				Help: "Total number of validation errors",
			},
		),
		rollbackCount: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: name + "_rollback_total",
				Help: "Total number of configuration rollbacks",
			},
		),
		loadErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: name + "_load_errors_total",
				Help: "Total number of configuration loading errors",
			},
		),
		watchErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: name + "_watch_errors_total",
				Help: "Total number of configuration watch errors",
			},
		),
		updateErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: name + "_update_errors_total",
				Help: "Total number of configuration update errors",
			},
		),
	}
}

// Register registers all metrics with the provided registry
func (m *Metrics) Register(reg prometheus.Registerer) error {
	metrics := []prometheus.Collector{
		m.configUpdates,
		m.configVersions,
		m.updateDuration,
		m.validationErrs,
		m.rollbackCount,
		m.loadErrors,
		m.watchErrors,
		m.updateErrors,
	}

	for _, metric := range metrics {
		if err := reg.Register(metric); err != nil {
			return err
		}
	}

	return nil
}

// Helper methods for the Metrics type
func (m *Metrics) IncValidationErrors() {
	if m != nil && m.validationErrs != nil {
		m.validationErrs.Inc()
	}
}

func (m *Metrics) IncRollbackCount() {
	if m != nil && m.rollbackCount != nil {
		m.rollbackCount.Inc()
	}
}

func (m *Metrics) IncLoadErrors() {
	if m != nil && m.loadErrors != nil {
		m.loadErrors.Inc()
	}
}

func (m *Metrics) IncWatchErrors() {
	if m != nil && m.watchErrors != nil {
		m.watchErrors.Inc()
	}
}

func (m *Metrics) IncUpdateErrors() {
	if m != nil && m.updateErrors != nil {
		m.updateErrors.Inc()
	}
}

func (m *Metrics) ObserveUpdateDuration(duration float64) {
	if m != nil && m.updateDuration != nil {
		m.updateDuration.Observe(duration)
	}
}

func (m *Metrics) SetConfigVersion(version float64) {
	if m != nil && m.configVersions != nil {
		m.configVersions.Set(version)
	}
}

func (m *Metrics) IncConfigUpdates(source, valid string) {
	if m != nil && m.configUpdates != nil {
		m.configUpdates.WithLabelValues(source, valid).Inc()
	}
}
