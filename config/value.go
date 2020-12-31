package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	_ Value = (*atomicValue)(nil)
	_ Value = (*errValue)(nil)
)

// Value is config value interface.
type Value interface {
	Bool() (bool, error)
	Int() (int64, error)
	Float() (float64, error)
	String() (string, error)
	Duration() (time.Duration, error)
	Scan(interface{}) error
	Map() (map[string]interface{}, error)
}

type atomicValue struct {
	raw atomic.Value
}

func (v *atomicValue) Bool() (bool, error) {
	switch val := v.raw.Load().(type) {
	case bool:
		return val, nil
	case int64, float64, string:
		return strconv.ParseBool(fmt.Sprint(val))
	}
	return false, fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.raw))
}
func (v *atomicValue) Int() (int64, error) {
	switch val := v.raw.Load().(type) {
	case int64:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.raw))
}
func (v *atomicValue) Float() (float64, error) {
	switch val := v.raw.Load().(type) {
	case float64:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 10)
	}
	return 0.0, fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.raw))
}
func (v *atomicValue) String() (string, error) {
	switch val := v.raw.Load().(type) {
	case string:
		return val, nil
	case bool, int64, float64:
		return fmt.Sprint(v.raw), nil
	}
	return "", fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.raw))
}
func (v *atomicValue) Duration() (time.Duration, error) {
	val, err := v.Int()
	if err != nil {
		return 0, err
	}
	return time.Duration(val), nil
}
func (v *atomicValue) Scan(obj interface{}) error {
	data, err := json.Marshal(v.raw.Load())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}
func (v *atomicValue) Map() (map[string]interface{}, error) {
	raw := v.raw.Load()
	if raw == nil {
		return nil, ErrNotFound
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, ErrTypeAssert
	}
	return m, nil
}

type errValue struct {
	err error
}

func (v errValue) Bool() (bool, error)                  { return false, v.err }
func (v errValue) Int() (int64, error)                  { return 0, v.err }
func (v errValue) Float() (float64, error)              { return 0.0, v.err }
func (v errValue) Duration() (time.Duration, error)     { return 0, v.err }
func (v errValue) String() (string, error)              { return "", v.err }
func (v errValue) Scan(interface{}) error               { return v.err }
func (v errValue) Map() (map[string]interface{}, error) { return nil, v.err }
