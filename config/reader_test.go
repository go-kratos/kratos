package config

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3/encoding"
)

func TestReaderMergeCallbacksCanReadConfig(t *testing.T) {
	for _, kind := range []string{"decoder", "merge"} {
		t.Run(kind, func(t *testing.T) {
			var c *config
			readConfig := false
			read := func() error {
				if !readConfig {
					return nil
				}
				value, err := c.Value("existing").String()
				if err != nil {
					return err
				}
				if value != "readable" {
					return fmt.Errorf("existing = %q, want readable", value)
				}
				return nil
			}
			var option Option
			if kind == "decoder" {
				option = WithDecoder(func(kv *KeyValue, dst map[string]any) error {
					if err := read(); err != nil {
						return err
					}
					return defaultDecoder(kv, dst)
				})
			} else {
				option = WithMergeFunc(func(dst, src any) error {
					if err := read(); err != nil {
						return err
					}
					return defaultMerge(dst, src)
				})
			}
			c = New(option).(*config)
			if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(`{"existing":"readable"}`)}); err != nil {
				t.Fatal(err)
			}
			// Keep existing out of the Config cache so Value must read the reader.
			readConfig = true
			done := make(chan error, 1)
			go func() {
				done <- c.reader.Merge(&KeyValue{Format: "json", Value: []byte(`{"next":"updated"}`)})
			}()
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(time.Second):
				t.Fatal("merge callback blocked while reading config")
			}
			if value, err := c.Value("next").String(); err != nil || value != "updated" {
				t.Fatalf("next = %q, %v; want updated", value, err)
			}
		})
	}
}

func TestReader_Merge(t *testing.T) {
	var (
		err error
		ok  bool
	)
	opts := options{
		decoder: func(kv *KeyValue, v map[string]any) error {
			if codec := encoding.GetCodec(kv.Format); codec != nil {
				return codec.Unmarshal(kv.Value, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", kv.Key, kv.Format)
		},
		resolver: defaultResolver,
		merge:    defaultMerge,
	}
	r := newReader(opts)
	err = r.Merge(&KeyValue{
		Key:    "a",
		Value:  []byte("bad"),
		Format: "json",
	})
	if err == nil {
		t.Fatal("err is nil")
	}

	err = r.Merge(&KeyValue{
		Key:    "b",
		Value:  []byte(`{"nice": "boat", "x": 1}`),
		Format: "json",
	})
	if err != nil {
		t.Fatal(err)
	}
	vv, ok := r.Value("nice")
	if !ok {
		t.Fatal("ok is false")
	}
	vvv, err := vv.String()
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}
	vv, ok = r.Value("x")
	if !ok {
		t.Fatal("ok is false")
	}
	vvx, err := vv.Int()
	if err != nil {
		t.Fatal(err)
	}
	if vvx != 2 {
		t.Fatal("vvx is not equal to 2")
	}
}

func TestReaderMergeCompactsLayerHistory(t *testing.T) {
	r := newReader(options{
		decoder:  defaultDecoder,
		resolver: defaultResolver,
		merge:    defaultMerge,
	}).(*reader)
	for i := 0; i < 100; i++ {
		if err := r.Merge(&KeyValue{
			Key:    "config",
			Value:  []byte(fmt.Sprintf(`{"value":%d}`, i)),
			Format: "json",
		}); err != nil {
			t.Fatal(err)
		}
	}
	if got := len(r.layers); got != 1 {
		t.Fatalf("expected one retained layer for repeated updates, got %d", got)
	}
}

func BenchmarkCompactLayersDistinctKeys(b *testing.B) {
	for _, count := range []int{500, 1000, 2000, 10000} {
		b.Run(fmt.Sprintf("keys_%d", count), func(b *testing.B) {
			layers := make([]map[string]any, count)
			for i := range layers {
				layers[i] = map[string]any{fmt.Sprintf("key_%d", i): "value"}
			}
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if got := compactLayers(layers); len(got) != count {
					b.Fatalf("retained %d layers, want %d", len(got), count)
				}
			}
		})
	}
}

func TestReader_Value(t *testing.T) {
	opts := options{
		decoder: func(kv *KeyValue, v map[string]any) error {
			if codec := encoding.GetCodec(kv.Format); codec != nil {
				return codec.Unmarshal(kv.Value, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", kv.Key, kv.Format)
		},
		resolver: defaultResolver,
		merge:    defaultMerge,
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
				t.Fatal(err)
			}
			vv, ok := r.Value("a.b.X")
			if !ok {
				t.Fatal("ok is false")
			}
			vvv, err := vv.Int()
			if err != nil {
				t.Fatal(err)
			}
			if int64(1) != vvv {
				t.Fatal("vvv is not equal to 1")
			}

			vv, ok = r.Value("a.b.Y")
			if !ok {
				t.Fatal("ok is false")
			}
			vvy, err := vv.String()
			if err != nil {
				t.Fatal(err)
			}
			if vvy != "lol" {
				t.Fatal(`vvy is not equal to "lol"`)
			}

			vv, ok = r.Value("a.b.z")
			if !ok {
				t.Fatal("ok is false")
			}
			vvz, err := vv.Bool()
			if err != nil {
				t.Fatal(err)
			}
			if !vvz {
				t.Fatal("vvz is not equal to true")
			}

			_, ok = r.Value("aasasdg=234l.asdfk,")
			if ok {
				t.Fatal("ok is true")
			}

			_, ok = r.Value("aas......asdg=234l.asdfk,")
			if ok {
				t.Fatal("ok is true")
			}

			_, ok = r.Value("a.b.Y.")
			if ok {
				t.Fatal("ok is true")
			}
		})
	}
}

func TestReader_Source(t *testing.T) {
	var err error
	opts := options{
		decoder: func(kv *KeyValue, v map[string]any) error {
			if codec := encoding.GetCodec(kv.Format); codec != nil {
				return codec.Unmarshal(kv.Value, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", kv.Key, kv.Format)
		},
		resolver: defaultResolver,
		merge:    defaultMerge,
	}
	r := newReader(opts)
	err = r.Merge(&KeyValue{
		Key:    "b",
		Value:  []byte(`{"a": {"b": {"X": 1}}}`),
		Format: "json",
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Source()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]byte(`{"a":{"b":{"X":1}}}`), b) {
		t.Fatal("[]byte(`{\"a\":{\"b\":{\"X\":1}}}`) is not equal to b")
	}
}

func TestCloneMap(t *testing.T) {
	tests := []struct {
		input map[string]any
		want  map[string]any
	}{
		{
			input: map[string]any{
				"a": 1,
				"b": "2",
				"c": true,
			},
			want: map[string]any{
				"a": 1,
				"b": "2",
				"c": true,
			},
		},
		{
			input: map[string]any{},
			want:  map[string]any{},
		},
		{
			input: nil,
			want:  map[string]any{},
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
		input any
		want  any
	}{
		{
			input: map[string]any{
				"a": 1,
				"b": "2",
				"c": true,
				"d": []byte{65, 66, 67},
			},
			want: map[string]any{
				"a": 1,
				"b": "2",
				"c": true,
				"d": "ABC",
			},
		},
		{
			input: []any{1, 2.0, "3", true, nil, []any{1, 2.0, "3", true, nil}},
			want:  []any{1, 2.0, "3", true, nil, []any{1, 2.0, "3", true, nil}},
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
	m := map[string]any{
		"a": 1,
		"b": map[string]any{
			"c": "3",
			"d": map[string]any{
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
