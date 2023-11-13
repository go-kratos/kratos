package form

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/go-kratos/kratos/v2/internal/testdata/complex"
)

func TestDecodeValues(t *testing.T) {
	form, err := url.ParseQuery("a=19&age=18&b=true&bool=false&byte=MTIz&bytes=MTIz&count=3&d=22.22&double=12.33&duration=" +
		"2m0.000000022s&field=1%2C2&float=12.34&id=2233&int32=32&int64=64&numberOne=2233&price=11.23&sex=woman&simples=3344&" +
		"simples=5566&string=go-kratos&timestamp=1970-01-01T00%3A00%3A20.000000002Z&uint32=32&uint64=64&very_simple.component=5566")
	if err != nil {
		t.Fatal(err)
	}

	comp := &complex.Complex{}
	err = DecodeValues(comp, form)
	if err != nil {
		t.Fatal(err)
	}
	if comp.Id != int64(2233) {
		t.Errorf("want %v, got %v", int64(2233), comp.Id)
	}
	if comp.NoOne != "2233" {
		t.Errorf("want %v, got %v", "2233", comp.NoOne)
	}
	if comp.Simple == nil {
		t.Fatalf("want %v, got %v", nil, comp.Simple)
	}
	if comp.Simple.Component != "5566" {
		t.Errorf("want %v, got %v", "5566", comp.Simple.Component)
	}
	if len(comp.Simples) != 2 {
		t.Fatalf("want %v, got %v", 2, len(comp.Simples))
	}
	if comp.Simples[0] != "3344" {
		t.Errorf("want %v, got %v", "3344", comp.Simples[0])
	}
	if comp.Simples[1] != "5566" {
		t.Errorf("want %v, got %v", "5566", comp.Simples[1])
	}
}

func TestGetFieldDescriptor(t *testing.T) {
	comp := &complex.Complex{}

	field := getFieldDescriptor(comp.ProtoReflect(), "id")
	if field.Kind() != protoreflect.Int64Kind {
		t.Errorf("want: %d, got: %d", protoreflect.Int64Kind, field.Kind())
	}

	field = getFieldDescriptor(comp.ProtoReflect(), "simples")
	if field.Kind() != protoreflect.StringKind {
		t.Errorf("want: %d, got: %d", protoreflect.StringKind, field.Kind())
	}
}

func TestPopulateRepeatedField(t *testing.T) {
	query, err := url.ParseQuery("simples=3344&simples=5566")
	if err != nil {
		t.Fatal(err)
	}
	comp := &complex.Complex{}
	field := getFieldDescriptor(comp.ProtoReflect(), "simples")

	err = populateRepeatedField(field, comp.ProtoReflect().Mutable(field).List(), query["simples"])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"3344", "5566"}, comp.GetSimples()) {
		t.Errorf("want: %v, got: %v", []string{"3344", "5566"}, comp.GetSimples())
	}
}

func TestPopulateMapField(t *testing.T) {
	query, err := url.ParseQuery("map%5Bkratos%5D=https://go-kratos.dev/")
	if err != nil {
		t.Fatal(err)
	}
	comp := &complex.Complex{}
	field := getFieldDescriptor(comp.ProtoReflect(), "map")
	// Fill the comp map field with the url query values
	err = populateMapField(field, comp.ProtoReflect().Mutable(field).Map(), []string{"map[kratos]"}, query["map[kratos]"])
	if err != nil {
		t.Fatal(err)
	}
	// Get the comp map field value
	if query["map[kratos]"][0] != comp.Map["kratos"] {
		t.Errorf("want: %s, got: %s", query["map[kratos]"], comp.Map["kratos"])
	}
}

func TestPopulateMapSepField(t *testing.T) {
	query, err := url.ParseQuery("map.name=kratos")
	if err != nil {
		t.Fatal(err)
	}
	comp := &complex.Complex{}
	field := getFieldDescriptor(comp.ProtoReflect(), "map")
	// Fill the comp map field with the url query values
	err = populateMapField(field, comp.ProtoReflect().Mutable(field).Map(), []string{"map.name"}, query["map.name"])
	if err != nil {
		t.Fatal(err)
	}
	// Get the comp map field value
	if query["map.name"][0] != comp.Map["name"] {
		t.Errorf("want: %s, got: %s", query, comp.Map)
	}
}

func TestParseField(t *testing.T) {
	tests := []struct {
		name                    string
		fieldName               string
		protoReflectKind        protoreflect.Kind
		value                   string
		targetProtoReflectValue protoreflect.Value
		targetErr               error
	}{
		{
			name:                    "BoolKind",
			fieldName:               "b",
			protoReflectKind:        protoreflect.BoolKind,
			value:                   "true",
			targetProtoReflectValue: protoreflect.ValueOfBool(true),
			targetErr:               nil,
		},
		{
			name:                    "BoolKind",
			fieldName:               "b",
			protoReflectKind:        protoreflect.BoolKind,
			value:                   "a",
			targetProtoReflectValue: protoreflect.Value{},
			targetErr:               &strconv.NumError{Func: "ParseBool", Num: "a", Err: strconv.ErrSyntax},
		},
		{
			name:                    "EnumKind",
			fieldName:               "sex",
			protoReflectKind:        protoreflect.EnumKind,
			value:                   "1",
			targetProtoReflectValue: protoreflect.ValueOfEnum(1),
			targetErr:               nil,
		},
		{
			name:                    "EnumKind",
			fieldName:               "sex",
			protoReflectKind:        protoreflect.EnumKind,
			value:                   "2",
			targetProtoReflectValue: protoreflect.Value{},
			targetErr:               fmt.Errorf("%q is not a valid value", "2"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			comp := &complex.Complex{}
			field := getFieldDescriptor(comp.ProtoReflect(), test.fieldName)
			if test.protoReflectKind != field.Kind() {
				t.Fatalf("want: %d, got: %d", test.protoReflectKind, field.Kind())
			}
			val, err := parseField(field, test.value)
			if !reflect.DeepEqual(test.targetErr, err) {
				t.Fatalf("want: %s, got: %s", test.targetErr, err)
			}
			if !reflect.DeepEqual(test.targetProtoReflectValue, val) {
				t.Errorf("want: %s, got: %s", test.targetProtoReflectValue, val)
			}
		})
	}
}

func TestJsonSnakeCase(t *testing.T) {
	tests := []struct {
		camelCase string
		snakeCase string
	}{
		{
			"userId", "user_id",
		},
		{
			"user", "user",
		},
		{
			"userIdAndUsername", "user_id_and_username",
		},
		{
			"", "",
		},
	}
	for _, test := range tests {
		t.Run(test.camelCase, func(t *testing.T) {
			snake := jsonSnakeCase(test.camelCase)
			if snake != test.snakeCase {
				t.Errorf("want: %s, got: %s", test.snakeCase, snake)
			}
		})
	}
}

func TestIsASCIIUpper(t *testing.T) {
	tests := []struct {
		b     byte
		upper bool
	}{
		{
			'A', true,
		},
		{
			'a', false,
		},
		{
			',', false,
		},
		{
			'1', false,
		},
		{
			' ', false,
		},
	}
	for _, test := range tests {
		t.Run(string(test.b), func(t *testing.T) {
			upper := isASCIIUpper(test.b)
			if test.upper != upper {
				t.Errorf("'%s' is not ascii upper", string(test.b))
			}
		})
	}
}

func TestParseURLQueryMapKey(t *testing.T) {
	tests := []struct {
		fieldName string
		field     string
		fieldKey  string
		err       error
	}{
		{
			fieldName: "map[kratos]", field: "map", fieldKey: "kratos", err: nil,
		},
		{
			fieldName: "map[]", field: "map", fieldKey: "", err: nil,
		},
		{
			fieldName: "", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
		{
			fieldName: "[[]", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
		{
			fieldName: "map[kratos]=", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
		{
			fieldName: "[kratos]", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
		{
			fieldName: "map", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
		{
			fieldName: "map[", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
		{
			fieldName: "]kratos[", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
		{
			fieldName: "[kratos", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
		{
			fieldName: "kratos]", field: "", fieldKey: "", err: errInvalidFormatMapKey,
		},
	}
	for _, test := range tests {
		t.Run(test.fieldName, func(t *testing.T) {
			fieldName, fieldKey, err := parseURLQueryMapKey(test.fieldName)
			if test.err != err {
				t.Fatalf("want: %s, got: %s", test.err, err)
			}
			if test.field != fieldName {
				t.Errorf("want: %s, got: %s", test.field, fieldName)
			}
			if test.fieldKey != fieldKey {
				t.Errorf("want: %s, got: %s", test.fieldKey, fieldKey)
			}
		})
	}
}
