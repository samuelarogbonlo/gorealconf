package dynconf

import (
    "time"
    "math/rand"
)

type RolloutStrategy interface {
    ShouldApply() bool
}

// Only one definition of CompositeStrategy
type CompositeStrategy struct {
    strategies []RolloutStrategy
}

func NewCompositeStrategy() *CompositeStrategy {
    return &CompositeStrategy{
        strategies: make([]RolloutStrategy, 0),
    }
}

func (cs *CompositeStrategy) Add(s RolloutStrategy) *CompositeStrategy {
    cs.strategies = append(cs.strategies, s)
    return cs
}

func (cs *CompositeStrategy) ShouldApply() bool {
    for _, s := range cs.strategies {
        if !s.ShouldApply() {
            return false
        }
    }
    return true
}

type PercentageStrategy struct {
    percentage float64
}

func NewPercentageStrategy(percentage float64) *PercentageStrategy {
    return &PercentageStrategy{
        percentage: percentage,
    }
}

func (s *PercentageStrategy) ShouldApply() bool {
    return rand.Float64() * 100 <= s.percentage
}

type TimeBasedStrategy struct {
    startTime time.Time
    duration  time.Duration
}

func NewTimeBasedStrategy(duration time.Duration) *TimeBasedStrategy {
    return &TimeBasedStrategy{
        startTime: time.Now(),
        duration:  duration,
    }
}

func (s *TimeBasedStrategy) ShouldApply() bool {
    return time.Since(s.startTime) >= s.duration
}