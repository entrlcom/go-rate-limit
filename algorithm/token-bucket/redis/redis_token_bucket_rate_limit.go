package redis_token_bucket_rate_limit

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/rueidis"

	"entrlcom.dev/rate-limit"
	"entrlcom.dev/rate-limit/ietf"

	_ "embed"
)

const n = 2

//go:embed script.lua
var script string

type RedisTokenBucketRateLimit struct { //nolint:govet // OK.
	client      rueidis.Client
	quotaPolicy ietf.QuotaPolicy
	script      *rueidis.Lua
}

func (x *RedisTokenBucketRateLimit) QuotaPolicy() ietf.QuotaPolicy {
	return x.quotaPolicy
}

func (x *RedisTokenBucketRateLimit) Take(ctx context.Context, key string, cost int64) (int64, error) {
	now := time.Now().UTC()

	args := []string{
		strconv.FormatInt(now.UnixMilli(), 10),
		strconv.FormatInt(x.quotaPolicy.RequestQuota().QuotaUnits(), 10),
		"1",
		strconv.FormatInt(x.quotaPolicy.TimeWindow().Milliseconds(), 10),
		strconv.FormatInt(cost, 10),
	}

	v, err := x.script.Exec(ctx, x.client, []string{key}, args).ToArray()
	if err != nil {
		return 0, errors.Join(err, rate_limit.ErrInternal)
	}

	if len(v) != n {
		return 0, rate_limit.ErrInternal
	}

	ok, err := v[0].AsInt64()
	if err != nil {
		return 0, errors.Join(err, rate_limit.ErrInternal)
	}

	if ok != 1 {
		return 0, rate_limit.ErrResourceExhausted
	}

	tokens, err := v[1].AsInt64()
	if err != nil {
		return 0, errors.Join(err, rate_limit.ErrInternal)
	}

	return tokens, nil
}

func NewRedisTokenBucketRateLimit(client rueidis.Client, quotaPolicy ietf.QuotaPolicy) RedisTokenBucketRateLimit {
	return RedisTokenBucketRateLimit{
		client:      client,
		quotaPolicy: quotaPolicy,
		script:      rueidis.NewLuaScript(script),
	}
}

var _ rate_limit.RateLimit = (*RedisTokenBucketRateLimit)(nil)
