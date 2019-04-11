package log

import (
	"math"
	"time"

	"github.com/bilibili/kratos/pkg/log/internal/core"
)

// D represents a map of entry level data used for structured logging.
// type D map[string]interface{}
type D = core.Field

// KVString construct Field with string value.
func KVString(key string, value string) D {
	return D{Key: key, Type: core.StringType, StringVal: value}
}

// KVInt construct Field with int value.
func KVInt(key string, value int) D {
	return D{Key: key, Type: core.IntTpye, Int64Val: int64(value)}
}

// KVInt64 construct D with int64 value.
func KVInt64(key string, value int64) D {
	return D{Key: key, Type: core.Int64Type, Int64Val: value}
}

// KVUint construct Field with uint value.
func KVUint(key string, value uint) D {
	return D{Key: key, Type: core.UintType, Int64Val: int64(value)}
}

// KVUint64 construct Field with uint64 value.
func KVUint64(key string, value uint64) D {
	return D{Key: key, Type: core.Uint64Type, Int64Val: int64(value)}
}

// KVFloat32 construct Field with float32 value.
func KVFloat32(key string, value float32) D {
	return D{Key: key, Type: core.Float32Type, Int64Val: int64(math.Float32bits(value))}
}

// KVFloat64 construct Field with float64 value.
func KVFloat64(key string, value float64) D {
	return D{Key: key, Type: core.Float64Type, Int64Val: int64(math.Float64bits(value))}
}

// KVDuration construct Field with Duration value.
func KVDuration(key string, value time.Duration) D {
	return D{Key: key, Type: core.DurationType, Int64Val: int64(value)}
}

// KV return a log kv for logging field.
// NOTE: use KV{type name} can avoid object alloc and get better performance. []~(￣▽￣)~*干杯
func KV(key string, value interface{}) D {
	return D{Key: key, Value: value}
}
