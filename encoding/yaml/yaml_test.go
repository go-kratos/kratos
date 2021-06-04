package yaml

import (
	"math"
	"reflect"
	"testing"
)

func TestCodec_Unmarshal(t *testing.T) {

	tests := []struct {
		data  string
		value interface{}
	}{
		{
			"",
			(*struct{})(nil),
		},
		{
			"{}", &struct{}{},
		}, {
			"v: hi",
			map[string]string{"v": "hi"},
		}, {
			"v: hi", map[string]interface{}{"v": "hi"},
		}, {
			"v: true",
			map[string]string{"v": "true"},
		}, {
			"v: true",
			map[string]interface{}{"v": true},
		}, {
			"v: 10",
			map[string]interface{}{"v": 10},
		}, {
			"v: 0b10",
			map[string]interface{}{"v": 2},
		}, {
			"v: 0xA",
			map[string]interface{}{"v": 10},
		}, {
			"v: 4294967296",
			map[string]int64{"v": 4294967296},
		}, {
			"v: 0.1",
			map[string]interface{}{"v": 0.1},
		}, {
			"v: .1",
			map[string]interface{}{"v": 0.1},
		}, {
			"v: .Inf",
			map[string]interface{}{"v": math.Inf(+1)},
		}, {
			"v: -.Inf",
			map[string]interface{}{"v": math.Inf(-1)},
		}, {
			"v: -10",
			map[string]interface{}{"v": -10},
		}, {
			"v: -.1",
			map[string]interface{}{"v": -0.1},
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
		A string
		B map[string]interface{}
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
