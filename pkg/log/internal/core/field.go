package core

import (
	"fmt"
	"math"
	"time"

	xtime "github.com/bilibili/kratos/pkg/time"
)

// FieldType represent D value type
type FieldType int32

// DType enum
const (
	UnknownType FieldType = iota
	StringType
	IntTpye
	Int64Type
	UintType
	Uint64Type
	Float32Type
	Float64Type
	DurationType
)

// Field is for encoder
type Field struct {
	Key       string
	Value     interface{}
	Type      FieldType
	StringVal string
	Int64Val  int64
}

// AddTo exports a field through the ObjectEncoder interface. It's primarily
// useful to library authors, and shouldn't be necessary in most applications.
func (f Field) AddTo(enc ObjectEncoder) {
	if f.Type == UnknownType {
		f.assertAddTo(enc)
		return
	}
	switch f.Type {
	case StringType:
		enc.AddString(f.Key, f.StringVal)
	case IntTpye:
		enc.AddInt(f.Key, int(f.Int64Val))
	case Int64Type:
		enc.AddInt64(f.Key, f.Int64Val)
	case UintType:
		enc.AddUint(f.Key, uint(f.Int64Val))
	case Uint64Type:
		enc.AddUint64(f.Key, uint64(f.Int64Val))
	case Float32Type:
		enc.AddFloat32(f.Key, math.Float32frombits(uint32(f.Int64Val)))
	case Float64Type:
		enc.AddFloat64(f.Key, math.Float64frombits(uint64(f.Int64Val)))
	case DurationType:
		enc.AddDuration(f.Key, time.Duration(f.Int64Val))
	default:
		panic(fmt.Sprintf("unknown field type: %v", f))
	}
}

func (f Field) assertAddTo(enc ObjectEncoder) {
	// assert interface
	switch val := f.Value.(type) {
	case bool:
		enc.AddBool(f.Key, val)
	case complex128:
		enc.AddComplex128(f.Key, val)
	case complex64:
		enc.AddComplex64(f.Key, val)
	case float64:
		enc.AddFloat64(f.Key, val)
	case float32:
		enc.AddFloat32(f.Key, val)
	case int:
		enc.AddInt(f.Key, val)
	case int64:
		enc.AddInt64(f.Key, val)
	case int32:
		enc.AddInt32(f.Key, val)
	case int16:
		enc.AddInt16(f.Key, val)
	case int8:
		enc.AddInt8(f.Key, val)
	case string:
		enc.AddString(f.Key, val)
	case uint:
		enc.AddUint(f.Key, val)
	case uint64:
		enc.AddUint64(f.Key, val)
	case uint32:
		enc.AddUint32(f.Key, val)
	case uint16:
		enc.AddUint16(f.Key, val)
	case uint8:
		enc.AddUint8(f.Key, val)
	case []byte:
		enc.AddByteString(f.Key, val)
	case uintptr:
		enc.AddUintptr(f.Key, val)
	case time.Time:
		enc.AddTime(f.Key, val)
	case xtime.Time:
		enc.AddTime(f.Key, val.Time())
	case time.Duration:
		enc.AddDuration(f.Key, val)
	case xtime.Duration:
		enc.AddDuration(f.Key, time.Duration(val))
	case error:
		enc.AddString(f.Key, val.Error())
	case fmt.Stringer:
		enc.AddString(f.Key, val.String())
	default:
		err := enc.AddReflected(f.Key, val)
		if err != nil {
			enc.AddString(fmt.Sprintf("%sError", f.Key), err.Error())
		}
	}
}
