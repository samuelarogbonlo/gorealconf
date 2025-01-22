package gorealconf

import "context"

// Source represents a configuration source
type Source[T any] interface {
	// Load loads the initial configuration
	Load(ctx context.Context) (T, error)

	// Watch watches for configuration changes
	Watch(ctx context.Context) (<-chan T, error)
}
