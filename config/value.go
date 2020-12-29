package config

import (
	"encoding/json"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

var (
	_ Value = (*jsonValue)(nil)
	_ Value = (*errValue)(nil)
)

// Value is config value interface.
type Value interface {
	Bool() (bool, error)
	Int() (int, error)
	Int32() (int32, error)
	Int64() (int64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	String() (string, error)
	Duration() (time.Duration, error)
	Scan(interface{}) error
}

type jsonValue struct {
	raw *simplejson.Json
}

func (v jsonValue) Bool() (bool, error) { return v.raw.Bool() }
func (v jsonValue) Int() (int, error)   { return v.raw.Int() }
func (v jsonValue) Int32() (int32, error) {
	val, err := v.raw.Int()
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}
func (v jsonValue) Int64() (int64, error) { return v.raw.Int64() }
func (v jsonValue) Float32() (float32, error) {
	val, err := v.raw.Float64()
	if err != nil {
		return 0.0, err
	}
	return float32(val), nil
}
func (v jsonValue) Float64() (float64, error) { return v.raw.Float64() }
func (v jsonValue) String() (string, error)   { return v.raw.String() }
func (v jsonValue) Duration() (time.Duration, error) {
	intVal, err := v.raw.Int64()
	if err != nil {
		return 0, err
	}
	return time.Duration(intVal), nil
}
func (v jsonValue) Scan(obj interface{}) error {
	data, err := v.raw.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

type errValue struct {
	err error
}

func (v errValue) Bool() (bool, error)              { return false, v.err }
func (v errValue) Int() (int, error)                { return 0, v.err }
func (v errValue) Int32() (int32, error)            { return 0, v.err }
func (v errValue) Int64() (int64, error)            { return 0, v.err }
func (v errValue) Float32() (float32, error)        { return 0.0, v.err }
func (v errValue) Float64() (float64, error)        { return 0.0, v.err }
func (v errValue) Duration() (time.Duration, error) { return 0, v.err }
func (v errValue) String() (string, error)          { return "", v.err }
func (v errValue) Scan(interface{}) error           { return v.err }
