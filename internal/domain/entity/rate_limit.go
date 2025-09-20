package entity

import "time"

type RateLimit struct {
	Limit    int
	Interval time.Duration
}

var DefaultRateLimits = map[string]RateLimit{
	"status":    {Limit: 2, Interval: time.Minute},
	"news":      {Limit: 1, Interval: 24 * time.Hour},
	"marketing": {Limit: 3, Interval: time.Hour},
}
