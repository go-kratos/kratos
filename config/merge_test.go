package config

import (
	"reflect"
	"testing"
)

func TestDefaultMerge(t *testing.T) {
	tests := []struct {
		name string
		dst  map[string]any
		src  any
		want map[string]any
	}{
		{
			name: "merge nested maps and override leaves",
			dst: map[string]any{
				"server": map[string]any{
					"http": map[string]any{
						"addr": "0.0.0.0",
						"port": 80,
					},
					"grpc": true,
				},
				"endpoints": []any{"a.example.com"},
			},
			src: map[string]any{
				"server": map[string]any{
					"http": map[string]any{
						"port": 8080,
						"tls":  true,
					},
				},
				"endpoints": []any{"b.example.com", "c.example.com"},
			},
			want: map[string]any{
				"server": map[string]any{
					"http": map[string]any{
						"addr": "0.0.0.0",
						"port": 8080,
						"tls":  true,
					},
					"grpc": true,
				},
				"endpoints": []any{"b.example.com", "c.example.com"},
			},
		},
		{
			name: "override type conflicts and nil values",
			dst: map[string]any{
				"map_to_scalar": map[string]any{"value": "old"},
				"scalar_to_map": "old",
				"nil_value":     "old",
			},
			src: map[string]any{
				"map_to_scalar": "new",
				"scalar_to_map": map[string]any{"value": "new"},
				"nil_value":     nil,
			},
			want: map[string]any{
				"map_to_scalar": "new",
				"scalar_to_map": map[string]any{"value": "new"},
				"nil_value":     nil,
			},
		},
		{
			name: "convert map keys",
			dst:  map[string]any{},
			src: map[any]any{
				"service": map[any]any{
					"name": "kratos",
				},
			},
			want: map[string]any{
				"service": map[string]any{
					"name": "kratos",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := defaultMerge(&tt.dst, tt.src); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.dst, tt.want) {
				t.Fatalf("defaultMerge() = %#v, want %#v", tt.dst, tt.want)
			}
		})
	}
}

func TestDefaultMergeClonesSourceValues(t *testing.T) {
	srcMap := map[string]any{"name": "kratos"}
	srcSlice := []any{map[string]any{"port": 8000}}
	dst := map[string]any{}

	if err := defaultMerge(&dst, map[string]any{
		"server":    srcMap,
		"listeners": srcSlice,
	}); err != nil {
		t.Fatal(err)
	}

	srcMap["name"] = "changed"
	srcSlice[0].(map[string]any)["port"] = 9000

	want := map[string]any{
		"server":    map[string]any{"name": "kratos"},
		"listeners": []any{map[string]any{"port": 8000}},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Fatalf("defaultMerge() retained source aliases: got %#v, want %#v", dst, want)
	}
}

func TestDefaultMergeInvalidInput(t *testing.T) {
	dst := map[string]any{}
	if err := defaultMerge(dst, map[string]any{}); err == nil {
		t.Fatal("defaultMerge() error is nil for non-pointer dst")
	}
	if err := defaultMerge(&dst, []any{}); err == nil {
		t.Fatal("defaultMerge() error is nil for non-map src")
	}
}
