package dynconf

import "fmt"

// ValidationError represents a configuration validation error
type ValidationError struct {
    Message string
    Old     interface{}
    New     interface{}
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed: %s", e.Message)
}

// RollbackError represents a rollback failure
type RollbackError struct {
    Message string
    Cause   error
}

func (e *RollbackError) Error() string {
    return fmt.Sprintf("rollback failed: %s: %v", e.Message, e.Cause)
}
