package rate_limit

import (
	"context"
	"errors"

	"entrlcom.dev/rate-limit/ietf"
)

var (
	ErrInternal          = errors.New("internal error")
	ErrResourceExhausted = errors.New("resource exhausted")
)

type RateLimit interface {
	QuotaPolicy() ietf.QuotaPolicy
	Take(ctx context.Context, key string, cost int64) (int64, error)
}
