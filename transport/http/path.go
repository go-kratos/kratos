package http

import (
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

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
// Protobuf fields may use their proto names or JSON names.
func WithOmitFields(fields ...string) BuildPathOption {
	return func(o *buildPathOptions) {
		o.omitFields = append(o.omitFields, fields...)
	}
}

// BuildPath builds an HTTP request path from a path template and request message.
// Protobuf field paths may use proto names or JSON names. Values are escaped as
// single segments, except resource-name templates which preserve slash separators.
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
			key := queryFieldName(msg, matches[1])
			pathParams[key] = struct{}{}
			return escapePathValue(queryParams.Get(key), matches[2])
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
		omitQueryParams(queryParams, options.omitFields, msg)
		if query := queryParams.Encode(); query != "" {
			path += "?" + query
		}
	}
	return path
}

// queryFieldName resolves proto field paths to the names used by the form encoder.
func queryFieldName(msg any, path string) string {
	message, ok := msg.(proto.Message)
	if !ok {
		return path
	}
	fields := message.ProtoReflect().Descriptor().Fields()
	parts := strings.Split(path, ".")
	for i, name := range parts {
		field := fields.ByName(protoreflect.Name(name))
		if field == nil {
			field = fields.ByJSONName(name)
		}
		if field == nil {
			return path
		}
		parts[i] = field.TextName()
		if field.HasJSONName() {
			parts[i] = field.JSONName()
		}
		if i < len(parts)-1 {
			if field.Message() == nil {
				return path
			}
			fields = field.Message().Fields()
		}
	}
	return strings.Join(parts, ".")
}

func escapePathValue(value, template string) string {
	parts := []string{value}
	// Resource-name templates preserve their hierarchy; a plain variable or *
	// represents one segment and must escape embedded slashes.
	if strings.Contains(template, "/") || strings.Contains(template, "**") {
		parts = strings.Split(value, "/")
	}
	for i, part := range parts {
		if part == "." || part == ".." {
			parts[i] = strings.Repeat("%2E", len(part))
		} else {
			parts[i] = url.PathEscape(part)
		}
	}
	return strings.Join(parts, "/")
}

func omitQueryParams(values map[string][]string, fields []string, msg any) {
	for _, field := range fields {
		if field == "" {
			continue
		}
		field = queryFieldName(msg, field)
		delete(values, field)
		prefix := field + "."
		for key := range values {
			if strings.HasPrefix(key, prefix) {
				delete(values, key)
			}
		}
	}
}
