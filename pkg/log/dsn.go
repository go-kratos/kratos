package log

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type verboseModule map[string]int32

type logFilter []string

func (f *logFilter) String() string {
	return fmt.Sprint(*f)
}

// Set sets the value of the named command-line flag.
// format: -log.filter key1,key2
func (f *logFilter) Set(value string) error {
	for _, i := range strings.Split(value, ",") {
		*f = append(*f, strings.TrimSpace(i))
	}
	return nil
}

func (m verboseModule) String() string {
	// FIXME strings.Builder
	var buf bytes.Buffer
	for k, v := range m {
		buf.WriteString(k)
		buf.WriteString(strconv.FormatInt(int64(v), 10))
		buf.WriteString(",")
	}
	return buf.String()
}

// Set sets the value of the named command-line flag.
// format: -log.module file=1,file2=2
func (m verboseModule) Set(value string) error {
	for _, i := range strings.Split(value, ",") {
		kv := strings.Split(i, "=")
		if len(kv) == 2 {
			if v, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
				m[strings.TrimSpace(kv[0])] = int32(v)
			}
		}
	}
	return nil
}
