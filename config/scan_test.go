package config

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/typepb"
)

func TestScanKeepsResolverDefaultsBeforeFirstLayer(t *testing.T) {
	c := New(WithResolver(func(values map[string]any) error {
		values["apply_schema"] = true
		return nil
	})).(*config)
	// An initially empty source still runs the resolver before its first update.
	if err := c.reader.Merge(); err != nil {
		t.Fatal(err)
	}
	if err := c.reader.Resolve(); err != nil {
		t.Fatal(err)
	}
	desc := testAliasMessageDescriptor(t)
	for _, body := range []string{"", `{"unrelated":"first update"}`} {
		if body != "" {
			if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(body)}); err != nil {
				t.Fatal(err)
			}
			if err := c.reader.Resolve(); err != nil {
				t.Fatal(err)
			}
		}
		value, err := c.Value("apply_schema").Bool()
		if err != nil || !value {
			t.Fatalf("Value = %v, %v; want true", value, err)
		}
		var plain struct {
			ApplySchema bool `json:"apply_schema"`
		}
		if err := c.Scan(&plain); err != nil {
			t.Fatal(err)
		}
		msg := dynamicpb.NewMessage(desc)
		if err := c.Scan(msg); err != nil {
			t.Fatal(err)
		}
		if !plain.ApplySchema || !msg.Get(desc.Fields().ByName("apply_schema")).Bool() {
			t.Fatalf("update %q lost resolver default: Go Scan = %v, protobuf Scan = %v", body, plain, msg)
		}
	}
}

func TestScanPreservesNullResetsAcrossAliases(t *testing.T) {
	for _, bodies := range [][]string{
		{`{"nestedConfig":{"db_host":"old"}}`, `{"nested_config":null}`, `{"nested_config":{"url":"new"}}`},
		{`{"nested_config":{"db_host":"old"}}`, `{"nestedConfig":null}`, `{"nestedConfig":{"url":"new"}}`},
		{`{"nested_config":{"db_host":"old"}}`, `{"nested_config":null}`, `{"nested_config":{"url":"new"}}`},
	} {
		c := New().(*config)
		for _, body := range bodies {
			if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(body)}); err != nil {
				t.Fatal(err)
			}
		}
		msg := dynamicpb.NewMessage(testResolutionDescriptor(t))
		if err := c.Scan(msg); err != nil {
			t.Fatal(err)
		}
		nested := msg.Get(msg.Descriptor().Fields().ByName("nested_config")).Message()
		old := nested.Get(nested.Descriptor().Fields().ByName("db_host")).String()
		latest := nested.Get(nested.Descriptor().Fields().ByName("url")).String()
		if old != "" || latest != "new" {
			t.Fatalf("layers %v: db_host = %q, url = %q; want only new url", bodies, old, latest)
		}
	}
}

func TestScanResolverDeletesOriginalSpelling(t *testing.T) {
	for _, test := range []struct {
		name   string
		layers []string
		path   []string
		want   string
	}{
		{"JSON name", []string{`{"db_host":"old"}`, `{"dbHost":"new"}`}, []string{"dbHost"}, `{"dbHost":"old"}`},
		{"proto name", []string{`{"dbHost":"old"}`, `{"db_host":"new"}`}, []string{"db_host"}, `{"dbHost":"old"}`},
		{"single layer", []string{`{"dbHost":"new"}`}, []string{"dbHost"}, `{}`},
		{"same spelling", []string{`{"db_host":"old"}`, `{"db_host":"new"}`}, []string{"db_host"}, `{}`},
		{
			"nested field",
			[]string{`{"nested_config":{"db_host":"old"}}`, `{"nested_config":{"dbHost":"new"}}`},
			[]string{"nested_config", "dbHost"},
			`{"nestedConfig":{"dbHost":"old"}}`,
		},
		{
			"whole object",
			[]string{`{"nested_config":{"db_host":"old"}}`, `{"nestedConfig":{"url":"new"}}`},
			[]string{"nestedConfig"},
			`{"nestedConfig":{"dbHost":"old"}}`,
		},
		{
			"reset and object",
			[]string{`{"nestedConfig":{"db_host":"old"}}`, `{"nested_config":null}`, `{"nested_config":{"url":"new"}}`},
			[]string{"nested_config"},
			`{"nestedConfig":{"dbHost":"old"}}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := New(WithResolver(func(values map[string]any) error {
				for _, key := range test.path[:len(test.path)-1] {
					values = values[key].(map[string]any)
				}
				delete(values, test.path[len(test.path)-1])
				return nil
			})).(*config)
			for _, body := range test.layers {
				if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(body)}); err != nil {
					t.Fatal(err)
				}
			}
			if err := c.reader.Resolve(); err != nil {
				t.Fatal(err)
			}
			desc := testResolutionDescriptor(t)
			want := dynamicpb.NewMessage(desc)
			if err := unmarshalJSON([]byte(test.want), want); err != nil {
				t.Fatal(err)
			}
			for range 2 {
				got := dynamicpb.NewMessage(desc)
				if err := c.Scan(got); err != nil {
					t.Fatal(err)
				}
				if !proto.Equal(got, want) {
					t.Fatalf("got %v, want %v", got, want)
				}
			}
		})
	}
}

func TestCompactLayersPreservesProtoReplay(t *testing.T) {
	// Exercise interleaved aliases, empty maps, and resets at two nesting levels.
	variants := []map[string]any{
		{"nested_config": nil},
		{"nestedConfig": nil},
		{"nested_config": map[string]any{"db_host": "old"}},
		{"nestedConfig": map[string]any{"dbHost": "new"}},
		{"nested_config": map[string]any{"url": "url"}},
		{"nestedConfig": map[string]any{}},
		{"nested_config": map[string]any{"nested_config": nil}},
		{"nestedConfig": map[string]any{"nestedConfig": map[string]any{"db_host": "deep"}}},
	}
	desc := testResolutionDescriptor(t)
	replay := func(layers []map[string]any) map[string]any {
		merged := make(map[string]any)
		for _, layer := range layers {
			canonical, err := canonicalizeProtoMap(layer, desc)
			if err != nil {
				t.Fatal(err)
			}
			mergeMap(merged, canonical)
		}
		return merged
	}
	const depth = 4
	count := 1
	for range depth {
		count *= len(variants)
	}
	for sequence := range count {
		layers := make([]map[string]any, 0, depth)
		var compacted []map[string]any
		index := sequence
		for range depth {
			layer := variants[index%len(variants)]
			index /= len(variants)
			layers = append(layers, layer)
			compacted = compactLayers(append(compacted, layer))
		}
		want := replay(layers)
		for _, candidate := range [][]map[string]any{compacted, compactLayers(layers)} {
			if got := replay(candidate); !reflect.DeepEqual(got, want) {
				t.Fatalf("layers %v: compacted result %v, want %v", layers, got, want)
			}
		}
	}
}

func TestScanCanonicalizesRegularGoogleProtobufMessages(t *testing.T) {
	for _, bodies := range [][]string{
		{`{"source_context":{"file_name":"old.proto"}}`, `{"sourceContext":{"fileName":"new.proto"}}`},
		{`{"sourceContext":{"fileName":"old.proto"}}`, `{"source_context":{"file_name":"new.proto"}}`},
		{`{"sourceContext":{"file_name":"old.proto"}}`, `{"sourceContext":{"fileName":"new.proto"}}`},
	} {
		c := New().(*config)
		for _, body := range bodies {
			if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(body)}); err != nil {
				t.Fatal(err)
			}
		}
		if err := c.reader.Resolve(); err != nil {
			t.Fatal(err)
		}
		var msg typepb.Type
		if err := c.Scan(&msg); err != nil {
			t.Fatalf("layers %v: %v", bodies, err)
		}
		if got := msg.GetSourceContext().GetFileName(); got != "new.proto" {
			t.Fatalf("layers %v: file name = %q, want new.proto", bodies, got)
		}
	}
}

func TestConfigScanKeepsGoJSONSemanticsAcrossSources(t *testing.T) {
	type bootstrap struct {
		ApplySchema bool `json:"apply_schema"`
	}

	c := New(WithSource(
		newTestJSONSource(`{"apply_schema":true}`),
		newTestJSONSource(`{"applySchema":false}`),
	))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	var got bootstrap
	if err := c.Scan(&got); err != nil {
		t.Fatal(err)
	}
	if !got.ApplySchema {
		t.Fatal("expected encoding/json to ignore the untagged applySchema key")
	}
}

func TestConfigScanHandlesInterfaceDestination(t *testing.T) {
	c := New(WithSource(newTestJSONSource(`{"value":"ok"}`)))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	var got any
	if err := c.Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got.(map[string]any)["value"] != "ok" {
		t.Fatalf("unexpected scan result: %#v", got)
	}
}

func TestConfigScanKeepsGoJSONSemanticsForEmbeddedFields(t *testing.T) {
	type Embedded struct {
		ApplySchema bool `json:"apply_schema"`
	}
	type bootstrap struct {
		Embedded
	}

	c := New(WithSource(
		newTestJSONSource(`{"apply_schema":true}`),
		newTestJSONSource(`{"applySchema":false}`),
	))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	var got bootstrap
	if err := c.Scan(&got); err != nil {
		t.Fatal(err)
	}
	if !got.ApplySchema {
		t.Fatal("expected encoding/json to ignore the untagged embedded applySchema key")
	}
}

func TestConfigScanKeepsGoJSONSemanticsWithinLayer(t *testing.T) {
	type bootstrap struct {
		ApplySchema bool `json:"apply_schema"`
	}
	c := New(WithSource(newTestJSONSource(`{"apply_schema":true,"ApplySchema":false}`)))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}
	var got bootstrap
	if err := c.Scan(&got); err != nil {
		t.Fatal(err)
	}
	if !got.ApplySchema {
		t.Fatal("expected the exact tagged JSON key to be decoded")
	}
}

func TestConfigScanIncludesResolverAddedValues(t *testing.T) {
	type bootstrap struct {
		Value   string `json:"value"`
		Derived string `json:"derived"`
	}
	c := New(
		WithSource(newTestJSONSource(`{"value":"source"}`)),
		WithResolver(func(values map[string]any) error {
			values["derived"] = "resolved"
			return nil
		}),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	var got bootstrap
	if err := c.Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got.Derived != "resolved" {
		t.Fatalf("expected resolver-added value, got %#v", got)
	}
}

func TestConfigScanResolvesProtoFieldAliasesAcrossSources(t *testing.T) {
	desc := testAliasMessageDescriptor(t)
	c := New(WithSource(
		newTestJSONSource(`{"apply_schema":true}`),
		newTestJSONSource(`{"applySchema":false}`),
	))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	msg := dynamicpb.NewMessage(desc)
	if err := c.Scan(msg); err != nil {
		t.Fatal(err)
	}
	field := desc.Fields().ByName("apply_schema")
	if msg.Get(field).Bool() {
		t.Fatal("expected the later applySchema source to override apply_schema")
	}
}

func TestConfigScanResolvesNestedProtoAliases(t *testing.T) {
	desc := testNestedAliasMessageDescriptor(t)
	c := New(WithSource(
		newTestJSONSource(`{"child":{"apply_schema":true},"children":[{"apply_schema":true}]}`),
		newTestJSONSource(`{"child":{"applySchema":false},"children":[{"applySchema":false}]}`),
	))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	msg := dynamicpb.NewMessage(desc)
	if err := c.Scan(msg); err != nil {
		t.Fatal(err)
	}
	child := desc.Fields().ByName("child")
	if msg.Get(child).Message().Get(child.Message().Fields().ByName("apply_schema")).Bool() {
		t.Fatal("expected nested proto alias to use the later value")
	}
	children := desc.Fields().ByName("children")
	if msg.Get(children).List().Get(0).Message().Get(children.Message().Fields().ByName("apply_schema")).Bool() {
		t.Fatal("expected repeated nested proto alias to use the later value")
	}
}

func TestConfigScanReplaysLatestWatchLayer(t *testing.T) {
	desc := testAliasMessageDescriptor(t)
	r := newReader(options{
		decoder:  defaultDecoder,
		resolver: defaultResolver,
		merge:    defaultMerge,
	}).(*reader)
	for _, value := range []string{`{"apply_schema":true}`, `{"applySchema":false}`, `{"applySchema":true}`} {
		if err := r.Merge(&KeyValue{Key: "watch", Value: []byte(value), Format: "json"}); err != nil {
			t.Fatal(err)
		}
	}
	if err := r.Resolve(); err != nil {
		t.Fatal(err)
	}
	msg := dynamicpb.NewMessage(desc)
	data, err := r.SourceFor(msg)
	if err != nil {
		t.Fatal(err)
	}
	if err := unmarshalJSON(data, msg); err != nil {
		t.Fatal(err)
	}
	if !msg.Get(desc.Fields().ByName("apply_schema")).Bool() {
		t.Fatal("expected the latest watch layer to win")
	}
}

func testAliasMessageDescriptor(t *testing.T) protoreflect.MessageDescriptor {
	t.Helper()
	fd, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:    proto.String("scan_alias.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("scan_alias"),
		MessageType: []*descriptorpb.DescriptorProto{{
			Name: proto.String("Bootstrap"),
			Field: []*descriptorpb.FieldDescriptorProto{{
				Name:     proto.String("apply_schema"),
				JsonName: proto.String("applySchema"),
				Number:   proto.Int32(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
			}},
		}},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return fd.Messages().ByName("Bootstrap")
}

func testNestedAliasMessageDescriptor(t *testing.T) protoreflect.MessageDescriptor {
	t.Helper()
	fd, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:    proto.String("nested_scan_alias.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("nested_scan_alias"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("Child"),
				Field: []*descriptorpb.FieldDescriptorProto{{
					Name:     proto.String("apply_schema"),
					JsonName: proto.String("applySchema"),
					Number:   proto.Int32(1),
					Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
				}},
			},
			{
				Name: proto.String("Bootstrap"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("child"),
						JsonName: proto.String("child"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName: proto.String(".nested_scan_alias.Child"),
					},
					{
						Name:     proto.String("children"),
						JsonName: proto.String("children"),
						Number:   proto.Int32(2),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName: proto.String(".nested_scan_alias.Child"),
					},
				},
			},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return fd.Messages().ByName("Bootstrap")
}

func TestScanPreservesPlaceholderPaths(t *testing.T) {
	desc := testResolutionDescriptor(t)
	c := New().(*config)
	for _, body := range []string{
		`{"db_host":"db","url":"${db_host}","nested_config":{"db_host":"nested","url":"${nested_config.db_host:fallback}"}}`,
		`{"items":["${db_host}","${missing:fallback}"]}`,
	} {
		if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(body)}); err != nil {
			t.Fatal(err)
		}
	}
	if err := c.reader.Resolve(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		msg := dynamicpb.NewMessage(desc)
		if err := c.Scan(msg); err != nil {
			t.Fatal(err)
		}
		if got := msg.Get(desc.Fields().ByName("url")).String(); got != "db" {
			t.Errorf("url = %q, want db", got)
		}
		nested := msg.Get(desc.Fields().ByName("nested_config")).Message()
		if got := nested.Get(desc.Fields().ByName("url")).String(); got != "nested" {
			t.Errorf("nested url = %q, want nested", got)
		}
		items := msg.Get(desc.Fields().ByName("items")).List()
		if items.Get(0).String() != "db" || items.Get(1).String() != "fallback" {
			t.Errorf("unexpected items: %v", items)
		}
	}
}

func TestScanDoesNotReplayResolver(t *testing.T) {
	calls := 0
	c := New(WithResolver(func(m map[string]any) error {
		calls++
		m["url"] = "resolved"
		m["db_host"] = "changed"
		delete(m, "items")
		return nil
	})).(*config)
	if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(`{"db_host":"original","items":["remove"]}`)}); err != nil {
		t.Fatal(err)
	}
	if err := c.reader.Resolve(); err != nil {
		t.Fatal(err)
	}
	desc := testResolutionDescriptor(t)
	for i := 0; i < 2; i++ {
		msg := dynamicpb.NewMessage(desc)
		if err := c.Scan(msg); err != nil {
			t.Fatal(err)
		}
		if msg.Get(desc.Fields().ByName("url")).String() != "resolved" ||
			msg.Get(desc.Fields().ByName("db_host")).String() != "changed" ||
			msg.Get(desc.Fields().ByName("items")).List().Len() != 0 {
			t.Fatalf("unexpected resolved message: %v", msg)
		}
		var plain map[string]any
		if err := c.Scan(&plain); err != nil {
			t.Fatal(err)
		}
	}
	if calls != 1 {
		t.Fatalf("resolver called %d times, want 1", calls)
	}
}

func TestScanPreservesCustomMerge(t *testing.T) {
	for _, firstWins := range []bool{false, true} {
		name := "append"
		if firstWins {
			name = "first wins"
		}
		t.Run(name, func(t *testing.T) {
			merge := func(dst, src any) error {
				d := dst.(*map[string]any)
				for k, v := range src.(map[string]any) {
					if old, ok := (*d)[k]; ok {
						if !firstWins {
							(*d)[k] = append(old.([]any), v.([]any)...)
						}
					} else {
						(*d)[k] = cloneMergeValue(v)
					}
				}
				return nil
			}
			c := New(WithMergeFunc(merge)).(*config)
			desc := testResolutionDescriptor(t)
			for i, body := range []string{`{"items":["a"]}`, `{"items":["b"]}`, `{"items":["c"]}`} {
				if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(body)}); err != nil {
					t.Fatal(err)
				}
				if err := c.reader.Resolve(); err != nil {
					t.Fatal(err)
				}
				want := i + 1
				if firstWins {
					want = 1
				}
				msg := dynamicpb.NewMessage(desc)
				if err := c.Scan(msg); err != nil {
					t.Fatal(err)
				}
				items := msg.Get(desc.Fields().ByName("items")).List()
				if items.Len() != want || items.Get(0).String() != "a" {
					t.Fatalf("update %d: proto items = %v, want length %d starting with a", i, items, want)
				}
				var plain struct {
					Items []string `json:"items"`
				}
				if err := c.Scan(&plain); err != nil {
					t.Fatal(err)
				}
				if len(plain.Items) != want || plain.Items[0] != "a" {
					t.Fatalf("update %d: Go items = %v", i, plain.Items)
				}
			}
		})
	}
}

func testResolutionDescriptor(t *testing.T) protoreflect.MessageDescriptor {
	t.Helper()
	fd, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name: proto.String("resolution.proto"), Syntax: proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("Settings"), Field: []*descriptorpb.FieldDescriptorProto{
			{Name: proto.String("db_host"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()},
			{Name: proto.String("url"), Number: proto.Int32(2), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()},
			{
				Name: proto.String("items"), Number: proto.Int32(3),
				Label: descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
			{
				Name: proto.String("nested_config"), Number: proto.Int32(4),
				Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".Settings"),
			},
		}}},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return fd.Messages().Get(0)
}

func TestScanResolvesPlaceholdersAcrossAliases(t *testing.T) {
	for _, bodies := range [][]string{
		{`{"db_host":"${base:old}"}`, `{"dbHost":"${overlay:new}"}`},
		{`{"db_host":"${base:old}"}`, `{"dbHost":"new"}`},
		{`{"dbHost":"${base:old}"}`, `{"db_host":"new"}`},
	} {
		c := New().(*config)
		for _, body := range bodies {
			if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(body)}); err != nil {
				t.Fatal(err)
			}
		}
		if err := c.reader.Resolve(); err != nil {
			t.Fatal(err)
		}
		msg := dynamicpb.NewMessage(testResolutionDescriptor(t))
		if err := c.Scan(msg); err != nil {
			t.Fatal(err)
		}
		if got := msg.Get(msg.Descriptor().Fields().ByName("db_host")).String(); got != "new" {
			t.Fatalf("layers %v: db_host = %q, want new", bodies, got)
		}
	}
}

func TestScanCustomMergeAfterResolution(t *testing.T) {
	calls := 0
	c := New(WithMergeFunc(func(dst, src any) error {
		d := dst.(*map[string]any)
		for k, v := range src.(map[string]any) {
			if old, ok := (*d)[k].([]any); ok {
				(*d)[k] = append(old, v.([]any)...)
			} else {
				(*d)[k] = cloneMergeValue(v)
			}
		}
		return nil
	}), WithResolver(func(m map[string]any) error {
		calls++
		m["items"] = []any{"resolved"}
		return nil
	})).(*config)
	desc := testResolutionDescriptor(t)
	for i := 0; i < 2; i++ {
		if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(`{"items":["raw"]}`)}); err != nil {
			t.Fatal(err)
		}
		if err := c.reader.Resolve(); err != nil {
			t.Fatal(err)
		}
		msg := dynamicpb.NewMessage(desc)
		if err := c.Scan(msg); err != nil {
			t.Fatal(err)
		}
		items := msg.Get(desc.Fields().ByName("items")).List()
		if items.Len() != 1 || items.Get(0).String() != "resolved" {
			t.Fatalf("unexpected resolved items: %v", items)
		}
	}
	if calls != 2 {
		t.Fatalf("resolver called %d times, want 2", calls)
	}
}

func TestScanCompactsResolvedWatchUpdates(t *testing.T) {
	c := New().(*config)
	desc := testResolutionDescriptor(t)
	for i := 0; i < 40; i++ {
		if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(`{"db_host":"${missing:db}"}`)}); err != nil {
			t.Fatal(err)
		}
		if err := c.reader.Resolve(); err != nil {
			t.Fatal(err)
		}
		msg := dynamicpb.NewMessage(desc)
		if err := c.Scan(msg); err != nil {
			t.Fatal(err)
		}
		if got := msg.Get(desc.Fields().ByName("db_host")).String(); got != "db" {
			t.Fatalf("db_host = %q, want db", got)
		}
	}
	if n := len(c.reader.(*reader).layers); n != 1 {
		t.Fatalf("retained %d layers, want 1", n)
	}
}

func TestScanKeepsWellKnownTypeKeys(t *testing.T) {
	fd, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name: proto.String("wkt_scan.proto"), Syntax: proto.String("proto3"), Package: proto.String("wkt_scan"),
		Dependency: []string{"google/protobuf/struct.proto"},
		MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("Settings"), Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name: proto.String("meta"), Number: proto.Int32(1),
				Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".google.protobuf.Struct"),
			},
			{
				Name: proto.String("val"), Number: proto.Int32(2),
				Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".google.protobuf.Value"),
			},
		}}},
	}, protoregistry.GlobalFiles)
	if err != nil {
		t.Fatal(err)
	}
	desc := fd.Messages().Get(0)
	c := New(WithSource(
		newTestJSONSource(`{"meta":{"fields":"a"},"val":{"string_value":"x"}}`),
		newTestJSONSource(`{"meta":{"apply_schema":true}}`),
	))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}
	msg := dynamicpb.NewMessage(desc)
	if err := c.Scan(msg); err != nil {
		t.Fatal(err)
	}
	// Nested well-known fields of a dynamic message are dynamic too; round-trip
	// them through protojson into the concrete types.
	var meta structpb.Struct
	var val structpb.Value
	for target, name := range map[proto.Message]protoreflect.Name{&meta: "meta", &val: "val"} {
		data, err := protojson.Marshal(msg.Get(desc.Fields().ByName(name)).Message().Interface())
		if err != nil {
			t.Fatal(err)
		}
		if err := protojson.Unmarshal(data, target); err != nil {
			t.Fatal(err)
		}
	}
	if meta.Fields["fields"].GetStringValue() != "a" || !meta.Fields["apply_schema"].GetBoolValue() {
		t.Fatalf("unexpected struct keys: %v", &meta)
	}
	if val.GetStructValue().GetFields()["string_value"].GetStringValue() != "x" {
		t.Fatalf("unexpected value keys: %v", &val)
	}
}

func TestCustomMergeDoesNotRetainLayers(t *testing.T) {
	c := New(WithMergeFunc(defaultMerge)).(*config)
	for i := 0; i < 100; i++ {
		if err := c.reader.Merge(&KeyValue{Format: "json", Value: []byte(`{"db_host":"${missing:db}"}`)}); err != nil {
			t.Fatal(err)
		}
		if err := c.reader.Resolve(); err != nil {
			t.Fatal(err)
		}
	}
	if n := len(c.reader.(*reader).layers); n != 0 {
		t.Fatalf("custom merge retained %d layers, want 0", n)
	}
	msg := dynamicpb.NewMessage(testResolutionDescriptor(t))
	if err := c.Scan(msg); err != nil {
		t.Fatal(err)
	}
	if got := msg.Get(msg.Descriptor().Fields().ByName("db_host")).String(); got != "db" {
		t.Fatalf("db_host = %q, want db", got)
	}
}

func TestReaderConcurrentMergeKeepsAllUpdates(t *testing.T) {
	r := newReader(options{decoder: defaultDecoder, resolver: defaultResolver, merge: defaultMerge}).(*reader)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := r.Merge(&KeyValue{Format: "json", Value: []byte(fmt.Sprintf(`{"k%d":%d}`, i, i))}); err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < 20; i++ {
		if _, ok := r.Value(fmt.Sprintf("k%d", i)); !ok {
			t.Fatalf("lost concurrent update k%d", i)
		}
	}
}
