package entity

import "time"

type RateLimit struct {
	Limit    int
	Interval time.Duration
}

var DefaultRateLimits = map[NotificationType]RateLimit{
	Status:    {Limit: 2, Interval: time.Minute},
	News:      {Limit: 1, Interval: 24 * time.Hour},
	Marketing: {Limit: 3, Interval: time.Hour},
}
