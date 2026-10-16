package config

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// SourceFor returns the merged source after resolving field aliases from the
// destination type. Sources are merged before JSON decoding so an overlay
// using a proto field's JSON name can replace one using its proto name.
func (r *reader) SourceFor(v any) ([]byte, error) {
	message, ok := v.(proto.Message)
	// Non-protobuf destinations and custom mergers keep the merged Source view.
	if !ok || r.opts.customMerge {
		return r.Source()
	}
	r.lock.Lock()
	if len(r.layers) == 0 {
		r.lock.Unlock()
		return r.Source()
	}
	layers := make([]map[string]any, len(r.layers))
	for i, layer := range r.layers {
		layers[i] = cloneMergeValue(layer).(map[string]any)
	}
	r.lock.Unlock()

	desc := message.ProtoReflect().Descriptor()
	merged := make(map[string]any)
	for _, layer := range layers {
		canonical, err := canonicalizeProtoMap(layer, desc)
		if err != nil {
			return nil, err
		}
		mergeMap(merged, canonical)
	}
	removeDeletedValues(merged)
	return marshalJSON(merged)
}

func removeDeletedValues(values map[string]any) {
	for key, value := range values {
		switch v := value.(type) {
		case deletedConfigValue:
			delete(values, key)
		case map[string]any:
			removeDeletedValues(v)
		}
	}
}

func canonicalizeProtoMap(src map[string]any, desc protoreflect.MessageDescriptor) (map[string]any, error) {
	// Only types with special protojson mappings bypass field canonicalization.
	// Other google.protobuf messages, such as Type, use ordinary object fields.
	switch desc.FullName() {
	case "google.protobuf.Any", "google.protobuf.Timestamp", "google.protobuf.Duration",
		"google.protobuf.BoolValue", "google.protobuf.Int32Value", "google.protobuf.Int64Value",
		"google.protobuf.UInt32Value", "google.protobuf.UInt64Value", "google.protobuf.FloatValue",
		"google.protobuf.DoubleValue", "google.protobuf.StringValue", "google.protobuf.BytesValue",
		"google.protobuf.Struct", "google.protobuf.ListValue", "google.protobuf.Value",
		"google.protobuf.FieldMask", "google.protobuf.Empty":
		return cloneMergeValue(src).(map[string]any), nil
	}
	dst := make(map[string]any, len(src))
	seen := make(map[string]string, len(src))
	for key, value := range src {
		field := desc.Fields().ByName(protoreflect.Name(key))
		if field == nil {
			field = desc.Fields().ByJSONName(key)
		}
		if field == nil {
			dst[key] = cloneMergeValue(value)
			continue
		}
		canonical := field.JSONName()
		if previous, ok := seen[canonical]; ok {
			return nil, fmt.Errorf("config: duplicate field %q in source (%q and %q)", canonical, previous, key)
		}
		seen[canonical] = key
		canonicalValue, err := canonicalizeProtoValue(value, field)
		if err != nil {
			return nil, err
		}
		dst[canonical] = canonicalValue
	}
	return dst, nil
}

func canonicalizeProtoValue(value any, field protoreflect.FieldDescriptor) (any, error) {
	if field.IsList() {
		items, ok := value.([]any)
		if !ok || field.Message() == nil {
			return cloneMergeValue(value), nil
		}
		result := make([]any, len(items))
		for i, item := range items {
			if object, ok := item.(map[string]any); ok {
				canonical, err := canonicalizeProtoMap(object, field.Message())
				if err != nil {
					return nil, err
				}
				result[i] = canonical
			} else {
				result[i] = cloneMergeValue(item)
			}
		}
		return result, nil
	}
	if field.IsMap() {
		values, ok := value.(map[string]any)
		if !ok || field.MapValue().Message() == nil {
			return cloneMergeValue(value), nil
		}
		result := make(map[string]any, len(values))
		for key, item := range values {
			if object, ok := item.(map[string]any); ok {
				canonical, err := canonicalizeProtoMap(object, field.MapValue().Message())
				if err != nil {
					return nil, err
				}
				result[key] = canonical
			} else {
				result[key] = cloneMergeValue(item)
			}
		}
		return result, nil
	}
	if field.Message() != nil {
		if object, ok := value.(map[string]any); ok {
			return canonicalizeProtoMap(object, field.Message())
		}
	}
	return cloneMergeValue(value), nil
}
