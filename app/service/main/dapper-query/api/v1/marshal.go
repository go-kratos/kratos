package v1

import (
	"encoding/json"

	"go-common/app/tool/protoc-gen-bm/jsonpb"
)

// MarshalJSON .
func (t *TagValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

// MarshalJSON .
func (t *TagValue_StringValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.StringValue)
}

// MarshalJSON .
func (t *TagValue_Int64Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Int64Value)
}

// MarshalJSON .
func (t *TagValue_BoolValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.BoolValue)
}

// MarshalJSON .
func (t *TagValue_FloatValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.FloatValue)
}

// MarshalJSONPB .
func (t *TagValue) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(t.Value)
}

// MarshalJSONPB .
func (t *TagValue_StringValue) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(t.StringValue)
}

// MarshalJSONPB .
func (t *TagValue_Int64Value) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(t.Int64Value)
}

// MarshalJSONPB .
func (t *TagValue_BoolValue) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(t.BoolValue)
}

// MarshalJSONPB .
func (t *TagValue_FloatValue) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(t.FloatValue)
}
