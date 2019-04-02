package core

import (
	"time"
)

// DefaultLineEnding defines the default line ending when writing logs.
// Alternate line endings specified in EncoderConfig can override this
// behavior.
const DefaultLineEnding = "\n"

// ObjectEncoder is a strongly-typed, encoding-agnostic interface for adding a
// map- or struct-like object to the logging context. Like maps, ObjectEncoders
// aren't safe for concurrent use (though typical use shouldn't require locks).
type ObjectEncoder interface {
	// Logging-specific marshalers.
	AddArray(key string, marshaler ArrayMarshaler) error
	AddObject(key string, marshaler ObjectMarshaler) error

	// Built-in types.
	AddBinary(key string, value []byte)     // for arbitrary bytes
	AddByteString(key string, value []byte) // for UTF-8 encoded bytes
	AddBool(key string, value bool)
	AddComplex128(key string, value complex128)
	AddComplex64(key string, value complex64)
	AddDuration(key string, value time.Duration)
	AddFloat64(key string, value float64)
	AddFloat32(key string, value float32)
	AddInt(key string, value int)
	AddInt64(key string, value int64)
	AddInt32(key string, value int32)
	AddInt16(key string, value int16)
	AddInt8(key string, value int8)
	AddString(key, value string)
	AddTime(key string, value time.Time)
	AddUint(key string, value uint)
	AddUint64(key string, value uint64)
	AddUint32(key string, value uint32)
	AddUint16(key string, value uint16)
	AddUint8(key string, value uint8)
	AddUintptr(key string, value uintptr)

	// AddReflected uses reflection to serialize arbitrary objects, so it's slow
	// and allocation-heavy.
	AddReflected(key string, value interface{}) error
	// OpenNamespace opens an isolated namespace where all subsequent fields will
	// be added. Applications can use namespaces to prevent key collisions when
	// injecting loggers into sub-components or third-party libraries.
	OpenNamespace(key string)
}

// ObjectMarshaler allows user-defined types to efficiently add themselves to the
// logging context, and to selectively omit information which shouldn't be
// included in logs (e.g., passwords).
type ObjectMarshaler interface {
	MarshalLogObject(ObjectEncoder) error
}

// ObjectMarshalerFunc is a type adapter that turns a function into an
// ObjectMarshaler.
type ObjectMarshalerFunc func(ObjectEncoder) error

// MarshalLogObject calls the underlying function.
func (f ObjectMarshalerFunc) MarshalLogObject(enc ObjectEncoder) error {
	return f(enc)
}

// ArrayMarshaler allows user-defined types to efficiently add themselves to the
// logging context, and to selectively omit information which shouldn't be
// included in logs (e.g., passwords).
type ArrayMarshaler interface {
	MarshalLogArray(ArrayEncoder) error
}

// ArrayMarshalerFunc is a type adapter that turns a function into an
// ArrayMarshaler.
type ArrayMarshalerFunc func(ArrayEncoder) error

// MarshalLogArray calls the underlying function.
func (f ArrayMarshalerFunc) MarshalLogArray(enc ArrayEncoder) error {
	return f(enc)
}

// ArrayEncoder is a strongly-typed, encoding-agnostic interface for adding
// array-like objects to the logging context. Of note, it supports mixed-type
// arrays even though they aren't typical in Go. Like slices, ArrayEncoders
// aren't safe for concurrent use (though typical use shouldn't require locks).
type ArrayEncoder interface {
	// Built-in types.
	PrimitiveArrayEncoder

	// Time-related types.
	AppendDuration(time.Duration)
	AppendTime(time.Time)

	// Logging-specific marshalers.
	AppendArray(ArrayMarshaler) error
	AppendObject(ObjectMarshaler) error

	// AppendReflected uses reflection to serialize arbitrary objects, so it's
	// slow and allocation-heavy.
	AppendReflected(value interface{}) error
}

// PrimitiveArrayEncoder is the subset of the ArrayEncoder interface that deals
// only in Go's built-in types. It's included only so that Duration- and
// TimeEncoders cannot trigger infinite recursion.
type PrimitiveArrayEncoder interface {
	// Built-in types.
	AppendBool(bool)
	AppendByteString([]byte) // for UTF-8 encoded bytes
	AppendComplex128(complex128)
	AppendComplex64(complex64)
	AppendFloat64(float64)
	AppendFloat32(float32)
	AppendInt(int)
	AppendInt64(int64)
	AppendInt32(int32)
	AppendInt16(int16)
	AppendInt8(int8)
	AppendString(string)
	AppendUint(uint)
	AppendUint64(uint64)
	AppendUint32(uint32)
	AppendUint16(uint16)
	AppendUint8(uint8)
	AppendUintptr(uintptr)
}

// An EncoderConfig allows users to configure the concrete encoders supplied by
// zapcore.
type EncoderConfig struct {
	EncodeTime     TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	/*EncodeLevel    LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime     TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller   CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`*/
}

// Encoder is a format-agnostic interface for all log entry marshalers. Since
// log encoders don't need to support the same wide range of use cases as
// general-purpose marshalers, it's possible to make them faster and
// lower-allocation.
//
// Implementations of the ObjectEncoder interface's methods can, of course,
// freely modify the receiver. However, the Clone and EncodeEntry methods will
// be called concurrently and shouldn't modify the receiver.
type Encoder interface {
	ObjectEncoder

	// Clone copies the encoder, ensuring that adding fields to the copy doesn't
	// affect the original.
	Clone() Encoder

	// EncodeEntry encodes an entry and fields, along with any accumulated
	// context, into a byte buffer and returns it.
	Encode(*Buffer, ...Field) error
}

// A TimeEncoder serializes a time.Time to a primitive type.
type TimeEncoder func(time.Time, PrimitiveArrayEncoder)

// A DurationEncoder serializes a time.Duration to a primitive type.
type DurationEncoder func(time.Duration, PrimitiveArrayEncoder)

// EpochTimeEncoder serializes a time.Time to a floating-point number of seconds
// since the Unix epoch.
func EpochTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) {
	//var d []byte
	enc.AppendString(t.Format("2006-01-02T15:04:05.999999"))
	//enc.AppendByteString(t.AppendFormat(d, "2006-01-02T15:04:05.999999"))
	/*nanos := t.UnixNano()
	sec := float64(nanos) / float64(time.Second)
	enc.AppendFloat64(sec)*/
}

// SecondsDurationEncoder serializes a time.Duration to a floating-point number of seconds elapsed.
func SecondsDurationEncoder(d time.Duration, enc PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Second))
}
