package config

import (
	"fmt"
	"testing"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/stretchr/testify/assert"
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
	assert.Error(t, err)

	err = r.Merge(&KeyValue{
		Key:    "b",
		Value:  []byte(`{"nice": "boat", "x": 1}`),
		Format: "json",
	})
	assert.NoError(t, err)
	vv, ok := r.Value("nice")
	assert.True(t, ok)
	vvv, err := vv.String()
	assert.NoError(t, err)
	assert.Equal(t, "boat", vvv)

	err = r.Merge(&KeyValue{
		Key:    "b",
		Value:  []byte(`{"x": 2}`),
		Format: "json",
	})
	assert.NoError(t, err)
	vv, ok = r.Value("x")
	assert.True(t, ok)
	vvx, err := vv.Int()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), vvx)
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
			assert.NoError(t, err)
			vv, ok := r.Value("a.b.X")
			assert.True(t, ok)
			vvv, err := vv.Int()
			assert.NoError(t, err)
			assert.Equal(t, int64(1), vvv)

			assert.NoError(t, err)
			vv, ok = r.Value("a.b.Y")
			assert.True(t, ok)
			vvy, err := vv.String()
			assert.NoError(t, err)
			assert.Equal(t, "lol", vvy)

			assert.NoError(t, err)
			vv, ok = r.Value("a.b.z")
			assert.True(t, ok)
			vvz, err := vv.Bool()
			assert.NoError(t, err)
			assert.Equal(t, true, vvz)

			vv, ok = r.Value("aasasdg=234l.asdfk,")
			assert.False(t, ok)

			vv, ok = r.Value("aas......asdg=234l.asdfk,")
			assert.False(t, ok)

			vv, ok = r.Value("a.b.Y.")
			assert.False(t, ok)
		})
	}
}

func TestReader_Source(t *testing.T) {
	var (
		err error
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
		Key:    "b",
		Value:  []byte(`{"a": {"b": {"X": 1}}}`),
		Format: "json",
	})
	b, err := r.Source()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"a":{"b":{"X":1}}}`), b)
}
