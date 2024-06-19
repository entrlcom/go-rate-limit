package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/redis/rueidis"

	"entrlcom.dev/rate-limit"
	"entrlcom.dev/rate-limit/algorithm/token-bucket/redis"
	"entrlcom.dev/rate-limit/ietf"
)

func main() {
	opt, err := rueidis.ParseURL("redis://:password@127.0.0.1:6379?protocol=3")
	if err != nil {
		// TODO: Handle error.
		return
	}

	client, err := rueidis.NewClient(opt)
	if err != nil {
		// TODO: Handle error.
		return
	}

	quotaPolicy := ietf.NewQuotaPolicy(ietf.NewRequestQuota(250), time.Second)
	rateLimit := redis_token_bucket_rate_limit.NewRedisTokenBucketRateLimit(client, quotaPolicy)

	handler := NewHandler(&rateLimit)

	// ...
}

type Handler struct {
	rateLimit rate_limit.RateLimit
}

func (x *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if tokens, err := x.rateLimit.Take(r.Context(), r.RemoteAddr, 5); err != nil {
		if errors.Is(err, rate_limit.ErrResourceExhausted) {
			headers := x.rateLimit.QuotaPolicy().Headers(time.Now().UTC(), tokens, 5)
			if err := headers.Write(w); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// ...
}

func NewHandler(rateLimit rate_limit.RateLimit) Handler {
	return Handler{rateLimit: rateLimit}
}
