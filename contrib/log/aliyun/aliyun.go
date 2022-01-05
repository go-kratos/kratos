package aliyun

import (
	"encoding/json"
	"strconv"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	klog "github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
)

// Log see more detail https://github.com/aliyun/aliyun-log-go-sdk
type Log interface {
	klog.Logger
	GetProducer() *producer.Producer
}

type aliyunLog struct {
	producer *producer.Producer
	config   *LogConfig
}

func (a *aliyunLog) GetProducer() *producer.Producer {
	return a.producer
}

type LogConfig struct {
	AccessKey    string
	AccessSecret string
	Endpoint     string
	Project      string
	Logstore     string
}

func DefaultConfig() *LogConfig {
	return &LogConfig{
		Project:  "projectName",
		Logstore: "app",
	}
}

func WithEndpoint(endpoint string) Option {
	return func(alc *LogConfig) {
		alc.Endpoint = endpoint
	}
}

func WithProject(project string) Option {
	return func(alc *LogConfig) {
		alc.Project = project
	}
}

func WithLogstore(logstore string) Option {
	return func(alc *LogConfig) {
		alc.Logstore = logstore
	}
}

func WithAccessKey(ak string) Option {
	return func(alc *LogConfig) {
		alc.AccessKey = ak
	}
}

func WithAccessSecret(as string) Option {
	return func(alc *LogConfig) {
		alc.AccessSecret = as
	}
}

type Option func(alc *LogConfig)

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

	err := a.producer.SendLog(a.config.Project, a.config.Logstore, "", "", logInst)
	return err
}

// NewAliyunLog new a aliyun logger with options.
func NewAliyunLog(options ...Option) Log {
	aliyunConfig := DefaultConfig()
	for _, o := range options {
		o(aliyunConfig)
	}

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = aliyunConfig.Endpoint
	producerConfig.AccessKeyID = aliyunConfig.AccessKey
	producerConfig.AccessKeySecret = aliyunConfig.AccessSecret
	producerInst := producer.InitProducer(producerConfig)

	return &aliyunLog{
		config:   aliyunConfig,
		producer: producerInst,
	}
}

// toString 任意类型转string
func toString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}
