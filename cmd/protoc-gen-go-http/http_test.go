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
	if !reflect.DeepEqual(replacePath("message.id", "test/*", path), "/test/{message.id:test/.*}") {
		t.Fatal(`replacePath("message.id", "test/*", path) should be "/test/{message.id:test/.*}"`)
	}
}

func TestReplacePath(t *testing.T) {
	path := "/test/{message.id}/{message.name=messages/*}"
	newPath := replacePath("message.name", "messages/*", path)
	if !reflect.DeepEqual("/test/{message.id}/{message.name:messages/.*}", newPath) {
		t.Fatal(`replacePath("message.name", "messages/*", path) should be "/test/{message.id}/{message.name:messages/.*}"`)
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
	if !reflect.DeepEqual("/test/{message.id}/{message.name:messages/.*}", path) {
		t.Fatal(`replacePath("message.name", "messages/*", path) should be "/test/{message.id}/{message.name:messages/.*}"`)
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
	if !reflect.DeepEqual("/test/{message.name:messages/.*}/books", path) {
		t.Fatal(`replacePath("message.name", "messages/*", path) should be "/test/{message.name:messages/.*}/books"`)
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
	if !reflect.DeepEqual("/test/{message.namespace:.*}/name/{message.name:.*}", path) {
		t.Fatal(`"/test/{message.namespace=*}/name/{message.name=*}" should be "/test/{message.namespace:.*}/name/{message.name:.*}"`)
	}
}

func TestGetterAccessor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{".Payload", ".GetPayload()"},
		{".UserProfile", ".GetUserProfile()"},
		{".User.Profile", ".GetUser().GetProfile()"},
	}
	for _, tt := range tests {
		got := toGetterAccessor(tt.input)
		if got != tt.expected {
			t.Fatalf("toGetterAccessor(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestTemplateBodyNamedField(t *testing.T) {
	sd := &serviceDesc{
		ServiceType: "TestSvc",
		ServiceName: "test.v1.TestSvc",
		Metadata:    "test/v1/test.proto",
		Methods: []*methodDesc{
			{
				Name:          "BotHandler",
				OriginalName: "BotHandler",
				Num:           0,
				Request:       "BotHandlerRequest",
				Reply:         "emptypb.Empty",
				Path:          "/api/v1/bot",
				Method:        "POST",
				HasBody:       true,
				Body:          ".Payload",
				BodyProtoName: "payload",
			},
		},
	}
	output := sd.execute()

	// Server handler should use proto reflection for named body field
	if !strings.Contains(output, `in.ProtoReflect().Mutable(in.ProtoReflect().Descriptor().Fields().ByName("payload")).Message().Interface()`) {
		t.Fatal("server handler should use ProtoReflect().Mutable() for named body field")
	}
	if strings.Contains(output, "&in.Payload") {
		t.Fatal("server handler should NOT use direct field access &in.Payload")
	}

	// Client should use getter for named body field
	if !strings.Contains(output, "in.GetPayload()") {
		t.Fatal("client should use getter in.GetPayload() for named body field")
	}
	if strings.Contains(output, "in.Payload") && !strings.Contains(output, "in.GetPayload()") {
		t.Fatal("client should NOT use direct field access in.Payload")
	}
}

func TestTemplateBodyStar(t *testing.T) {
	sd := &serviceDesc{
		ServiceType: "TestSvc",
		ServiceName: "test.v1.TestSvc",
		Metadata:    "test/v1/test.proto",
		Methods: []*methodDesc{
			{
				Name:         "Create",
				OriginalName: "Create",
				Num:          0,
				Request:      "CreateRequest",
				Reply:        "CreateReply",
				Path:         "/api/v1/create",
				Method:       "POST",
				HasBody:      true,
				Body:         "",
			},
		},
	}
	output := sd.execute()

	// body="*" should use &in for server bind and in for client
	if !strings.Contains(output, "ctx.Bind(&in)") {
		t.Fatal("body=* should use ctx.Bind(&in)")
	}
	if strings.Contains(output, "ProtoReflect") {
		t.Fatal("body=* should NOT use ProtoReflect")
	}
}
