package ietf

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	HeaderRateLimitLimit     = "RateLimit-Limit"
	HeaderRateLimitRemaining = "RateLimit-Remaining"
	HeaderRateLimitReset     = "RateLimit-Reset"
	HeaderRetryAfter         = "Retry-After"
)

type QuotaPolicy struct {
	quotaComments []QuotaComment
	requestQuota  RequestQuota
	timeWindow    time.Duration
}

func (x QuotaPolicy) Headers(now time.Time, tokens, cost int64) http.Header {
	header := make(http.Header)

	//nolint:canonicalheader // OK.
	header.Add(HeaderRateLimitLimit, x.HeaderLimit())

	//nolint:canonicalheader // OK.
	header.Add(HeaderRateLimitRemaining, strconv.FormatInt(tokens, 10))

	//nolint:canonicalheader // OK.
	header.Add(HeaderRateLimitReset, strconv.FormatInt(int64(x.timeWindow.Seconds()), 10))

	timeWindow := x.timeWindow
	if tokens < cost {
		timeWindow *= time.Duration(cost - tokens)

		header.Add(HeaderRetryAfter, now.Add(timeWindow).Format(http.TimeFormat))
	}

	//nolint:canonicalheader // OK.
	header.Add(HeaderRateLimitReset, strconv.FormatInt(int64(timeWindow.Seconds()), 10))

	return header
}

func (x QuotaPolicy) RequestQuota() RequestQuota {
	return x.requestQuota
}

func (x QuotaPolicy) HeaderLimit() string {
	var b strings.Builder

	//nolint:errcheck // OK.
	_, _ = b.WriteString(strconv.FormatInt(x.RequestQuota().QuotaUnits(), 10))

	//nolint:errcheck // OK.
	_, _ = b.WriteString(", ")

	//nolint:errcheck // OK.
	_, _ = b.WriteString(strconv.FormatInt(x.RequestQuota().QuotaUnits(), 10))

	//nolint:errcheck // OK.
	_, _ = b.WriteString(";w=")

	//nolint:errcheck // OK.
	_, _ = b.WriteString(strconv.FormatInt(int64(x.TimeWindow().Seconds()), 10))

	for _, quotaComment := range x.quotaComments {
		//nolint:errcheck // OK.
		_, _ = b.WriteRune(';')

		//nolint:errcheck // OK.
		_, _ = b.WriteString(quotaComment.Token())

		//nolint:errcheck // OK.
		_, _ = b.WriteRune('=')

		//nolint:errcheck // OK.
		_, _ = b.WriteString(quotaComment.Value())
	}

	return b.String()
}

func (x QuotaPolicy) TimeWindow() time.Duration {
	return x.timeWindow
}

func NewQuotaPolicy(requestQuota RequestQuota, timeWindow time.Duration, opts ...QuotaPolicyOption) QuotaPolicy {
	x := QuotaPolicy{
		quotaComments: nil,
		requestQuota:  requestQuota,
		timeWindow:    timeWindow,
	}

	for _, opt := range opts {
		opt(&x)
	}

	return x
}

type QuotaPolicyOption func(x *QuotaPolicy)

func WithQuotaComments(quotaComments ...QuotaComment) QuotaPolicyOption {
	return func(x *QuotaPolicy) {
		x.quotaComments = quotaComments
	}
}
