package yaml

import (
	"math"
	"reflect"
	"testing"
)

func TestCodec_Unmarshal(t *testing.T) {
	tests := []struct {
		value interface{}
		data  string
	}{
		{
			data:  "",
			value: (*struct{})(nil),
		},
		{
			data: "{}", value: &struct{}{},
		},
		{
			data:  "v: hi",
			value: map[string]string{"v": "hi"},
		},
		{
			data: "v: hi", value: map[string]interface{}{"v": "hi"},
		},
		{
			data:  "v: true",
			value: map[string]string{"v": "true"},
		},
		{
			data:  "v: true",
			value: map[string]interface{}{"v": true},
		},
		{
			data:  "v: 10",
			value: map[string]interface{}{"v": 10},
		},
		{
			data:  "v: 0b10",
			value: map[string]interface{}{"v": 2},
		},
		{
			data:  "v: 0xA",
			value: map[string]interface{}{"v": 10},
		},
		{
			data:  "v: 4294967296",
			value: map[string]int64{"v": 4294967296},
		},
		{
			data:  "v: 0.1",
			value: map[string]interface{}{"v": 0.1},
		},
		{
			data:  "v: .1",
			value: map[string]interface{}{"v": 0.1},
		},
		{
			data:  "v: .Inf",
			value: map[string]interface{}{"v": math.Inf(+1)},
		},
		{
			data:  "v: -.Inf",
			value: map[string]interface{}{"v": math.Inf(-1)},
		},
		{
			data:  "v: -10",
			value: map[string]interface{}{"v": -10},
		},
		{
			data:  "v: -.1",
			value: map[string]interface{}{"v": -0.1},
		},
	}
	for _, tt := range tests {
		v := reflect.ValueOf(tt.value).Type()
		value := reflect.New(v)
		err := (codec{}).Unmarshal([]byte(tt.data), value.Interface())
		if err != nil {
			t.Fatalf("(codec{}).Unmarshal should not return err")
		}
	}
	spec := struct {
		B map[string]interface{}
		A string
	}{A: "a"}
	err := (codec{}).Unmarshal([]byte("v: hi"), &spec.B)
	if err != nil {
		t.Fatalf("(codec{}).Unmarshal should not return err")
	}
}

func TestCodec_Marshal(t *testing.T) {
	value := map[string]string{"v": "hi"}
	got, err := (codec{}).Marshal(value)
	if err != nil {
		t.Fatalf("should not return err")
	}
	if string(got) != "v: hi\n" {
		t.Fatalf("want \"v: hi\n\" return \"%s\"", string(got))
	}
}
