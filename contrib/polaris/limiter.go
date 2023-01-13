package polaris

import (
	"time"

	"github.com/go-kratos/aegis/ratelimit"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"
)

type (
	// LimiterOption function for polaris limiter
	LimiterOption func(*limiterOptions)
)

type limiterOptions struct {
	// required, polaris limit namespace
	namespace string

	// required, polaris limit service name
	service string

	// optional, polaris limit request timeout
	// max value is (1+RetryCount) * Timeout
	timeout time.Duration

	// optional, polaris limit retryCount
	// init by polaris config
	retryCount int

	// optional, request limit quota
	token uint32
}

// WithLimiterNamespace with limiter namespace.
func WithLimiterNamespace(namespace string) LimiterOption {
	return func(o *limiterOptions) {
		o.namespace = namespace
	}
}

// WithLimiterService with limiter service.
func WithLimiterService(service string) LimiterOption {
	return func(o *limiterOptions) {
		o.service = service
	}
}

// WithLimiterTimeout with limiter arguments.
func WithLimiterTimeout(timeout time.Duration) LimiterOption {
	return func(o *limiterOptions) {
		o.timeout = timeout
	}
}

// WithLimiterRetryCount with limiter retryCount.
func WithLimiterRetryCount(retryCount int) LimiterOption {
	return func(o *limiterOptions) {
		o.retryCount = retryCount
	}
}

// WithLimiterToken with limiter token.
func WithLimiterToken(token uint32) LimiterOption {
	return func(o *limiterOptions) {
		o.token = token
	}
}

type Limiter struct {
	// polaris limit api
	limitAPI polaris.LimitAPI

	opts limiterOptions
}

// init quotaRequest
func buildRequest(opts limiterOptions) polaris.QuotaRequest {
	quotaRequest := polaris.NewQuotaRequest()
	quotaRequest.SetNamespace(opts.namespace)
	quotaRequest.SetRetryCount(opts.retryCount)
	quotaRequest.SetService(opts.service)
	quotaRequest.SetTimeout(opts.timeout)
	quotaRequest.SetToken(opts.token)
	return quotaRequest
}

// Allow interface impl
func (l *Limiter) Allow(method string, argument ...model.Argument) (ratelimit.DoneFunc, error) {
	request := buildRequest(l.opts)
	request.SetMethod(method)
	for _, arg := range argument {
		request.AddArgument(arg)
	}
	resp, err := l.limitAPI.GetQuota(request)
	if err != nil {
		// ignore err
		return func(ratelimit.DoneInfo) {}, nil
	}
	if resp.Get().Code == model.QuotaResultOk {
		return func(ratelimit.DoneInfo) {}, nil
	}
	return nil, ratelimit.ErrLimitExceed
}
