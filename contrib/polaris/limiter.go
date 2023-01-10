package polaris

import (
	"time"

	"github.com/go-kratos/aegis/ratelimit"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"
)

type (
	// LimiterOption function for polaris limiter
	LimiterOption func(*options)
)

type options struct {
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
	return func(o *options) {
		o.namespace = namespace
	}
}

// WithLimiterService with limiter service.
func WithLimiterService(service string) LimiterOption {
	return func(o *options) {
		o.service = service
	}
}

// WithLimiterTimeout with limiter arguments.
func WithLimiterTimeout(timeout time.Duration) LimiterOption {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithLimiterRetryCount with limiter retryCount.
func WithLimiterRetryCount(retryCount int) LimiterOption {
	return func(o *options) {
		o.retryCount = retryCount
	}
}

// WithLimiterToken with limiter token.
func WithLimiterToken(token uint32) LimiterOption {
	return func(o *options) {
		o.token = token
	}
}

type Limiter struct {
	// polaris limit api
	limitAPI polaris.LimitAPI

	opts options
}

// init quotaRequest
func buildRequest(opts options) polaris.QuotaRequest {
	quotaRequest := polaris.NewQuotaRequest()
	quotaRequest.SetNamespace(opts.namespace)
	quotaRequest.SetRetryCount(opts.retryCount)
	quotaRequest.SetService(opts.service)
	quotaRequest.SetTimeout(opts.timeout)
	quotaRequest.SetToken(opts.token)
	return quotaRequest
}

// NewLimiter New a Polaris limiter impl.
func NewLimiter(sdk api.SDKContext, opts ...LimiterOption) *Limiter {
	opt := options{
		namespace:  "default",
		service:    "",
		retryCount: 3,
		timeout:    3 * time.Second,
		token:      1,
	}
	for _, o := range opts {
		o(&opt)
	}
	return &Limiter{
		limitAPI: New(sdk, WithNamespace(opt.namespace)).limit,
		opts:     opt,
	}
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
		return nil, err
	}
	if resp.Get().Code == model.QuotaResultOk {
		return func(ratelimit.DoneInfo) {}, nil
	}
	return nil, ratelimit.ErrLimitExceed
}
