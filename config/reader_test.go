package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/encoding"
)

func TestReader_Merge(t *testing.T) {
	var (
		err error
		ok  bool
	)
	opts := options{
		decoder: func(kv *KeyValue, v map[string]interface{}) error {
			if codec := encoding.GetCodec(kv.Format); codec != nil {
				return codec.Unmarshal(kv.Value, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", kv.Key, kv.Format)
		},
		resolver: defaultResolver,
	}
	r := newReader(opts)
	err = r.Merge(&KeyValue{
		Key:    "a",
		Value:  []byte("bad"),
		Format: "json",
	})
	if err == nil {
		t.Fatal(`err is nil`)
	}

	err = r.Merge(&KeyValue{
		Key:    "b",
		Value:  []byte(`{"nice": "boat", "x": 1}`),
		Format: "json",
	})
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	vv, ok := r.Value("nice")
	if !ok {
		t.Fatal(`ok is false`)
	}
	vvv, err := vv.String()
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	if vvv != "boat" {
		t.Fatal(`vvv is not equal to "boat"`)
	}

	err = r.Merge(&KeyValue{
		Key:    "b",
		Value:  []byte(`{"x": 2}`),
		Format: "json",
	})
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	vv, ok = r.Value("x")
	if !ok {
		t.Fatal(`ok is false`)
	}
	vvx, err := vv.Int()
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	if int64(2) != vvx {
		t.Fatal(`vvx is not equal to 2`)
	}
}

func TestReader_Value(t *testing.T) {
	opts := options{
		decoder: func(kv *KeyValue, v map[string]interface{}) error {
			if codec := encoding.GetCodec(kv.Format); codec != nil {
				return codec.Unmarshal(kv.Value, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", kv.Key, kv.Format)
		},
		resolver: defaultResolver,
	}

	ymlval := `
a: 
  b: 
    X: 1
    Y: "lol"
    z: true
`
	tests := []struct {
		name string
		kv   KeyValue
	}{
		{
			name: "json value",
			kv: KeyValue{
				Key:    "config",
				Value:  []byte(`{"a": {"b": {"X": 1, "Y": "lol", "z": true}}}`),
				Format: "json",
			},
		},
		{
			name: "yaml value",
			kv: KeyValue{
				Key:    "config",
				Value:  []byte(ymlval),
				Format: "yaml",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := newReader(opts)
			err := r.Merge(&test.kv)
			if err != nil {
				t.Fatal(`err is not nil`)
			}
			vv, ok := r.Value("a.b.X")
			if !ok {
				t.Fatal(`ok is false`)
			}
			vvv, err := vv.Int()
			if err != nil {
				t.Fatal(`err is not nil`)
			}
			if int64(1) != vvv {
				t.Fatal(`vvv is not equal to 1`)
			}

			vv, ok = r.Value("a.b.Y")
			if !ok {
				t.Fatal(`ok is false`)
			}
			vvy, err := vv.String()
			if err != nil {
				t.Fatal(`err is not nil`)
			}
			if vvy != "lol" {
				t.Fatal(`vvy is not equal to "lol"`)
			}

			vv, ok = r.Value("a.b.z")
			if !ok {
				t.Fatal(`ok is false`)
			}
			vvz, err := vv.Bool()
			if err != nil {
				t.Fatal(`err is not nil`)
			}
			if !vvz {
				t.Fatal(`vvz is not equal to true`)
			}

			_, ok = r.Value("aasasdg=234l.asdfk,")
			if ok {
				t.Fatal(`ok is true`)
			}

			_, ok = r.Value("aas......asdg=234l.asdfk,")
			if ok {
				t.Fatal(`ok is true`)
			}

			_, ok = r.Value("a.b.Y.")
			if ok {
				t.Fatal(`ok is true`)
			}
		})
	}
}

func TestReader_Source(t *testing.T) {
	var err error
	opts := options{
		decoder: func(kv *KeyValue, v map[string]interface{}) error {
			if codec := encoding.GetCodec(kv.Format); codec != nil {
				return codec.Unmarshal(kv.Value, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", kv.Key, kv.Format)
		},
		resolver: defaultResolver,
	}
	r := newReader(opts)
	err = r.Merge(&KeyValue{
		Key:    "b",
		Value:  []byte(`{"a": {"b": {"X": 1}}}`),
		Format: "json",
	})
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	b, err := r.Source()
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	if !reflect.DeepEqual([]byte(`{"a":{"b":{"X":1}}}`), b) {
		t.Fatal("[]byte(`{\"a\":{\"b\":{\"X\":1}}}`) is not equal to b")
	}
}

func TestCloneMap(t *testing.T) {
	tests := []struct {
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{
			input: map[string]interface{}{
				"a": 1,
				"b": "2",
				"c": true,
			},
			want: map[string]interface{}{
				"a": 1,
				"b": "2",
				"c": true,
			},
		},
		{
			input: map[string]interface{}{},
			want:  map[string]interface{}{},
		},
		{
			input: nil,
			want:  map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		if got, err := cloneMap(tt.input); err != nil {
			t.Errorf("expect no err, got %v", err)
		} else if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("cloneMap(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestConvertMap(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{
			input: map[string]interface{}{
				"a": 1,
				"b": "2",
				"c": true,
				"d": []byte{65, 66, 67},
			},
			want: map[string]interface{}{
				"a": 1,
				"b": "2",
				"c": true,
				"d": "ABC",
			},
		},
		{
			input: []interface{}{1, 2.0, "3", true, nil, []interface{}{1, 2.0, "3", true, nil}},
			want:  []interface{}{1, 2.0, "3", true, nil, []interface{}{1, 2.0, "3", true, nil}},
		},
		{
			input: []byte{65, 66, 67},
			want:  "ABC",
		},
	}
	for _, tt := range tests {
		if got := convertMap(tt.input); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("convertMap(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestReadValue(t *testing.T) {
	m := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"c": "3",
			"d": map[string]interface{}{
				"e": true,
			},
		},
	}
	va := atomicValue{}
	va.Store(1)

	vbc := atomicValue{}
	vbc.Store("3")

	vbde := atomicValue{}
	vbde.Store(true)

	tests := []struct {
		path string
		want atomicValue
	}{
		{
			path: "a",
			want: va,
		},
		{
			path: "b.c",
			want: vbc,
		},
		{
			path: "b.d.e",
			want: vbde,
		},
	}
	for _, tt := range tests {
		if got, found := readValue(m, tt.path); !found {
			t.Errorf("expect found %v in %v, but not.", tt.path, m)
		} else if got.Load() != tt.want.Load() {
			t.Errorf("readValue(%v, %v) = %v, want %v", m, tt.path, got, tt.want)
		}
	}
}
