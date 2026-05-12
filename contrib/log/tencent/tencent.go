package tencent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"google.golang.org/protobuf/proto"

	klog "github.com/go-kratos/kratos/v2/log"
)

// Handler writes slog records to Tencent CLS.
type Handler struct {
	producer *cls.AsyncProducerClient
	opts     *options
	attrs    []groupedAttr
	groups   []string
}

type groupedAttr struct {
	groups []string
	attr   slog.Attr
}

func (h *Handler) GetProducer() *cls.AsyncProducerClient {
	return h.producer
}

type options struct {
	topicID      string
	accessKey    string
	accessSecret string
	endpoint     string
	logOptions   []klog.Option
}

func defaultOptions() *options {
	return &options{}
}

func WithEndpoint(endpoint string) Option {
	return func(cls *options) {
		cls.endpoint = endpoint
	}
}

func WithTopicID(topicID string) Option {
	return func(cls *options) {
		cls.topicID = topicID
	}
}

func WithAccessKey(ak string) Option {
	return func(cls *options) {
		cls.accessKey = ak
	}
}

func WithAccessSecret(as string) Option {
	return func(cls *options) {
		cls.accessSecret = as
	}
}

// WithLogOptions applies Kratos core log builder options to NewLogger.
func WithLogOptions(opts ...klog.Option) Option {
	return func(o *options) {
		o.logOptions = append(o.logOptions, opts...)
	}
}

// Option configures the Tencent CLS handler or the NewLogger wrapper.
type Option func(o *options)

func (h *Handler) Close() error {
	return h.producer.Close(5000)
}

func (h *Handler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	contents := make([]*cls.Log_Content, 0, len(h.attrs)+record.NumAttrs()+2)
	contents = append(contents, &cls.Log_Content{
		Key:   newString(slog.LevelKey),
		Value: newString(levelString(record.Level)),
	})
	if record.Message != "" {
		contents = append(contents, &cls.Log_Content{
			Key:   newString("msg"),
			Value: newString(record.Message),
		})
	}
	for _, attr := range h.attrs {
		appendAttr(&contents, attr.groups, attr.attr)
	}
	record.Attrs(func(attr slog.Attr) bool {
		appendAttr(&contents, h.groups, attr)
		return true
	})

	logInst := &cls.Log{
		Time:     proto.Int64(time.Now().Unix()),
		Contents: contents,
	}
	return h.producer.SendLog(h.opts.topicID, logInst, nil)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = append([]groupedAttr{}, h.attrs...)
	for _, attr := range attrs {
		next.attrs = append(next.attrs, groupedAttr{
			groups: append([]string{}, h.groups...),
			attr:   attr,
		})
	}
	return &next
}

func (h *Handler) WithGroup(name string) slog.Handler {
	next := *h
	next.groups = append(append([]string{}, h.groups...), name)
	return &next
}

// NewHandler returns a slog handler backed by Tencent CLS.
func NewHandler(options ...Option) (*Handler, error) {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}
	return newHandler(opts)
}

func newHandler(opts *options) (*Handler, error) {
	producerConfig := cls.GetDefaultAsyncProducerClientConfig()
	producerConfig.AccessKeyID = opts.accessKey
	producerConfig.AccessKeySecret = opts.accessSecret
	producerConfig.Endpoint = opts.endpoint
	producerInst, err := cls.NewAsyncProducerClient(producerConfig)
	if err != nil {
		return nil, err
	}
	producerInst.Start()
	return &Handler{
		producer: producerInst,
		opts:     opts,
	}, nil
}

// NewLogger returns a slog logger backed by Tencent CLS.
func NewLogger(options ...Option) (*slog.Logger, error) {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}
	handler, err := newHandler(opts)
	if err != nil {
		return nil, err
	}
	logOptions := append([]klog.Option{}, opts.logOptions...)
	logOptions = append(logOptions, klog.WithHandler(handler))
	return klog.NewLogger(logOptions...), nil
}

func newString(s string) *string {
	return &s
}

func appendAttr(contents *[]*cls.Log_Content, groups []string, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Value.Kind() == slog.KindGroup {
		nextGroups := groups
		if attr.Key != "" {
			nextGroups = append(append([]string{}, groups...), attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			appendAttr(contents, nextGroups, groupAttr)
		}
		return
	}
	key := attr.Key
	if len(groups) > 0 {
		key = strings.Join(append(append([]string{}, groups...), key), ".")
	}
	*contents = append(*contents, &cls.Log_Content{
		Key:   newString(key),
		Value: newString(toString(attr.Value.Any())),
	})
}

func levelString(level slog.Level) string {
	switch {
	case level <= slog.LevelDebug:
		return "DEBUG"
	case level < slog.LevelWarn:
		return "INFO"
	case level < slog.LevelError:
		return "WARN"
	case level < slog.LevelError+4:
		return "ERROR"
	default:
		return "FATAL"
	}
}

// toString convert any type to string
func toString(v any) string {
	var key string
	if v == nil {
		return key
	}
	switch v := v.(type) {
	case float64:
		key = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		key = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		key = strconv.Itoa(v)
	case uint:
		key = strconv.FormatUint(uint64(v), 10)
	case int8:
		key = strconv.Itoa(int(v))
	case uint8:
		key = strconv.FormatUint(uint64(v), 10)
	case int16:
		key = strconv.Itoa(int(v))
	case uint16:
		key = strconv.FormatUint(uint64(v), 10)
	case int32:
		key = strconv.Itoa(int(v))
	case uint32:
		key = strconv.FormatUint(uint64(v), 10)
	case int64:
		key = strconv.FormatInt(v, 10)
	case uint64:
		key = strconv.FormatUint(v, 10)
	case string:
		key = v
	case bool:
		key = strconv.FormatBool(v)
	case []byte:
		key = string(v)
	case fmt.Stringer:
		key = v.String()
	default:
		newValue, _ := json.Marshal(v)
		key = string(newValue)
	}
	return key
}
