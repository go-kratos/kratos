package polaris

import (
	"time"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/config"
	"github.com/polarismesh/polaris-go/pkg/model"
)

// check

type (
	// Option function for polaris limiter
	Option func(*options)
)

type options struct {
	// polaris config
	conf config.Configuration

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

// WithConfig with polaris config.
func WithConfig(conf config.Configuration) Option {
	return func(o *options) {
		o.conf = conf
	}
}

// WithNamespace with limiter namespace.
func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
	}
}

// WithService with limiter service.
func WithService(service string) Option {
	return func(o *options) {
		o.service = service
	}
}

// WithTimeout with limiter arguments.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithRetryCount with limiter retryCount.
func WithRetryCount(retryCount int) Option {
	return func(o *options) {
		o.retryCount = retryCount
	}
}

// WithToken with limiter token.
func WithToken(token uint32) Option {
	return func(o *options) {
		o.token = token
	}
}

type PolarisLimiter struct {
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
func NewLimiter(opts ...Option) *PolarisLimiter {
	opt := options{
		conf:       nil,
		namespace:  "default",
		service:    "",
		retryCount: 3,
		timeout:    3 * time.Second,
		token:      1,
	}
	for _, o := range opts {
		o(&opt)
	}
	// new limitApi
	limitAPI, err := polaris.NewLimitAPIByConfig(opt.conf)
	if err != nil {
		panic(err)
	}
	return &PolarisLimiter{
		limitAPI: limitAPI,
		opts:     opt,
	}
}

// Allow interface impl
func (l *PolarisLimiter) Allow(method string, argument ...model.Argument) (ratelimit.DoneFunc, error) {
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
