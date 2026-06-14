package zerolog

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	zerolog.Logger
}

func NewLogger(zlog zerolog.Logger) *Logger {
	return &Logger{zlog}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	keylen := len(keyvals)
	if keylen == 0 {
		return nil
	}
	if (keylen & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
		if level < log.LevelWarn {
			level = log.LevelWarn
		}
	}

	ev := l.WithLevel(zerolog.Level(level + 1))

	for i := 0; i < keylen; i += 2 {
		var key string
		switch kk := keyvals[i].(type) {
		case string:
			key = kk
		default:
			key = fmt.Sprintf("%v", kk)
		}

		switch val := keyvals[i+1].(type) {
		case string:
			ev = ev.Str(key, val)
		case []byte:
			ev = ev.Bytes(key, val)
		case error:
			ev = ev.AnErr(key, val)
		case []error:
			ev = ev.Errs(key, val)
		case bool:
			ev = ev.Bool(key, val)
		case int:
			ev = ev.Int(key, val)
		case int8:
			ev = ev.Int8(key, val)
		case int16:
			ev = ev.Int16(key, val)
		case int32:
			ev = ev.Int32(key, val)
		case int64:
			ev = ev.Int64(key, val)
		case uint:
			ev = ev.Uint(key, val)
		case uint8:
			ev = ev.Uint8(key, val)
		case uint16:
			ev = ev.Uint16(key, val)
		case uint32:
			ev = ev.Uint32(key, val)
		case uint64:
			ev = ev.Uint64(key, val)
		case float32:
			ev = ev.Float32(key, val)
		case float64:
			ev = ev.Float64(key, val)
		case time.Time:
			ev = ev.Time(key, val)
		case time.Duration:
			ev = ev.Dur(key, val)
		case *string:
			if val != nil {
				ev = ev.Str(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *bool:
			if val != nil {
				ev = ev.Bool(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *int:
			if val != nil {
				ev = ev.Int(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *int8:
			if val != nil {
				ev = ev.Int8(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *int16:
			if val != nil {
				ev = ev.Int16(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *int32:
			if val != nil {
				ev = ev.Int32(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *int64:
			if val != nil {
				ev = ev.Int64(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *uint:
			if val != nil {
				ev = ev.Uint(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *uint8:
			if val != nil {
				ev = ev.Uint8(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *uint16:
			if val != nil {
				ev = ev.Uint16(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *uint32:
			if val != nil {
				ev = ev.Uint32(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *uint64:
			if val != nil {
				ev = ev.Uint64(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *float32:
			if val != nil {
				ev = ev.Float32(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *float64:
			if val != nil {
				ev = ev.Float64(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *time.Time:
			if val != nil {
				ev = ev.Time(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case *time.Duration:
			if val != nil {
				ev = ev.Dur(key, *val)
			} else {
				ev = ev.Str(key, "nil")
			}
		case []string:
			ev = ev.Strs(key, val)
		case []bool:
			ev = ev.Bools(key, val)
		case []int:
			ev = ev.Ints(key, val)
		case []int8:
			ev = ev.Ints8(key, val)
		case []int16:
			ev = ev.Ints16(key, val)
		case []int32:
			ev = ev.Ints32(key, val)
		case []int64:
			ev = ev.Ints64(key, val)
		case []uint:
			ev = ev.Uints(key, val)
		// case []uint8:
		// 	dst = enc.AppendUints8(dst, val)
		case []uint16:
			ev = ev.Uints16(key, val)
		case []uint32:
			ev = ev.Uints32(key, val)
		case []uint64:
			ev = ev.Uints64(key, val)
		case []float32:
			ev = ev.Floats32(key, val)
		case []float64:
			ev = ev.Floats64(key, val)
		case []time.Time:
			ev = ev.Times(key, val)
		case []time.Duration:
			ev = ev.Durs(key, val)
		case nil:
			ev = ev.Str(key, "nil")
		case net.IP:
			ev = ev.IPAddr(key, val)
		case net.IPNet:
			ev = ev.IPPrefix(key, val)
		case net.HardwareAddr:
			ev = ev.MACAddr(key, val)
		case json.RawMessage:
			ev = ev.RawJSON(key, val)
		default:
			if strer, ok := keyvals[i].(fmt.Stringer); ok {
				ev = ev.Stringer(key, strer)
				continue
			}
			ev = ev.Interface(key, val)
		}
	}
	ev.Send()
	return nil
}
