package generator

import (
	"reflect"
	"strings"

	"github.com/bilibili/kratos/tool/protobuf/pkg/extensions/gogoproto"
	"github.com/bilibili/kratos/tool/protobuf/pkg/tag"
	"github.com/bilibili/kratos/tool/protobuf/pkg/typemap"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// GetJSONFieldName get name from gogoproto.jsontag
// else the original name
func GetJSONFieldName(field *descriptor.FieldDescriptorProto) string {
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

// GetFormOrJSONName get name from form tag, then json tag
// then original name
func GetFormOrJSONName(field *descriptor.FieldDescriptorProto) string {
	if field == nil {
		return ""
	}
	tags := tag.GetMoreTags(field)
	if tags != nil {
		tag := reflect.StructTag(*tags)
		fName := tag.Get("form")
		if fName != "" {
			i := strings.Index(fName, ",")
			if i != -1 {
				fName = fName[:i]
			}
			return fName
		}
	}
	return GetJSONFieldName(field)
}

// IsScalar Is this field a scalar numeric type?
func IsScalar(field *descriptor.FieldDescriptorProto) bool {
	if field.Type == nil {
		return false
	}
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_BYTES,
		descriptor.FieldDescriptorProto_TYPE_STRING:
		return true
	default:
		return false
	}
}

// IsMap is protocol buffer map
func IsMap(field *descriptor.FieldDescriptorProto, reg *typemap.Registry) bool {
	if field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return false
	}
	md := reg.MessageDefinition(field.GetTypeName())
	if md == nil || !md.Descriptor.GetOptions().GetMapEntry() {
		return false
	}
	return true
}

// IsRepeated Is this field repeated?
func IsRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED

}

// GetFieldRequired is field required?
// eg. validate="required"
func GetFieldRequired(
	f *descriptor.FieldDescriptorProto,
	reg *typemap.Registry,
	md *typemap.MessageDefinition,
) bool {
	fComment, _ := reg.FieldComments(md, f)
	var tags []reflect.StructTag
	{
		//get required info from gogoproto.moretags
		moretags := tag.GetMoreTags(f)
		if moretags != nil {
			tags = []reflect.StructTag{reflect.StructTag(*moretags)}
		}
	}
	if len(tags) == 0 {
		tags = tag.GetTagsInComment(fComment.Leading)
	}
	validateTag := tag.GetTagValue("validate", tags)
	var validateRules []string
	if validateTag != "" {
		validateRules = strings.Split(validateTag, ",")
	}
	required := false
	for _, rule := range validateRules {
		if rule == "required" {
			required = true
		}
	}
	return required
}

func MakeIndentStr(i int) string {
	return strings.Repeat(" ", i)
}
