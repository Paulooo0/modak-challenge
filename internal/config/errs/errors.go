package errs

import "errors"

var (
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
	ErrInvalidNotification = errors.New("invalid notification")
)
