package tencent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"google.golang.org/protobuf/proto"

	"github.com/go-kratos/kratos/v2/log"
)

type Logger interface {
	log.Logger

	GetProducer() *cls.AsyncProducerClient
	Close() error
}

type tencentLog struct {
	producer *cls.AsyncProducerClient
	opts     *options
}

func (log *tencentLog) GetProducer() *cls.AsyncProducerClient {
	return log.producer
}

type options struct {
	topicID      string
	accessKey    string
	accessSecret string
	endpoint     string
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

type Option func(cls *options)

func (log *tencentLog) Close() error {
	return log.producer.Close(5000)
}

func (log *tencentLog) Log(level log.Level, keyvals ...interface{}) error {
	contents := make([]*cls.Log_Content, 0, len(keyvals)/2+1)

	contents = append(contents, &cls.Log_Content{
		Key:   newString(level.Key()),
		Value: newString(level.String()),
	})
	for i := 0; i < len(keyvals); i += 2 {
		contents = append(contents, &cls.Log_Content{
			Key:   newString(toString(keyvals[i])),
			Value: newString(toString(keyvals[i+1])),
		})
	}

	logInst := &cls.Log{
		Time:     proto.Int64(time.Now().Unix()),
		Contents: contents,
	}
	return log.producer.SendLog(log.opts.topicID, logInst, nil)
}

func NewLogger(options ...Option) (Logger, error) {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}
	producerConfig := cls.GetDefaultAsyncProducerClientConfig()
	producerConfig.AccessKeyID = opts.accessKey
	producerConfig.AccessKeySecret = opts.accessSecret
	producerConfig.Endpoint = opts.endpoint
	producerInst, err := cls.NewAsyncProducerClient(producerConfig)
	if err != nil {
		return nil, err
	}
	producerInst.Start()
	return &tencentLog{
		producer: producerInst,
		opts:     opts,
	}, nil
}

func newString(s string) *string {
	return &s
}

// toString convert any type to string
func toString(v interface{}) string {
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
