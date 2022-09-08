package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

// NewStdLogger new a logger with writer.
func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		log: log.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log print the kv pairs log.
func (l *stdLogger) Log(level Level, keyvals ...interface{}) error {
	if len(keyvals)&1 == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	kvMap := make(map[string]interface{}, len(keyvals)/2+1)
	kvMap[level.Key()] = level.String()
	for i := 0; i < len(keyvals); i += 2 {
		key, val := toString(keyvals[i]), toString(keyvals[i+1])
		kvMap[key] = val
	}
	b, err := json.Marshal(kvMap)
	if err != nil {
		return err
	}

	buf := l.pool.Get().(*bytes.Buffer)
	_, _ = buf.Write(append(b, '\n'))
	_ = l.log.Output(4, buf.String()) //nolint:gomnd
	buf.Reset()
	l.pool.Put(buf)
	return nil
}

func toString(v interface{}) string {
	var str string
	if v == nil {
		return str
	}
	switch v := v.(type) {
	case float64:
		str = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		str = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		str = strconv.Itoa(v)
	case uint:
		str = strconv.FormatUint(uint64(v), 10)
	case int8:
		str = strconv.Itoa(int(v))
	case uint8:
		str = strconv.FormatUint(uint64(v), 10)
	case int16:
		str = strconv.Itoa(int(v))
	case uint16:
		str = strconv.FormatUint(uint64(v), 10)
	case int32:
		str = strconv.Itoa(int(v))
	case uint32:
		str = strconv.FormatUint(uint64(v), 10)
	case int64:
		str = strconv.FormatInt(v, 10)
	case uint64:
		str = strconv.FormatUint(v, 10)
	case string:
		str = v
	case bool:
		str = strconv.FormatBool(v)
	case []byte:
		str = string(v)
	case error:
		str = v.Error()
	case fmt.Stringer:
		str = v.String()
	default:
		newValue, _ := json.Marshal(v)
		str = string(newValue)
	}
	return str
}

func (l *stdLogger) Close() error {
	return nil
}
