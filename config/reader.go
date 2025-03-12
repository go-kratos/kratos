package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/go-kratos/kratos/v2/log"
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
			log.Errorf("Failed to config decode error: %v key: %s value: %s", err, kv.Key, string(kv.Value))
			return err
		}
		if err := r.opts.merge(&merged, convertMap(next)); err != nil {
			log.Errorf("Failed to config merge error: %v key: %s value: %s", err, kv.Key, string(kv.Value))
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
	return json.Unmarshal(data, v)
}
