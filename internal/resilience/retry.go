package resilience

import (
	"context"
	"errors"
	"time"
)

// TemporaryError preserves retry classification and the downstream error chain.
type TemporaryError struct {
	After time.Duration
	Err   error
}

func (e *TemporaryError) Error() string { return e.Err.Error() }
func (e *TemporaryError) Unwrap() error { return e.Err }

func IsTemporary(err error) (*TemporaryError, bool) {
	var target *TemporaryError
	ok := errors.As(err, &target)
	return target, ok
}

// Retry calls work until it succeeds, becomes non-retryable, or the request ends.
func Retry(ctx context.Context, max int, work func(context.Context, int) error) error {
	var last error
	for attempt := 1; attempt <= max; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		last = work(ctx, attempt)
		if last == nil {
			return nil
		}
		temporary, ok := IsTemporary(last)
		if !ok {
			return last
		}
		time.Sleep(temporary.After)
	}
	return last
}
