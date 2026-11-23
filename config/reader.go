package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/go-kratos/kratos/v3/log"
)

// Reader is config reader.
type Reader interface {
	Merge(...*KeyValue) error
	Value(string) (Value, bool)
	Source() ([]byte, error)
	Resolve() error
}

type reader struct {
	opts   options
	values map[string]any
	lock   sync.Mutex
}

func newReader(opts options) Reader {
	return &reader{
		opts:   opts,
		values: make(map[string]any),
		lock:   sync.Mutex{},
	}
}

func (r *reader) Merge(kvs ...*KeyValue) error {
	merged, err := r.cloneMap()
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		next := make(map[string]any)
		if err := r.opts.decoder(kv, next); err != nil {
			log.Error("failed to decode config", "error", err, "key", kv.Key, "value", string(kv.Value))
			return err
		}
		if err := r.opts.merge(&merged, convertMap(next)); err != nil {
			log.Error("failed to merge config", "error", err, "key", kv.Key, "value", string(kv.Value))
			return err
		}
	}
	r.lock.Lock()
	r.values = merged
	r.lock.Unlock()
	return nil
}

func (r *reader) Value(path string) (Value, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return readValue(r.values, path)
}

func (r *reader) Source() ([]byte, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return marshalJSON(convertMap(r.values))
}

func (r *reader) Resolve() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.opts.resolver(r.values)
}

func (r *reader) cloneMap() (map[string]any, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return cloneMap(r.values)
}

func cloneMap(src map[string]any) (map[string]any, error) {
	// https://gist.github.com/soroushjp/0ec92102641ddfc3ad5515ca76405f4d
	var buf bytes.Buffer
	gob.Register(map[string]any{})
	gob.Register([]any{})
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(src)
	if err != nil {
		return nil, err
	}
	var clone map[string]any
	err = dec.Decode(&clone)
	if err != nil {
		return nil, err
	}
	return clone, nil
}

func convertMap(src any) any {
	switch m := src.(type) {
	case map[string]any:
		dst := make(map[string]any, len(m))
		for k, v := range m {
			dst[k] = convertMap(v)
		}
		return dst
	case map[any]any:
		dst := make(map[string]any, len(m))
		for k, v := range m {
			dst[fmt.Sprint(k)] = convertMap(v)
		}
		return dst
	case []any:
		dst := make([]any, len(m))
		for k, v := range m {
			dst[k] = convertMap(v)
		}
		return dst
	case []byte:
		// there will be no binary data in the config data
		return string(m)
	default:
		return src
	}
}

// readValue read Value in given map[string]interface{}
// by the given path, will return false if not found.
func readValue(values map[string]any, path string) (Value, bool) {
	var (
		next = values
		keys = strings.Split(path, ".")
		last = len(keys) - 1
	)
	for idx, key := range keys {
		value, ok := next[key]
		if !ok {
			return nil, false
		}
		if idx == last {
			av := &atomicValue{}
			av.Store(value)
			return av, true
		}
		switch vm := value.(type) {
		case map[string]any:
			next = vm
		default:
			return nil, false
		}
	}
	return nil, false
}

func marshalJSON(v any) ([]byte, error) {
	if m, ok := v.(proto.Message); ok {
		return protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(m)
	}
	return json.Marshal(v)
}

func unmarshalJSON(data []byte, v any) error {
	if m, ok := v.(proto.Message); ok {
		return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, m)
	}

	// For struct pointers, use reflection-based scanning that converts
	// string values to the target field's type. This enables env variable
	// overrides (which are always strings) to work with bool, int, and
	// float fields, while preserving string fields that look numeric
	// (e.g. passwords, version numbers).
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() && rv.Elem().Kind() == reflect.Struct {
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		return scanStruct(rv.Elem(), m)
	}

	return json.Unmarshal(data, v)
}

// scanStruct scans a map into a struct value using reflection.
// String values are coerced to the target field's type (bool/int/uint/float/string).
// This is type-safe: a string field will always remain a string, even if its
// value looks like a number or boolean.
func scanStruct(v reflect.Value, data map[string]any) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		// Handle embedded (anonymous) struct fields by scanning their
		// fields against the same data map, mirroring json.Unmarshal behavior.
		if field.Anonymous {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					continue
				}
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				if err := scanStruct(fv, data); err != nil {
					return err
				}
			}
			continue
		}

		key := field.Tag.Get("json")
		if key == "-" {
			continue
		}
		if idx := strings.IndexByte(key, ','); idx != -1 {
			key = key[:idx]
		}
		if key == "" {
			key = field.Name
		}

		val, ok := data[key]
		if !ok {
			continue
		}

		if err := setStructField(v.Field(i), val); err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}
	return nil
}

// setStructField sets a struct field from a map value, with type coercion
// for scalar types.
func setStructField(field reflect.Value, val any) error {
	// Handle nil values for pointer fields
	if val == nil {
		return nil
	}

	// Handle pointer fields
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setStructField(field.Elem(), val)
	}

	switch field.Kind() {
	case reflect.Bool:
		b, err := valueToBool(val)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := valueToInt(val)
		if err != nil {
			return err
		}
		field.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := valueToUint(val)
		if err != nil {
			return err
		}
		field.SetUint(n)

	case reflect.Float32, reflect.Float64:
		f, err := valueToFloat(val)
		if err != nil {
			return err
		}
		field.SetFloat(f)

	case reflect.String:
		s, err := valueToString(val)
		if err != nil {
			return err
		}
		field.SetString(s)

	case reflect.Struct:
		m, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object for struct field, got %T", val)
		}
		return scanStruct(field, m)

	case reflect.Slice, reflect.Map, reflect.Interface:
		// For complex types, fall back to JSON round-trip.
		// These are rarely overridden via env vars, and the standard
		// JSON unmarshaler handles them correctly in the common case.
		data, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, field.Addr().Interface())

	default:
		// Fall back to JSON round-trip for other types.
		data, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, field.Addr().Interface())
	}
	return nil
}

func valueToBool(val any) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case float64:
		return v != 0, nil
	}
	return false, fmt.Errorf("cannot convert %T to bool", val)
}

func valueToInt(val any) (int64, error) {
	switch v := val.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	}
	return 0, fmt.Errorf("cannot convert %T to int64", val)
}

func valueToUint(val any) (uint64, error) {
	switch v := val.(type) {
	case float64:
		return uint64(v), nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	}
	return 0, fmt.Errorf("cannot convert %T to uint64", val)
}

func valueToFloat(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	}
	return 0, fmt.Errorf("cannot convert %T to float64", val)
}

func valueToString(val any) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	}
	return "", fmt.Errorf("cannot convert %T to string", val)
}
