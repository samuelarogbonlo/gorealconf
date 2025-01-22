// rollout.go
package gorealconf

type Rollout[T any] struct {
	config    *Config[T]
	strategy  RolloutStrategy // Changed from Strategy to RolloutStrategy
	validator func(T) error
	threshold float64
}

func NewRollout[T any](config *Config[T]) *Rollout[T] {
	return &Rollout[T]{
		config: config,
	}
}

func (r *Rollout[T]) WithStrategy(strategy RolloutStrategy) *Rollout[T] { // Changed here too
	r.strategy = strategy
	return r
}

func (r *Rollout[T]) WithValidation(validator func(T) error) *Rollout[T] {
	r.validator = validator
	return r
}

func (r *Rollout[T]) WithRollbackThreshold(threshold float64) *Rollout[T] {
	r.threshold = threshold
	return r
}

func (r *Rollout[T]) ShouldApply() bool {
	if r.strategy == nil {
		return true
	}
	return r.strategy.ShouldApply()
}
