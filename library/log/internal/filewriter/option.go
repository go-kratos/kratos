package filewriter

import (
	"fmt"
	"strings"
	"time"
)

// RotateFormat
const (
	RotateDaily = "2006-01-02"
)

var defaultOption = option{
	RotateFormat:   RotateDaily,
	MaxSize:        1 << 30,
	ChanSize:       1024 * 8,
	RotateInterval: 10 * time.Second,
}

type option struct {
	RotateFormat string
	MaxFile      int
	MaxSize      int64
	ChanSize     int

	// TODO export Option
	RotateInterval time.Duration
	WriteTimeout   time.Duration
}

// Option filewriter option
type Option func(opt *option)

// RotateFormat e.g 2006-01-02 meaning rotate log file every day.
// NOTE: format can't contain ".", "." will cause panic ヽ(*。>Д<)o゜.
func RotateFormat(format string) Option {
	if strings.Contains(format, ".") {
		panic(fmt.Sprintf("rotate format can't contain '.' format: %s", format))
	}
	return func(opt *option) {
		opt.RotateFormat = format
	}
}

// MaxFile default 999, 0 meaning unlimit.
// TODO: don't create file list if MaxSize is unlimt.
func MaxFile(n int) Option {
	return func(opt *option) {
		opt.MaxFile = n
	}
}

// MaxSize set max size for single log file,
// defult 1GB, 0 meaning unlimit.
func MaxSize(n int64) Option {
	return func(opt *option) {
		opt.MaxSize = n
	}
}

// ChanSize set internal chan size default 8192 use about 64k memory on x64 platfrom static,
// because filewriter has internal object pool, change chan size bigger may cause filewriter use
// a lot of memory, because sync.Pool can't set expire time memory won't free until program exit.
func ChanSize(n int) Option {
	return func(opt *option) {
		opt.ChanSize = n
	}
}
