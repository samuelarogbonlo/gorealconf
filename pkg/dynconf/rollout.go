package dynconf

type Rollout[T any] struct {
    config    *Config[T]
    strategy  RolloutStrategy
    validator func(T) error
}

func NewRollout[T any](config *Config[T]) *Rollout[T] {
    return &Rollout[T]{
        config: config,
    }
}

func (r *Rollout[T]) WithStrategy(strategy RolloutStrategy) *Rollout[T] {
    r.strategy = strategy
    return r
}

func (r *Rollout[T]) WithValidation(validator func(T) error) *Rollout[T] {
    r.validator = validator
    return r
}
