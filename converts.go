package kratos

import (
	"fmt"
	"strconv"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// String converts.
func String(val string) (string, error) {
	return val, nil
}

// Bool converts.
func Bool(val string) (bool, error) {
	return strconv.ParseBool(val)
}

// Float64 converts.
func Float64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

// Float32 converts.
func Float32(val string) (float32, error) {
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

// Int64 converts.
func Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 0, 64)
}

// Int32 converts.
func Int32(val string) (int32, error) {
	i, err := strconv.ParseInt(val, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// Uint64 converts.
func Uint64(val string) (uint64, error) {
	return strconv.ParseUint(val, 0, 64)
}

// Uint32 converts.
func Uint32(val string) (uint32, error) {
	i, err := strconv.ParseUint(val, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

// Timestamp converts.
func Timestamp(val string) (*timestamppb.Timestamp, error) {
	var r timestamppb.Timestamp
	unmarshaler := &protojson.UnmarshalOptions{}
	err := unmarshaler.Unmarshal([]byte(val), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Duration converts.
func Duration(val string) (*durationpb.Duration, error) {
	var r durationpb.Duration
	unmarshaler := &protojson.UnmarshalOptions{}
	err := unmarshaler.Unmarshal([]byte(val), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Enum converts.
func Enum(val string, enumValMap map[string]int32) (int32, error) {
	e, ok := enumValMap[val]
	if ok {
		return e, nil
	}
	i, err := Int32(val)
	if err != nil {
		return 0, fmt.Errorf("invalid value: %s", val)
	}
	for _, v := range enumValMap {
		if v == i {
			return i, nil
		}
	}
	return 0, fmt.Errorf("invalid value: %s", val)
}
