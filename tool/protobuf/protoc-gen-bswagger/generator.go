package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/bilibili/kratos/tool/protobuf/pkg/gen"
	"github.com/bilibili/kratos/tool/protobuf/pkg/generator"
	"github.com/bilibili/kratos/tool/protobuf/pkg/naming"
	"github.com/bilibili/kratos/tool/protobuf/pkg/tag"
	"github.com/bilibili/kratos/tool/protobuf/pkg/typemap"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type swaggerGen struct {
	generator.Base
	// defsMap will fill into swagger's definitions
	// key is full qualified proto name
	defsMap map[string]*typemap.MessageDefinition
}

// NewSwaggerGenerator a swagger generator
func NewSwaggerGenerator() *swaggerGen {
	return &swaggerGen{}
}

func (t *swaggerGen) Generate(in *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	t.Setup(in)
	resp := &plugin.CodeGeneratorResponse{}
	for _, f := range t.GenFiles {
		if len(f.Service) == 0 {
			continue
		}
		respFile := t.generateSwagger(f)
		if respFile != nil {
			resp.File = append(resp.File, respFile)
		}
	}
	return resp
}

func (t *swaggerGen) generateSwagger(file *descriptor.FileDescriptorProto) *plugin.CodeGeneratorResponse_File {
	var pkg = file.GetPackage()
	r := regexp.MustCompile("v(\\d+)$")
	strs := r.FindStringSubmatch(pkg)
	var vStr string
	if len(strs) >= 2 {
		vStr = strs[1]
	} else {
		vStr = ""
	}
	var swaggerObj = &swaggerObject{
		Paths:   swaggerPathsObject{},
		Swagger: "2.0",
		Info: swaggerInfoObject{
			Title:   file.GetName(),
			Version: vStr,
		},
		Schemes:  []string{"http", "https"},
		Consumes: []string{"application/json", "multipart/form-data"},
		Produces: []string{"application/json"},
	}
	t.defsMap = map[string]*typemap.MessageDefinition{}

	out := &plugin.CodeGeneratorResponse_File{}
	name := naming.GenFileName(file, ".swagger.json")
	for _, svc := range file.Service {
		for _, meth := range svc.Method {
			if !t.ShouldGenForMethod(file, svc, meth) {
				continue
			}
			apiInfo := t.GetHttpInfoCached(file, svc, meth)
			pathItem := swaggerPathItemObject{}
			if originPathItem, ok := swaggerObj.Paths[apiInfo.Path]; ok {
				pathItem = originPathItem
			}
			op := t.getOperationByHTTPMethod(apiInfo.HttpMethod, &pathItem)
			op.Summary = apiInfo.Title
			op.Description = apiInfo.Description
			swaggerObj.Paths[apiInfo.Path] = pathItem
			op.Tags = []string{pkg + "." + svc.GetName()}

			// request
			request := t.Reg.MessageDefinition(meth.GetInputType())
			// request cannot represent by simple form
			isComplexRequest := false
			for _, field := range request.Descriptor.Field {
				if !generator.IsScalar(field) {
					isComplexRequest = true
					break
				}
			}
			if !isComplexRequest && apiInfo.HttpMethod == "GET" {
				for _, field := range request.Descriptor.Field {
					if !generator.IsScalar(field) {
						continue
					}
					p := t.getQueryParameter(file, request, field)
					op.Parameters = append(op.Parameters, p)
				}
			} else {
				p := swaggerParameterObject{}
				p.In = "body"
				p.Required = true
				p.Name = "body"
				p.Schema = &swaggerSchemaObject{}
				p.Schema.Ref = "#/definitions/" + meth.GetInputType()
				op.Parameters = []swaggerParameterObject{p}
			}

			// response
			resp := swaggerResponseObject{}
			resp.Description = "A successful response."

			// proto 里面的response只定义data里面的
			// 所以需要把code msg data 这一级加上
			resp.Schema.Type = "object"
			resp.Schema.Properties = &swaggerSchemaObjectProperties{}
			p := keyVal{Key: "code", Value: &schemaCore{Type: "integer"}}
			*resp.Schema.Properties = append(*resp.Schema.Properties, p)
			p = keyVal{Key: "message", Value: &schemaCore{Type: "string"}}
			*resp.Schema.Properties = append(*resp.Schema.Properties, p)
			p = keyVal{Key: "data", Value: schemaCore{Ref: "#/definitions/" + meth.GetOutputType()}}
			*resp.Schema.Properties = append(*resp.Schema.Properties, p)
			op.Responses = swaggerResponsesObject{"200": resp}
		}
	}

	// walk though definitions
	t.walkThroughFileDefinition(file)
	defs := swaggerDefinitionsObject{}
	swaggerObj.Definitions = defs
	for typ, msg := range t.defsMap {
		def := swaggerSchemaObject{}
		def.Properties = new(swaggerSchemaObjectProperties)
		def.Description = strings.Trim(msg.Comments.Leading, "\n\r ")
		for _, field := range msg.Descriptor.Field {
			p := keyVal{Key: generator.GetFormOrJSONName(field)}
			schema := t.schemaForField(file, msg, field)
			if generator.GetFieldRequired(field, t.Reg, msg) {
				def.Required = append(def.Required, p.Key)
			}
			p.Value = schema
			*def.Properties = append(*def.Properties, p)
		}
		def.Type = "object"
		defs[typ] = def
	}
	b, _ := json.MarshalIndent(swaggerObj, "", "    ")
	str := string(b)
	out.Name = &name
	out.Content = &str
	return out
}

func (t *swaggerGen) getOperationByHTTPMethod(httpMethod string, pathItem *swaggerPathItemObject) *swaggerOperationObject {
	var op = &swaggerOperationObject{}
	switch httpMethod {
	case http.MethodGet:
		pathItem.Get = op
	case http.MethodPost:
		pathItem.Post = op
	case http.MethodPut:
		pathItem.Put = op
	case http.MethodDelete:
		pathItem.Delete = op
	case http.MethodPatch:
		pathItem.Patch = op
	default:
		pathItem.Get = op
	}
	return op
}

func (t *swaggerGen) getQueryParameter(file *descriptor.FileDescriptorProto,
	input *typemap.MessageDefinition,
	field *descriptor.FieldDescriptorProto) swaggerParameterObject {
	p := swaggerParameterObject{}
	p.Name = generator.GetFormOrJSONName(field)
	fComment, _ := t.Reg.FieldComments(input, field)
	cleanComment := tag.GetCommentWithoutTag(fComment.Leading)

	p.Description = strings.Trim(strings.Join(cleanComment, "\n"), "\n\r ")
	p.In = "query"
	p.Required = generator.GetFieldRequired(field, t.Reg, input)
	typ, isArray, format := getFieldSwaggerType(field)
	if isArray {
		p.Items = &swaggerItemsObject{}
		p.Type = "array"
		p.Items.Type = typ
		p.Items.Format = format
	} else {
		p.Type = typ
		p.Format = format
	}
	return p
}

func (t *swaggerGen) schemaForField(file *descriptor.FileDescriptorProto,
	msg *typemap.MessageDefinition,
	field *descriptor.FieldDescriptorProto) swaggerSchemaObject {
	schema := swaggerSchemaObject{}
	fComment, err := t.Reg.FieldComments(msg, field)
	if err != nil {
		gen.Error(err, "comment not found err %+v")
	}
	schema.Description = strings.Trim(fComment.Leading, "\n\r ")
	typ, isArray, format := getFieldSwaggerType(field)
	if !generator.IsScalar(field) {
		if generator.IsMap(field, t.Reg) {
			schema.Type = "object"
			mapMsg := t.Reg.MessageDefinition(field.GetTypeName())
			mapValueField := mapMsg.Descriptor.Field[1]
			valSchema := t.schemaForField(file, mapMsg, mapValueField)
			schema.AdditionalProperties = &valSchema
		} else {
			if isArray {
				schema.Items = &swaggerItemsObject{}
				schema.Type = "array"
				schema.Items.Ref = "#/definitions/" + field.GetTypeName()
			} else {
				schema.Ref = "#/definitions/" + field.GetTypeName()
			}
		}
	} else {
		if isArray {
			schema.Items = &swaggerItemsObject{}
			schema.Type = "array"
			schema.Items.Type = typ
			schema.Items.Format = format
		} else {
			schema.Type = typ
			schema.Format = format
		}
	}
	return schema
}

func (t *swaggerGen) walkThroughFileDefinition(file *descriptor.FileDescriptorProto) {
	for _, svc := range file.Service {
		for _, meth := range svc.Method {
			shouldGen := t.ShouldGenForMethod(file, svc, meth)
			if !shouldGen {
				continue
			}
			t.walkThroughMessages(t.Reg.MessageDefinition(meth.GetOutputType()))
			t.walkThroughMessages(t.Reg.MessageDefinition(meth.GetInputType()))
		}
	}
}

func (t *swaggerGen) walkThroughMessages(msg *typemap.MessageDefinition) {
	_, ok := t.defsMap[msg.ProtoName()]
	if ok {
		return
	}
	if !msg.Descriptor.GetOptions().GetMapEntry() {
		t.defsMap[msg.ProtoName()] = msg
	}
	for _, field := range msg.Descriptor.Field {
		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			t.walkThroughMessages(t.Reg.MessageDefinition(field.GetTypeName()))
		}
	}
}

func getFieldSwaggerType(field *descriptor.FieldDescriptorProto) (typeName string, isArray bool, formatName string) {
	typeName = "unknown"
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		typeName = "boolean"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		typeName = "number"
		formatName = "double"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		typeName = "number"
		formatName = "float"
	case
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		typeName = "integer"
	case
		descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_BYTES:
		typeName = "string"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		typeName = "object"
	}
	if field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		isArray = true
	}
	return
}
