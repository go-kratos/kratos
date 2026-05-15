package http

import (
	"reflect"
	"regexp"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/go-kratos/kratos/v3/encoding/form"
)

var pathTemplateParamRE = regexp.MustCompile(`{([.\w]+)(=[^{}]*)?}`)

// BuildPathOption configures path construction.
type BuildPathOption func(*buildPathOptions)

type buildPathOptions struct {
	queryParams bool
	omitFields  []string
}

// WithQueryParams appends request fields that are not bound in the path as query parameters.
func WithQueryParams() BuildPathOption {
	return func(o *buildPathOptions) {
		o.queryParams = true
	}
}

// WithOmitFields excludes fields from generated query parameters.
func WithOmitFields(fields ...string) BuildPathOption {
	return func(o *buildPathOptions) {
		o.omitFields = append(o.omitFields, fields...)
	}
}

// BuildPath builds an HTTP request path from a path template and request message.
func BuildPath(pathTemplate string, msg any, opts ...BuildPathOption) string {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Pointer && reflect.ValueOf(msg).IsNil()) {
		return pathTemplate
	}

	options := buildPathOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	queryParams, _ := form.EncodeValues(msg)
	pathParams := make(map[string]struct{})
	path := pathTemplate
	if strings.ContainsRune(pathTemplate, '{') {
		path = pathTemplateParamRE.ReplaceAllStringFunc(pathTemplate, func(in string) string {
			matches := pathTemplateParamRE.FindStringSubmatch(in)
			key := matches[1]
			pathParams[key] = struct{}{}
			return queryParams.Get(key)
		})
	}

	if !options.queryParams {
		if v, ok := msg.(proto.Message); ok {
			if query := form.EncodeFieldMask(v.ProtoReflect()); query != "" {
				return path + "?" + query
			}
		}
		return path
	}
	if len(queryParams) > 0 {
		for key := range pathParams {
			delete(queryParams, key)
		}
		omitQueryParams(queryParams, options.omitFields)
		if query := queryParams.Encode(); query != "" {
			path += "?" + query
		}
	}
	return path
}

func omitQueryParams(values map[string][]string, fields []string) {
	for _, field := range fields {
		if field == "" {
			continue
		}
		delete(values, field)
		prefix := field + "."
		for key := range values {
			if strings.HasPrefix(key, prefix) {
				delete(values, key)
			}
		}
	}
}
