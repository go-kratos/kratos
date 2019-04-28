package paladin

import (
	"encoding"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// ErrNotExist value key not exist.
var (
	ErrNotExist       = errors.New("paladin: value key not exist")
	ErrTypeAssertion  = errors.New("paladin: value type assertion no match")
	ErrDifferentTypes = errors.New("paladin: value different types")
)

// Value is config value, maybe a json/toml/ini/string file.
type Value struct {
	val   interface{}
	slice interface{}
	raw   string
}

// Bool return bool value.
func (v *Value) Bool() (bool, error) {
	if v.val == nil {
		return false, ErrNotExist
	}
	b, ok := v.val.(bool)
	if !ok {
		return false, ErrTypeAssertion
	}
	return b, nil
}

// Int return int value.
func (v *Value) Int() (int, error) {
	i, err := v.Int64()
	if err != nil {
		return 0, nil
	}
	return int(i), nil
}

// Int32 return int32 value.
func (v *Value) Int32() (int32, error) {
	i, err := v.Int64()
	if err != nil {
		return 0, nil
	}
	return int32(i), nil
}

// Int64 return int64 value.
func (v *Value) Int64() (int64, error) {
	if v.val == nil {
		return 0, ErrNotExist
	}
	i, ok := v.val.(int64)
	if !ok {
		return 0, ErrTypeAssertion
	}
	return i, nil
}

// Float32 return float32 value.
func (v *Value) Float32() (float32, error) {
	f, err := v.Float64()
	if err != nil {
		return 0.0, err
	}
	return float32(f), nil
}

// Float64 return float64 value.
func (v *Value) Float64() (float64, error) {
	if v.val == nil {
		return 0.0, ErrNotExist
	}
	f, ok := v.val.(float64)
	if !ok {
		return 0.0, ErrTypeAssertion
	}
	return f, nil
}

// String return string value.
func (v *Value) String() (string, error) {
	if v.val == nil {
		return "", ErrNotExist
	}
	s, ok := v.val.(string)
	if !ok {
		return "", ErrTypeAssertion
	}
	return s, nil
}

// Duration parses a duration string. A duration string is a possibly signed sequence of decimal numbers
// each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (v *Value) Duration() (time.Duration, error) {
	s, err := v.String()
	if err != nil {
		return time.Duration(0), err
	}
	return time.ParseDuration(s)
}

// Raw return raw value.
func (v *Value) Raw() (string, error) {
	if v.val == nil {
		return "", ErrNotExist
	}
	return v.raw, nil
}

// Slice scan a slcie interface.
func (v *Value) Slice(dst interface{}) error {
	// NOTE: val is []interface{}, slice is []type
	if v.val == nil {
		return ErrNotExist
	}
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return ErrDifferentTypes
	}
	el := rv.Elem()
	kind := el.Type().Elem().Kind()
	if v.slice == nil {
		src, ok := v.val.([]interface{})
		if !ok {
			return ErrDifferentTypes
		}
		for _, s := range src {
			if reflect.TypeOf(s).Kind() != kind {
				return ErrTypeAssertion
			}
			el = reflect.Append(el, reflect.ValueOf(s))
		}
		v.slice = el.Interface()
		rv.Elem().Set(el)
		return nil
	}
	sv := reflect.ValueOf(v.slice)
	if sv.Type().Elem().Kind() != kind {
		return ErrTypeAssertion
	}
	rv.Elem().Set(sv)
	return nil
}

// Unmarshal is the interface implemented by an object that can unmarshal a textual representation of itself.
func (v *Value) Unmarshal(un encoding.TextUnmarshaler) error {
	text, err := v.Raw()
	if err != nil {
		return err
	}
	return un.UnmarshalText([]byte(text))
}

// UnmarshalTOML unmarhsal toml to struct.
func (v *Value) UnmarshalTOML(dst interface{}) error {
	text, err := v.Raw()
	if err != nil {
		return err
	}
	return toml.Unmarshal([]byte(text), dst)
}
