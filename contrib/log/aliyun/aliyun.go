package aliyun

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"google.golang.org/protobuf/proto"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*aliyunLog)(nil)

// Logger see more detail https://github.com/aliyun/aliyun-log-go-sdk
type Logger interface {
	log.Logger

	GetClient() sls.ClientInterface
	Close() error
}

type aliyunLog struct {
	client sls.ClientInterface
	opts   *options
}

func (a *aliyunLog) GetClient() sls.ClientInterface {
	return a.client
}

type options struct {
	accessKey     string
	accessSecret  string
	endpoint      string
	project       string
	logstore      string
	securityToken string
}

func defaultOptions() *options {
	return &options{
		project:  "projectName",
		logstore: "app",
	}
}

func WithEndpoint(endpoint string) Option {
	return func(alc *options) {
		alc.endpoint = endpoint
	}
}

func WithProject(project string) Option {
	return func(alc *options) {
		alc.project = project
	}
}

func WithLogstore(logstore string) Option {
	return func(alc *options) {
		alc.logstore = logstore
	}
}

func WithAccessKey(ak string) Option {
	return func(alc *options) {
		alc.accessKey = ak
	}
}

func WithAccessSecret(as string) Option {
	return func(alc *options) {
		alc.accessSecret = as
	}
}

func WithSecurityToken(st string) Option {
	return func(alc *options) {
		alc.securityToken = st
	}
}

type Option func(alc *options)

func (a *aliyunLog) Close() error {
	return a.client.Close()
}

func (a *aliyunLog) Log(level log.Level, keyvals ...any) error {
	contents := make([]*sls.LogContent, 0, len(keyvals)/2+1)

	contents = append(contents, &sls.LogContent{
		Key:   newString(level.Key()),
		Value: newString(level.String()),
	})

	for i := 0; i < len(keyvals); i += 2 {
		contents = append(contents, &sls.LogContent{
			Key:   newString(toString(keyvals[i])),
			Value: newString(toString(keyvals[i+1])),
		})
	}

	logInst := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: contents,
	}

	return a.client.PutLogs(a.opts.project, a.opts.logstore, &sls.LogGroup{Logs: []*sls.Log{logInst}})
}

// NewAliyunLog new aliyun logger with options.
func NewAliyunLog(options ...Option) Logger {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}
	provider := sls.NewStaticCredentialsProvider(opts.accessKey, opts.accessSecret, opts.securityToken)
	client := sls.CreateNormalInterfaceV2(opts.endpoint, provider)

	return &aliyunLog{
		opts:   opts,
		client: client,
	}
}

// newString string convert to *string
func newString(s string) *string {
	return &s
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
