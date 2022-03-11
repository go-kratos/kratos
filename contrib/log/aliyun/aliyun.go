package aliyun

import (
	"encoding/json"
	"strconv"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	klog "github.com/SeeMusic/kratos/v2/log"
	"google.golang.org/protobuf/proto"
)

// Log see more detail https://github.com/aliyun/aliyun-log-go-sdk
type Log interface {
	klog.Logger
	GetProducer() *producer.Producer
}

type aliyunLog struct {
	producer *producer.Producer
	opts     *options
}

func (a *aliyunLog) GetProducer() *producer.Producer {
	return a.producer
}

type options struct {
	accessKey    string
	accessSecret string
	endpoint     string
	project      string
	logstore     string
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

type Option func(alc *options)

func (a *aliyunLog) Log(level klog.Level, keyvals ...interface{}) error {
	buf := level.String()
	levelTitle := "level"

	contents := make([]*sls.LogContent, 0)

	contents = append(contents, &sls.LogContent{
		Key:   &levelTitle,
		Value: &buf,
	})

	for i := 0; i < len(keyvals); i += 2 {
		key := toString(keyvals[i])
		value := toString(keyvals[i+1])
		contents = append(contents, &sls.LogContent{
			Key:   &key,
			Value: &value,
		})
	}

	logInst := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: contents,
	}

	err := a.producer.SendLog(a.opts.project, a.opts.logstore, "", "", logInst)
	return err
}

// NewAliyunLog new a aliyun logger with options.
func NewAliyunLog(options ...Option) Log {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = opts.endpoint
	producerConfig.AccessKeyID = opts.accessKey
	producerConfig.AccessKeySecret = opts.accessSecret
	producerInst := producer.InitProducer(producerConfig)

	return &aliyunLog{
		opts:     opts,
		producer: producerInst,
	}
}

// toString 任意类型转string
func toString(v interface{}) string {
	var key string
	if v == nil {
		return key
	}
	switch v := v.(type) {
	case float64:
		key = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		key = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case int:
		key = strconv.Itoa(v)
	case uint:
		key = strconv.Itoa(int(v))
	case int8:
		key = strconv.Itoa(int(v))
	case uint8:
		key = strconv.Itoa(int(v))
	case int16:
		key = strconv.Itoa(int(v))
	case uint16:
		key = strconv.Itoa(int(v))
	case int32:
		key = strconv.Itoa(int(v))
	case uint32:
		key = strconv.Itoa(int(v))
	case int64:
		key = strconv.FormatInt(v, 10)
	case uint64:
		key = strconv.FormatUint(v, 10)
	case string:
		key = v
	case []byte:
		key = string(v)
	default:
		newValue, _ := json.Marshal(v)
		key = string(newValue)
	}
	return key
}
