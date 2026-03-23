package apollo

import (
	"testing"

	"github.com/apolloconfig/agollo/v4/storage"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"
)

func Test_onChange(t *testing.T) {
	s := map[string]struct {
		Name string `yaml:"name"`
	}{
		"app": {
			Name: "new",
		},
	}
	codec := encoding.GetCodec(yaml)
	val, _ := codec.Marshal(s)
	c := customChangeListener{}
	tests := []struct {
		name      string
		namespace string
		changes   map[string]*storage.ConfigChange
		kvs       []*config.KeyValue
	}{
		{
			"test yaml onChange",
			"app.yaml",
			map[string]*storage.ConfigChange{
				"name": {
					OldValue:   "old",
					NewValue:   "new",
					ChangeType: storage.MODIFIED,
				},
			},
			[]*config.KeyValue{
				{
					Key:    "app.yaml",
					Value:  val,
					Format: yaml,
				},
			},
		},
		{
			"test json onChange",
			"app.json",
			map[string]*storage.ConfigChange{
				"content": {
					OldValue:   `{"name":"old"}`,
					NewValue:   `{"name":"new"}`,
					ChangeType: storage.MODIFIED,
				},
			},
			[]*config.KeyValue{
				{
					Key:    "app.json",
					Value:  []byte(`{"name":"new"}`),
					Format: json,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kvs := c.onChange(tt.namespace, tt.changes)
			if len(kvs) != len(tt.kvs) {
				t.Errorf("len(kvs) = %v, want %v", len(kvs), len(tt.kvs))
			}
			for i := range kvs {
				if kvs[i].Format != tt.kvs[i].Format || kvs[i].Key != tt.kvs[i].Key || string(kvs[i].Value) != string(tt.kvs[i].Value) {
					t.Errorf("got %v, want %v", kvs[i], tt.kvs[i])
				}
			}
		})
	}
}

func Test_onChange_deletedContent(t *testing.T) {
	c := customChangeListener{}

	t.Run("json content deleted should not panic", func(t *testing.T) {
		changes := map[string]*storage.ConfigChange{
			"content": {
				OldValue:   `{"name":"old"}`,
				NewValue:   nil,
				ChangeType: storage.DELETED,
			},
		}
		kvs := c.onChange("app.json", changes)
		// NewValue is nil, so the original config path is skipped;
		// falls through to resolve path which also skips nil NewValue.
		if len(kvs) != 1 {
			t.Fatalf("expected 1 kv, got %d", len(kvs))
		}
	})

	t.Run("yaml content deleted should not panic", func(t *testing.T) {
		changes := map[string]*storage.ConfigChange{
			"content": {
				OldValue:   "name: old",
				NewValue:   nil,
				ChangeType: storage.DELETED,
			},
		}
		kvs := c.onChange("app.yaml", changes)
		if len(kvs) != 1 {
			t.Fatalf("expected 1 kv, got %d", len(kvs))
		}
	})

	t.Run("properties key deleted should not panic", func(t *testing.T) {
		changes := map[string]*storage.ConfigChange{
			"name": {
				OldValue:   "old",
				NewValue:   nil,
				ChangeType: storage.DELETED,
			},
		}
		kvs := c.onChange("app", changes)
		if len(kvs) != 1 {
			t.Fatalf("expected 1 kv, got %d", len(kvs))
		}
	})
}

func Test_onChange_nonStringNewValue(t *testing.T) {
	c := customChangeListener{}

	t.Run("json content with non-string NewValue should not panic", func(t *testing.T) {
		changes := map[string]*storage.ConfigChange{
			"content": {
				OldValue:   `{"name":"old"}`,
				NewValue:   12345,
				ChangeType: storage.MODIFIED,
			},
		}
		// Should not panic; falls through to resolve path
		kvs := c.onChange("app.json", changes)
		if kvs == nil {
			t.Fatal("expected non-nil kvs")
		}
	})
}
