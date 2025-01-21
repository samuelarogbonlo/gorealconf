package dynconf

import (
    "time"
)

type RolloutStrategy interface {
    ShouldApply() bool
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
    // Implementation using random number to determine if should apply
    return true
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