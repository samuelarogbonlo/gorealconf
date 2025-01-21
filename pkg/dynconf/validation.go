package dynconf

import (
    "context"
    "time"
)

// Validator represents a configuration validator
type Validator[T any] interface {
    Validate(old, new T) error
}

// AsyncValidator represents an asynchronous configuration validator
type AsyncValidator[T any] interface {
    ValidateAsync(ctx context.Context, old, new T) error
}

// ValidationResult represents the result of a validation
type ValidationResult struct {
    Valid   bool
    Message string
}

// ValidationOptions provides options for validation
type ValidationOptions struct {
    Timeout     time.Duration
    RetryCount  int
    RetryDelay  time.Duration
}

// DefaultValidationOptions returns default validation options
func DefaultValidationOptions() ValidationOptions {
    return ValidationOptions{
        Timeout:    5 * time.Second,
        RetryCount: 3,
        RetryDelay: time.Second,
    }
}