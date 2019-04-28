package common

import (
	"bytes"
	"fmt"
	"crypto/md5"
	"io/ioutil"
)

// importantLog check if log level is above ERROR
func CriticalLog(level []byte) bool {
	if bytes.Equal(level, []byte("WARN")) || bytes.Equal(level, []byte("ERROR")) || bytes.Equal(level, []byte("FATAL")) {
		return true
	}
	return false
}

// GetPriority get priority value from json body
func GetPriority(logBody []byte) (value []byte, err error) {
	return SeekValue([]byte(`"priority":`), logBody)
}

// seekValue seek value by key from json
func SeekValue(key []byte, logBody []byte) (value []byte, err error) {
	var (
		b, logLen, begin, end int
	)
	b = bytes.Index(logBody, key)
	if b != -1 {
		logLen = len(logBody)
		for begin = b + len(key); begin < logLen && logBody[begin] != byte('"'); begin++ {
		}
		if begin >= logLen {
			err = fmt.Errorf("beginning of value not found by key: %s", string(key))
			return
		}
		begin++ // begin position of value of appid
		for end = begin; end < logLen && logBody[end] != byte('"'); end++ {
		}
		if end >= logLen {
			err = fmt.Errorf("end of value not found by key: %s", string(key))
			return
		}
		value = logBody[begin:end]
		return
	} else {
		err = fmt.Errorf("key %s not found", string(key))
		return
	}
}

func FileMd5(filePath string) string {
	data, _ := ioutil.ReadFile(filePath)
	value := md5.Sum(data)
	return fmt.Sprintf("%x", value)
}
