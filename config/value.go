package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var (
	_ Value = (*jsonValue)(nil)
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
}

type jsonValue struct {
	raw interface{}
}

func (v *jsonValue) Map() (map[string]interface{}, error) {
	if m, ok := (v.raw).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, ErrTypeAssert
}
func (v *jsonValue) Get(key string) *jsonValue {
	if m, err := v.Map(); err == nil {
		if val, ok := m[key]; ok {
			return &jsonValue{raw: val}
		}
	}
	return &jsonValue{raw: nil}
}
func (v *jsonValue) GetPath(path ...string) *jsonValue {
	var next = v
	for _, key := range path {
		next = next.Get(key)
	}
	return next
}
func (v *jsonValue) Bool() (bool, error) {
	switch val := v.raw.(type) {
	case bool:
		return val, nil
	case int64, float64, string:
		return strconv.ParseBool(fmt.Sprint(val))
	}
	return false, fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.raw))
}
func (v *jsonValue) Int() (int64, error) {
	switch val := v.raw.(type) {
	case int64:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.raw))
}
func (v *jsonValue) Float() (float64, error) {
	switch val := v.raw.(type) {
	case float64:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 10)
	}
	return 0.0, fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.raw))
}
func (v *jsonValue) String() (string, error) {
	switch val := v.raw.(type) {
	case string:
		return val, nil
	case bool, int64, float64:
		return fmt.Sprint(v.raw), nil
	}
	return "", fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.raw))
}
func (v *jsonValue) Duration() (time.Duration, error) {
	val, err := v.Int()
	if err != nil {
		return 0, err
	}
	return time.Duration(val), nil
}
func (v *jsonValue) Scan(obj interface{}) error {
	data, err := json.Marshal(v.raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

type errValue struct {
	err error
}

func (v errValue) Bool() (bool, error)              { return false, v.err }
func (v errValue) Int() (int64, error)              { return 0, v.err }
func (v errValue) Float() (float64, error)          { return 0.0, v.err }
func (v errValue) Duration() (time.Duration, error) { return 0, v.err }
func (v errValue) String() (string, error)          { return "", v.err }
func (v errValue) Scan(interface{}) error           { return v.err }
