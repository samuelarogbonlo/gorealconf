package gorealconf

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Config[T any] struct {
	mu             sync.RWMutex
	current        atomic.Pointer[T]
	version        uint64
	subscribers    map[chan T]struct{}
	validator      func(old, new T) error
	sources        []Source[T]
	enableRollback bool
	metrics        *Metrics
}

type Option[T any] func(*Config[T])

func New[T any](opts ...Option[T]) *Config[T] {
	cfg := &Config[T]{
		subscribers: make(map[chan T]struct{}),
		sources:     make([]Source[T], 0),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// Load initializes the configuration from all sources
func (c *Config[T]) Load(ctx context.Context) error {
	for _, source := range c.sources {
		value, err := source.Load(ctx)
		if err != nil {
			if c.metrics != nil {
				c.metrics.loadErrors.Inc()
			}
			return fmt.Errorf("failed to load config from source: %w", err)
		}

		// Apply validation if configured
		if c.validator != nil {
			oldValue := c.Get(ctx)
			if err := c.validator(oldValue, value); err != nil {
				if c.metrics != nil {
					c.metrics.validationErrs.Inc()
				}
				if c.enableRollback {
					if c.metrics != nil {
						c.metrics.rollbackCount.Inc()
					}
					return fmt.Errorf("validation failed: %w", err)
				}
			}
		}

		if err := c.Update(ctx, value); err != nil {
			return fmt.Errorf("failed to update config: %w", err)
		}
	}

	// Start watching all sources
	for _, source := range c.sources {
		go c.watchSource(ctx, source)
	}

	return nil
}

func (c *Config[T]) watchSource(ctx context.Context, source Source[T]) {
	changes, err := source.Watch(ctx)
	if err != nil {
		if c.metrics != nil {
			c.metrics.watchErrors.Inc()
		}
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-changes:
			if !ok {
				return
			}
			if err := c.Update(ctx, value); err != nil {
				if c.metrics != nil {
					c.metrics.updateErrors.Inc()
				}
			}
		}
	}
}

func WithValidation[T any](validator func(old, new T) error) Option[T] {
	return func(c *Config[T]) {
		c.validator = validator
	}
}

func WithRollback[T any](enable bool) Option[T] {
	return func(c *Config[T]) {
		c.enableRollback = enable
	}
}

func WithMetrics[T any](metrics *Metrics) Option[T] {
	return func(c *Config[T]) {
		c.metrics = metrics
	}
}

func WithSource[T any](source Source[T]) Option[T] {
	return func(c *Config[T]) {
		c.AddSource(source)
	}
}

func (c *Config[T]) Get(ctx context.Context) T {
	value := c.current.Load()
	if value == nil {
		var zero T
		return zero
	}
	return *value
}

func (c *Config[T]) Update(ctx context.Context, newValue T) error {
	start := time.Now()
	oldValue := c.Get(ctx)

	if c.validator != nil {
		if err := c.validator(oldValue, newValue); err != nil {
			if c.metrics != nil {
				c.metrics.validationErrs.Inc()
			}
			if c.enableRollback {
				if c.metrics != nil {
					c.metrics.rollbackCount.Inc()
				}
				return err
			}
		}
	}

	c.mu.Lock()
	c.current.Store(&newValue)
	newVersion := atomic.AddUint64(&c.version, 1)

	if c.metrics != nil {
		c.metrics.configUpdates.WithLabelValues("manual", "true").Inc()
		c.metrics.configVersions.Set(float64(newVersion))
		c.metrics.updateDuration.Observe(time.Since(start).Seconds())
	}

	subscribers := make([]chan T, 0, len(c.subscribers))
	for sub := range c.subscribers {
		subscribers = append(subscribers, sub)
	}
	c.mu.Unlock()

	go func() {
		for _, sub := range subscribers {
			select {
			case sub <- newValue:
			default:
				c.mu.Lock()
				delete(c.subscribers, sub)
				c.mu.Unlock()
				close(sub)
			}
		}
	}()

	return nil
}

func (c *Config[T]) Subscribe(ctx context.Context) (<-chan T, func()) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan T, 1)
	c.subscribers[ch] = struct{}{}

	if current := c.current.Load(); current != nil {
		select {
		case ch <- *current:
		default:
		}
	}

	return ch, func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.subscribers, ch)
		close(ch)
	}
}

func (c *Config[T]) Watch(ctx context.Context) (<-chan T, error) {
	ch, cleanup := c.Subscribe(ctx)

	watchCh := make(chan T)

	go func() {
		defer cleanup()
		defer close(watchCh)

		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch:
				if !ok {
					return
				}
				select {
				case watchCh <- val:
				default:
				}
			}
		}
	}()

	return watchCh, nil
}

func (c *Config[T]) AddSource(source Source[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sources = append(c.sources, source)
}
