package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
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
	layers []map[string]any
	lock   sync.Mutex
	// Serialize writes separately so merge callbacks can read published values.
	updateLock sync.Mutex
}

// deletedConfigValue records a resolver deletion until layers are replayed.
type deletedConfigValue struct{}

func newReader(opts options) Reader {
	return &reader{
		opts:   opts,
		values: make(map[string]any),
		lock:   sync.Mutex{},
	}
}

func (r *reader) Merge(kvs ...*KeyValue) error {
	// Serialize with other merges and resolution without holding the data lock
	// across decoder or merge callbacks, which may read the current config.
	r.updateLock.Lock()
	defer r.updateLock.Unlock()
	r.lock.Lock()
	merged, err := cloneMap(r.values)
	r.lock.Unlock()
	if err != nil {
		return err
	}
	layers := make([]map[string]any, 0, len(kvs))
	for _, kv := range kvs {
		next := make(map[string]any)
		if err := r.opts.decoder(kv, next); err != nil {
			log.Error("failed to decode config", "error", err, "key", kv.Key, "value", string(kv.Value))
			return err
		}
		source, ok := convertMap(next).(map[string]any)
		if !ok {
			return fmt.Errorf("config: decoded source must be map[string]interface{}, got %T", next)
		}
		// Custom mergers cannot be replayed per layer, so they keep only the
		// merged view and Scan falls back to Source.
		if !r.opts.customMerge {
			layers = append(layers, cloneMergeValue(source).(map[string]any))
		}
		if err := r.opts.merge(&merged, source); err != nil {
			log.Error("failed to merge config", "error", err, "key", kv.Key, "value", string(kv.Value))
			return err
		}
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.values = merged
	if len(layers) > 0 {
		r.layers = compactLayers(append(r.layers, layers...))
	}
	return nil
}

// compactLayers drops historical values that have since been replaced by a
// later source. Different spellings are kept because their destination field
// is not known until Scan.
func compactLayers(layers []map[string]any) []map[string]any {
	if len(layers) <= 1 {
		return layers
	}
	suffix := make(map[string]any)
	retained := make([]map[string]any, 0, len(layers))
	for i := len(layers) - 1; i >= 0; i-- {
		layer := filterOverridden(layers[i], suffix)
		if len(layer) > 0 {
			retained = append(retained, layer)
		}
		base := cloneMergeValue(layers[i]).(map[string]any)
		mergeMap(base, suffix)
		suffix = base
	}
	for left, right := 0, len(retained)-1; left < right; left, right = left+1, right-1 {
		retained[left], retained[right] = retained[right], retained[left]
	}
	return retained
}

func filterOverridden(current, later map[string]any) map[string]any {
	filtered := make(map[string]any, len(current))
	for key, value := range current {
		laterValue, overridden := later[key]
		if !overridden {
			filtered[key] = cloneMergeValue(value)
			continue
		}
		currentMap, currentIsMap := value.(map[string]any)
		laterMap, laterIsMap := laterValue.(map[string]any)
		if !currentIsMap || !laterIsMap {
			continue
		}
		if nested := filterOverridden(currentMap, laterMap); len(nested) > 0 {
			filtered[key] = nested
		}
	}
	return filtered
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
	r.updateLock.Lock()
	defer r.updateLock.Unlock()
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.opts.customMerge {
		return r.opts.resolver(r.values)
	}
	before := cloneMergeValue(r.values).(map[string]any)
	err := r.opts.resolver(r.values)
	// Write the resolver's edits back to the layers that own the original paths,
	// so Scan keeps source precedence without rerunning the resolver.
	if changes := resolvedChanges(before, r.values); len(changes) > 0 {
		for i := len(r.layers) - 1; i >= 0; i-- {
			mergeMap(r.layers[i], takeResolvedChanges(changes, r.layers[i]))
		}
		if len(changes) > 0 {
			r.layers = append(r.layers, changes)
		}
		r.layers = compactLayers(r.layers)
	}
	return err
}

// takeResolvedChanges assigns edits to the latest layer containing their
// original path. Keeping these edits in place preserves precedence when an
// earlier spelling contains a placeholder and a later alias is already literal.
func takeResolvedChanges(changes, source map[string]any) map[string]any {
	selected := make(map[string]any)
	for key, value := range changes {
		original, exists := source[key]
		if !exists {
			continue
		}
		changeMap, isMap := value.(map[string]any)
		sourceMap, sourceIsMap := original.(map[string]any)
		if isMap && sourceIsMap && len(changeMap) > 0 {
			if nested := takeResolvedChanges(changeMap, sourceMap); len(nested) > 0 {
				selected[key] = nested
			}
			if len(changeMap) == 0 {
				delete(changes, key)
			}
		} else {
			selected[key] = value
			delete(changes, key)
		}
	}
	return selected
}

// resolvedChanges records only the resolver's edits, including removed keys.
func resolvedChanges(before, after map[string]any) map[string]any {
	changes := make(map[string]any)
	for key, old := range before {
		next, exists := after[key]
		if !exists {
			changes[key] = deletedConfigValue{}
			continue
		}
		oldMap, oldIsMap := old.(map[string]any)
		nextMap, nextIsMap := next.(map[string]any)
		if oldIsMap && nextIsMap {
			if nested := resolvedChanges(oldMap, nextMap); len(nested) > 0 {
				changes[key] = nested
			}
		} else if !reflect.DeepEqual(old, next) {
			changes[key] = cloneMergeValue(next)
		}
	}
	for key, next := range after {
		if _, exists := before[key]; !exists {
			changes[key] = cloneMergeValue(next)
		}
	}
	return changes
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
