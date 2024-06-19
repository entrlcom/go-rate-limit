package redis_token_bucket_rate_limit_test

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/redis/rueidis"

	"entrlcom.dev/rate-limit"
	"entrlcom.dev/rate-limit/algorithm/token-bucket/redis"
	"entrlcom.dev/rate-limit/ietf"
	"entrlcom.dev/rate-limit/internal/docker"
)

var _ = ginkgo.Describe("Take", func() {
	type input struct { //nolint:govet // OK.
		cost        int64
		key         string
		quotaPolicy ietf.QuotaPolicy
	}

	var client rueidis.Client

	ginkgo.BeforeEach(func(ctx ginkgo.SpecContext) {
		db, err := docker.NewRedis(ctx)
		gomega.Expect(err).Should(gomega.Succeed())

		client = db.GetClient()
	})

	ginkgo.DescribeTable("", func(ctx context.Context, in input) {
		quota := redis_token_bucket_rate_limit.NewRedisTokenBucketRateLimit(client, in.quotaPolicy)

		for i := 0; i < int(in.quotaPolicy.RequestQuota().QuotaUnits()/in.cost); i++ {
			_, err := quota.Take(ctx, in.key, in.cost)
			gomega.Expect(err).Should(gomega.Succeed())
		}

		tokens, err := quota.Take(ctx, in.key, in.cost)
		gomega.Expect(err).Should(gomega.MatchError(rate_limit.ErrResourceExhausted))
		gomega.Expect(tokens).Should(gomega.BeZero())

		time.Sleep(in.quotaPolicy.TimeWindow() * time.Duration(in.cost))

		_, err = quota.Take(ctx, in.key, in.cost)
		gomega.Expect(err).Should(gomega.Succeed())

		tokens, err = quota.Take(ctx, in.key, in.cost)
		gomega.Expect(err).Should(gomega.MatchError(rate_limit.ErrResourceExhausted))
		gomega.Expect(tokens).Should(gomega.BeZero())
	},
		ginkgo.Entry("", input{
			cost:        int64(gofakeit.IntRange(1, 5)),
			key:         gofakeit.UUID(),
			quotaPolicy: ietf.NewQuotaPolicy(ietf.NewRequestQuota(int64(gofakeit.IntRange(10, 100))), time.Second),
		}),
	)
})

// func logHeaders(headers http.Header) {
//	log.Println(ietf.HeaderRateLimitLimit+":", headers.Get(ietf.HeaderRateLimitLimit))
//	log.Println(ietf.HeaderRateLimitRemaining+":", headers.Get(ietf.HeaderRateLimitRemaining))
//	log.Println(ietf.HeaderRateLimitReset+":", headers.Get(ietf.HeaderRateLimitReset))
//	log.Println(ietf.HeaderRetryAfter+":", headers.Get(ietf.HeaderRetryAfter))
// }.
