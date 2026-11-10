package config

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// coerceScalarTypes converts string leaf values in the serialized config to the
// scalar types declared by the destination (v), so that string-only sources
// (e.g. env) can populate bool/int/float fields during Scan.
//
// The coercion is driven exclusively by the destination's declared type — never
// by guessing the type from the value's text shape — so string fields whose
// values happen to look numeric ("1234", "0.5", "8000") are left untouched.
//
// If no string leaf needed conversion the original bytes are returned unchanged
// (the common case, so scanning is a pure round-trip of the source JSON).
func coerceScalarTypes(data []byte, v any) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var root any
	if err := dec.Decode(&root); err != nil {
		return data, err
	}

	var changed bool
	if pb, ok := v.(proto.Message); ok {
		root, changed = coerceProtoTree(root, pb.ProtoReflect().Descriptor())
	} else {
		root, changed = coerceStructTree(root, reflect.TypeOf(v))
	}
	if !changed {
		return data, nil
	}
	out, err := json.Marshal(root)
	if err != nil {
		return data, err
	}
	return out, nil
}

// coerceStructTree walks a JSON tree (decoded with UseNumber) alongside the
// destination's reflect.Type, converting string leaves to the declared scalar
// kind. Returns the (possibly rewritten) tree and whether any leaf changed.
func coerceStructTree(v any, t reflect.Type) (any, bool) {
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t == nil {
		return v, false
	}
	switch t.Kind() {
	case reflect.Struct:
		m, ok := v.(map[string]any)
		if !ok {
			return v, false
		}
		changed := false
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" { // unexported
				continue
			}
			for _, name := range jsonFieldNames(f) {
				sub, ok := m[name]
				if !ok {
					continue
				}
				nv, ch := coerceStructTree(sub, f.Type)
				if ch {
					m[name] = nv
					changed = true
				}
				break
			}
		}
		return m, changed
	case reflect.Slice, reflect.Array:
		arr, ok := v.([]any)
		if !ok {
			return v, false
		}
		changed := false
		for i := range arr {
			nv, ch := coerceStructTree(arr[i], t.Elem())
			if ch {
				arr[i] = nv
				changed = true
			}
		}
		return arr, changed
	case reflect.Map:
		m, ok := v.(map[string]any)
		if !ok {
			return v, false
		}
		changed := false
		for k, sub := range m {
			nv, ch := coerceStructTree(sub, t.Elem())
			if ch {
				m[k] = nv
				changed = true
			}
		}
		return m, changed
	default:
		return convertStringLeaf(v, t)
	}
}

// convertStringLeaf converts a string leaf to the destination scalar kind.
// Returns (converted, true) only when the leaf is a string AND the destination
// kind is bool or numeric AND the string parses. It is a no-op for string
// destinations, so numeric-looking strings stay strings.
func convertStringLeaf(v any, t reflect.Type) (any, bool) {
	s, ok := v.(string)
	if !ok {
		return v, false
	}
	switch t.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return v, false
		}
		return b, true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, t.Bits())
		if err != nil {
			return v, false
		}
		return i, true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, t.Bits())
		if err != nil {
			return v, false
		}
		return u, true
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return v, false
		}
		return f, true
	default:
		return v, false
	}
}

// coerceProtoTree walks a JSON tree alongside a protobuf message descriptor,
// converting string leaves to the field's declared kind.
func coerceProtoTree(v any, md protoreflect.MessageDescriptor) (any, bool) {
	m, ok := v.(map[string]any)
	if !ok {
		return v, false
	}
	changed := false
	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		for _, name := range protoFieldNames(fd) {
			sub, ok := m[name]
			if !ok {
				continue
			}
			nv, ch := coerceProtoField(sub, fd)
			if ch {
				m[name] = nv
				changed = true
			}
			break
		}
	}
	return m, changed
}

func coerceProtoField(v any, fd protoreflect.FieldDescriptor) (any, bool) {
	if fd.IsMap() {
		m, ok := v.(map[string]any)
		if !ok {
			return v, false
		}
		changed := false
		for k, sub := range m {
			valDesc := fd.MapValue()
			if valDesc.Kind() == protoreflect.MessageKind || valDesc.Kind() == protoreflect.GroupKind {
				nv, ch := coerceProtoTree(sub, valDesc.Message())
				if ch {
					m[k] = nv
					changed = true
				}
				continue
			}
			nv, ch := convertStringLeaf(sub, kindToType(valDesc.Kind()))
			if ch {
				m[k] = nv
				changed = true
			}
		}
		return m, changed
	}
	if fd.IsList() {
		arr, ok := v.([]any)
		if !ok {
			return v, false
		}
		changed := false
		for i := range arr {
			nv, ch := coerceProtoField(arr[i], fd)
			if ch {
				arr[i] = nv
				changed = true
			}
		}
		return arr, changed
	}
	if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
		return coerceProtoTree(v, fd.Message())
	}
	return convertStringLeaf(v, kindToType(fd.Kind()))
}

// kindToType maps a protobuf scalar kind onto a Go reflect.Type usable by
// convertStringLeaf. Message/group kinds are not expected here.
func kindToType(k protoreflect.Kind) reflect.Type {
	switch k {
	case protoreflect.BoolKind:
		return reflect.TypeOf(false)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return reflect.TypeOf(int32(0))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return reflect.TypeOf(int64(0))
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return reflect.TypeOf(uint32(0))
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return reflect.TypeOf(uint64(0))
	case protoreflect.FloatKind:
		return reflect.TypeOf(float32(0))
	case protoreflect.DoubleKind:
		return reflect.TypeOf(float64(0))
	default:
		return reflect.TypeOf("")
	}
}

// jsonFieldNames returns the JSON key candidates for a struct field: the json
// tag name if present, otherwise the field name (encoding/json matches field
// names case-insensitively).
func jsonFieldNames(f reflect.StructField) []string {
	tag := f.Tag.Get("json")
	if tag != "" {
		// "name,omitempty" or "name"
		if i := indexByte(tag, ','); i >= 0 {
			tag = tag[:i]
		}
		if tag == "-" {
			return nil
		}
		if tag != "" {
			return []string{tag}
		}
	}
	return []string{f.Name, toLowerFirst(f.Name)}
}

// protoFieldNames returns the JSON key candidates for a protobuf field:
// protojson accepts both the json_name and the original proto field name.
func protoFieldNames(fd protoreflect.FieldDescriptor) []string {
	return []string{fd.JSONName(), string(fd.Name())}
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'A' && b[0] <= 'Z' {
		b[0] = b[0] + ('a' - 'A')
	}
	return string(b)
}
