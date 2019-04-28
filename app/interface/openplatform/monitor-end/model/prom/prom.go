package prom

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"go-common/app/interface/openplatform/monitor-end/conf"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	_department   = "open"
	_typeCommon   = "Common"
	_typeDetailed = "Detailed"
	_NaN          = "NaN"
)

var (
	// HTTPClientSum HTTP Client request cost sum.
	HTTPClientSum *Prom
	// HTTPClientCount HTTP Client request count.
	HTTPClientCount *Prom
	// HTTPClientCode HTTP Client request server code count.
	HTTPClientCode *Prom
	// HTTPClientStatus HTTP Client request status.
	HTTPClientStatus *Prom
	// HTTPClientSummary HTTP Client request quantiles.
	HTTPClientSummary *Prom
)

var c *conf.Config

// Prom struct info
type Prom struct {
	timer   *prometheus.HistogramVec
	counter *prometheus.CounterVec
	state   *prometheus.GaugeVec
	summary *prometheus.SummaryVec
}

// Init .
func Init(ce *conf.Config) {
	c = ce
	NewProm(true)
	go clearMemory()
}

func clearMemory() {
	for {
		// default 512MB accloc
		var limit = uint64(512 * 1024 * 1024)
		memStat := new(runtime.MemStats)
		runtime.ReadMemStats(memStat)
		used := memStat.Alloc
		if c.Prom.Limit > 0 && c.Prom.Limit < 3072 {
			limit = uint64(c.Prom.Limit * 1024 * 1024)
		}
		if used > limit {
			NewProm(false)
		}
		time.Sleep(time.Minute)
	}
}

// NewProm .
func NewProm(isFirst bool) {
	if isFirst {
		HTTPClientSum = New().WithCounter("http_client_sum", []string{"target", "client_app", "department", "type", "method", "event", "version", "detail"})
		HTTPClientCount = New().WithCounter("http_client_count", []string{"target", "client_app", "department", "type", "method", "event", "version", "detail"})
		HTTPClientCode = New().WithCounter("http_client_code", []string{"target", "client_app", "department", "type", "method", "event", "version", "code"})
		HTTPClientStatus = New().WithCounter("http_client_status", []string{"target", "client_app", "department", "type", "method", "event", "version", "status"})
		HTTPClientSummary = New().WithQuantile("http_client_summary", []string{"target", "client_app", "department", "type", "method", "event", "version", "detail"})
		return
	}
	var mutex sync.Mutex
	mutex.Lock()
	HTTPClientSum.Unregister()
	HTTPClientSum = New().WithCounter("http_client_sum", []string{"target", "client_app", "department", "type", "method", "event", "version", "detail"})
	HTTPClientCount.Unregister()
	HTTPClientCount = New().WithCounter("http_client_count", []string{"target", "client_app", "department", "type", "method", "event", "version", "detail"})
	HTTPClientCode.Unregister()
	HTTPClientCode = New().WithCounter("http_client_code", []string{"target", "client_app", "department", "type", "method", "event", "version", "code"})
	HTTPClientStatus.Unregister()
	HTTPClientStatus = New().WithCounter("http_client_status", []string{"target", "client_app", "department", "type", "method", "event", "version", "status"})
	HTTPClientSummary.Unregister()
	HTTPClientSummary = New().WithQuantile("http_client_summary", []string{"target", "client_app", "department", "type", "method", "event", "version", "detail"})
	mutex.Unlock()
}

// New creates a Prom instance.
func New() *Prom {
	return &Prom{}
}

// WithTimer with summary timer
func (p *Prom) WithTimer(name string, labels []string) *Prom {
	if p == nil || p.timer != nil {
		return p
	}
	p.timer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.timer)
	return p
}

// WithCounter sets counter.
func (p *Prom) WithCounter(name string, labels []string) *Prom {
	if p == nil || p.counter != nil {
		return p
	}
	p.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.counter)
	return p
}

// WithState sets state.
func (p *Prom) WithState(name string, labels []string) *Prom {
	if p == nil || p.state != nil {
		return p
	}
	p.state = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.state)
	return p
}

// WithQuantile sets quantiles.
func (p *Prom) WithQuantile(name string, labels []string) *Prom {
	if p == nil || p.summary != nil {
		return p
	}
	p.summary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.summary)
	return p
}

// Unregister .
func (p *Prom) Unregister() {
	if p.counter != nil {
		prometheus.Unregister(p.counter)
	}
	if p.state != nil {
		prometheus.Unregister(p.state)
	}
	if p.timer != nil {
		prometheus.Unregister(p.timer)
	}
	if p.summary != nil {
		prometheus.Unregister(p.summary)
	}
}

// Timing log timing information (in milliseconds) without sampling
func (p *Prom) Timing(name string, time int64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.timer != nil {
		p.timer.WithLabelValues(label...).Observe(float64(time))
	}
}

// Incr increments one stat counter without sampling
func (p *Prom) Incr(name string, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counter != nil {
		p.counter.WithLabelValues(label...).Inc()
	}
	if p.state != nil {
		p.state.WithLabelValues(label...).Inc()
	}
}

// Decr decrements one stat counter without sampling
func (p *Prom) Decr(name string, extra ...string) {
	if p.state != nil {
		label := append([]string{name}, extra...)
		p.state.WithLabelValues(label...).Dec()
	}
}

// State set state
func (p *Prom) State(name string, v int64, extra ...string) {
	if p.state != nil {
		label := append([]string{name}, extra...)
		p.state.WithLabelValues(label...).Set(float64(v))
	}
}

// Add add count    v must > 0
func (p *Prom) Add(name string, v int64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counter != nil {
		p.counter.WithLabelValues(label...).Add(float64(v))
	}
	if p.state != nil {
		p.state.WithLabelValues(label...).Add(float64(v))
	}
}

// AddCommonLog .
func AddCommonLog(target string, app string, method string, event string, version string, v int64) {
	if HTTPClientCount.counter == nil || HTTPClientSum.counter == nil || HTTPClientSummary.summary == nil {
		return
	}
	if target == "" || app == "" || method == "" || event == "" {
		return
	}
	var ok = true
	if c.Prom.Factor > 0 && c.Prom.Factor < 100 {
		if rand.Intn(100) >= c.Prom.Factor {
			ok = false
		}
	}
	if version == "" {
		version = _NaN
	}
	// default 1ms per request
	i := float64(v)
	if i <= 0 {
		i = 1
	}
	labels := []string{target, app, _department, _typeCommon, method, event, version, _NaN}
	// labels := append(label, _NaN)
	HTTPClientCount.counter.WithLabelValues(labels...).Inc()
	HTTPClientSum.counter.WithLabelValues(labels...).Add(i)
	if ok {
		HTTPClientSummary.summary.WithLabelValues(labels...).Observe(i)
	}
}

// AddDetailedLog .
func AddDetailedLog(target string, app string, method string, event string, version string, details map[string]int64) {
	if HTTPClientCount.counter == nil || HTTPClientSum.counter == nil || HTTPClientSummary.summary == nil {
		return
	}
	if target == "" || app == "" || method == "" || event == "" || details == nil {
		return
	}
	var ok = true
	if c.Prom.Factor > 0 && c.Prom.Factor < 100 {
		if rand.Intn(100) >= c.Prom.Factor {
			ok = false
		}
	}
	if version == "" {
		version = _NaN
	}
	label := []string{target, app, _department, _typeDetailed, method, event, version}
	for k, v := range details {
		labels := append(label, k)
		// default 1ms per request
		i := float64(v)
		if i <= 0 {
			i = 1
		}
		HTTPClientCount.counter.WithLabelValues(labels...).Inc()
		HTTPClientSum.counter.WithLabelValues(labels...).Add(i)
		if ok {
			HTTPClientSummary.summary.WithLabelValues(labels...).Observe(i)
		}
	}
}

// AddHTTPCode .
func AddHTTPCode(target string, app string, method string, event string, version string, code string) {
	if HTTPClientStatus.counter == nil {
		return
	}
	if target == "" || app == "" || method == "" || event == "" || code == "" {
		return
	}
	if version == "" {
		version = _NaN
	}
	label := []string{target, app, _department, _typeCommon, method, event, version, code}
	HTTPClientStatus.counter.WithLabelValues(label...).Inc()
}

// AddCode .
func AddCode(target string, app string, method string, event string, version string, code string) {
	if HTTPClientCode.counter == nil {
		return
	}
	if target == "" || app == "" || method == "" || event == "" || code == "" {
		return
	}
	if version == "" {
		version = _NaN
	}
	label := []string{target, app, _department, _typeCommon, method, event, version, code}
	HTTPClientCode.counter.WithLabelValues(label...).Inc()
}
