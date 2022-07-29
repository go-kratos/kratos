package tencent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	log "github.com/go-kratos/kratos/v2/log"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"google.golang.org/protobuf/proto"
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
	err := log.producer.Close(5000)
	return err
}

func (log *tencentLog) Log(level log.Level, keyvals ...interface{}) error {
	buf := level.String()
	levelTitle := "level"
	contents := make([]*cls.Log_Content, 0)

	contents = append(contents, &cls.Log_Content{
		Key:   &levelTitle,
		Value: &buf,
	})
	for i := 0; i < len(keyvals); i += 2 {
		key := toString(keyvals[i])
		value := toString(keyvals[i+1])
		contents = append(contents, &cls.Log_Content{
			Key:   &key,
			Value: &value,
		})
	}

	logInst := &cls.Log{
		Time:     proto.Int64(time.Now().Unix()),
		Contents: contents,
	}
	err := log.producer.SendLog(log.opts.topicID, logInst, nil)
	return err
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
	return &tencentLog{
		producer: producerInst,
		opts:     opts,
	}, nil
}

// toString any type to string
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
	case fmt.Stringer:
		key = v.String()
	default:
		newValue, _ := json.Marshal(v)
		key = string(newValue)
	}
	return key
}
