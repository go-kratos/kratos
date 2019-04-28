package main

import (
	"fmt"
	"net/http"
	"strings"

	"go-common/app/tool/bmproto/protoc-gen-bm/extensions/gogoproto"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func getMoreTags(field *descriptor.FieldDescriptorProto) *string {
	if field == nil {
		return nil
	}
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, gogoproto.E_Moretags)
		if err == nil && v.(*string) != nil {
			return v.(*string)
		}
	}
	return nil
}

func getJsonTag(field *descriptor.FieldDescriptorProto) string {
	if field == nil {
		return ""
	}
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, gogoproto.E_Jsontag)
		if err == nil && v.(*string) != nil {
			ret := *(v.(*string))
			i := strings.Index(ret, ",")
			if i != -1 {
				ret = ret[:i]
			}
			return ret
		}
	}
	return field.GetName()
}

type googleMethodOptionInfo struct {
	Method      string
	PathPattern string
	HTTPRule    *annotations.HttpRule
}

// ParseBMMethod parse BMMethodDescriptor form method descriptor proto
func ParseBMMethod(method *descriptor.MethodDescriptorProto) (*googleMethodOptionInfo, error) {
	ext, err := proto.GetExtension(method.GetOptions(), annotations.E_Http)
	if err != nil {
		return nil, fmt.Errorf("get extension error: %s", err)
	}
	rule := ext.(*annotations.HttpRule)
	var httpMethod string
	var pathPattern string
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		pathPattern = pattern.Get
		httpMethod = http.MethodGet
	case *annotations.HttpRule_Put:
		pathPattern = pattern.Put
		httpMethod = http.MethodPut
	case *annotations.HttpRule_Post:
		pathPattern = pattern.Post
		httpMethod = http.MethodPost
	case *annotations.HttpRule_Patch:
		pathPattern = pattern.Patch
		httpMethod = http.MethodPatch
	case *annotations.HttpRule_Delete:
		pathPattern = pattern.Delete
		httpMethod = http.MethodDelete
	default:
		return nil, fmt.Errorf("unsupport http pattern %s", rule.Pattern)
	}
	bmMethod := &googleMethodOptionInfo{
		Method:      httpMethod,
		PathPattern: pathPattern,
		HTTPRule:    rule,
	}
	return bmMethod, nil
}
