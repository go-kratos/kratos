package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestNoParameters(t *testing.T) {
	path := "/test/noparams"
	m := buildPathVars(path)
	if !reflect.DeepEqual(m, map[string]*string{}) {
		t.Fatalf("Map should be empty")
	}
}

func TestSingleParam(t *testing.T) {
	path := "/test/{message.id}"
	m := buildPathVars(path)
	if !reflect.DeepEqual(len(m), 1) {
		t.Fatalf("len(m) not is 1")
	}
	if m["message.id"] != nil {
		t.Fatalf(`m["message.id"] should be empty`)
	}
}

func TestTwoParametersReplacement(t *testing.T) {
	path := "/test/{message.id}/{message.name=messages/*}"
	m := buildPathVars(path)
	if len(m) != 2 {
		t.Fatal("len(m) should be 2")
	}
	if m["message.id"] != nil {
		t.Fatal(`m["message.id"] should be nil`)
	}
	if m["message.name"] == nil {
		t.Fatal(`m["message.name"] should not be nil`)
	}
	if *m["message.name"] != "messages/*" {
		t.Fatal(`m["message.name"] should be "messages/*"`)
	}
}

func TestNoReplacePath(t *testing.T) {
	path := "/test/{message.id=test}"
	if !reflect.DeepEqual(replacePath("message.id", "test", path), "/test/{message.id:test}") {
		t.Fatal(`replacePath("message.id", "test", path) should be "/test/{message.id:test}"`)
	}
	path = "/test/{message.id=test/*}"
	if !reflect.DeepEqual(replacePath("message.id", "test/*", path), "/test/{message.id:test/[^/]+}") {
		t.Fatal(`replacePath("message.id", "test/*", path) should be "/test/{message.id:test/[^/]+}"`)
	}
}

func TestReplacePath(t *testing.T) {
	path := "/test/{message.id}/{message.name=messages/*}"
	newPath := replacePath("message.name", "messages/*", path)
	if !reflect.DeepEqual("/test/{message.id}/{message.name:messages/[^/]+}", newPath) {
		t.Fatal(`replacePath("message.name", "messages/*", path) should be "/test/{message.id}/{message.name:messages/[^/]+}"`)
	}
}

func TestIteration(t *testing.T) {
	path := "/test/{message.id}/{message.name=messages/*}"
	vars := buildPathVars(path)
	for v, s := range vars {
		if s != nil {
			path = replacePath(v, *s, path)
		}
	}
	if !reflect.DeepEqual("/test/{message.id}/{message.name:messages/[^/]+}", path) {
		t.Fatal(`replacePath("message.name", "messages/*", path) should be "/test/{message.id}/{message.name:messages/[^/]+}"`)
	}
}

func TestIterationMiddle(t *testing.T) {
	path := "/test/{message.name=messages/*}/books"
	vars := buildPathVars(path)
	for v, s := range vars {
		if s != nil {
			path = replacePath(v, *s, path)
		}
	}
	if !reflect.DeepEqual("/test/{message.name:messages/[^/]+}/books", path) {
		t.Fatal(`replacePath("message.name", "messages/*", path) should be "/test/{message.name:messages/[^/]+}/books"`)
	}
}

func TestReplaceBoundary(t *testing.T) {
	path := "/test/{message.namespace=*}/name/{message.name=*}"
	vars := buildPathVars(path)
	for v, s := range vars {
		if s != nil {
			path = replacePath(v, *s, path)
		}
	}
	if !reflect.DeepEqual("/test/{message.namespace:[^/]+}/name/{message.name:[^/]+}", path) {
		t.Fatal(`"/test/{message.namespace=*}/name/{message.name=*}" should be "/test/{message.namespace:[^/]+}/name/{message.name:[^/]+}"`)
	}
}

func TestPathTemplateRegex(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "single segment",
			value: "messages/*",
			want:  "messages/[^/]+",
		},
		{
			name:  "multi segment",
			value: "messages/**",
			want:  "messages/.*",
		},
		{
			name:  "literal",
			value: "v1.0/*",
			want:  `v1\.0/[^/]+`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathTemplateRegex(tt.value); got != tt.want {
				t.Errorf("expected %s got %s", tt.want, got)
			}
		})
	}
}

func TestHTTPTemplateClientUsesBuildPathAndProtoJSONHeaders(t *testing.T) {
	sd := &serviceDesc{
		ServiceType: "Greeter",
		ServiceName: "helloworld.Greeter",
		Methods: []*methodDesc{
			{
				Name:         "SayHello",
				OriginalName: "SayHello",
				Request:      "HelloRequest",
				Reply:        "HelloReply",
				Path:         "/helloworld/{name}",
				PathTemplate: "/helloworld/{name}",
				Method:       "GET",
				HasVars:      true,
			},
			{
				Name:         "CreateHello",
				OriginalName: "CreateHello",
				Request:      "CreateHelloRequest",
				Reply:        "HelloReply",
				Path:         "/helloworld",
				PathTemplate: "/helloworld",
				Method:       "POST",
				HasBody:      true,
				Body:         "",
			},
		},
	}
	got := sd.execute()
	for _, want := range []string{
		`path := http.BuildPath(pattern, in, http.WithQueryParams())`,
		`path := http.BuildPath(pattern, in)`,
		`http.Accept("application/protojson")`,
		`http.ContentType("application/protojson")`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated template missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "binding.") {
		t.Fatalf("generated template should not reference binding package:\n%s", got)
	}
}
