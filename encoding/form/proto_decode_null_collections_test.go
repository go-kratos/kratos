package form

import (
	"net/url"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Complement the singular-field null tests by exercising the separate list and
// map assignment guards without generated fixtures or servers.
func TestDecodeValuesNullCollections(t *testing.T) {
	for _, tc := range []struct{ name, dependency, message, valid string }{
		{"Timestamp", "google/protobuf/timestamp.proto", ".google.protobuf.Timestamp", "2026-01-02T03:04:05Z"},
		{"Duration", "google/protobuf/duration.proto", ".google.protobuf.Duration", "1.250s"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optional := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
			repeated := descriptorpb.FieldDescriptorProto_LABEL_REPEATED
			message := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
			str := descriptorpb.FieldDescriptorProto_TYPE_STRING
			fd, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
				Name: proto.String("null_collection_" + tc.name + ".proto"), Package: proto.String("formregression"), Syntax: proto.String("proto3"), Dependency: []string{tc.dependency},
				MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("Container"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{Name: proto.String("items"), Number: proto.Int32(1), Label: &repeated, Type: &message, TypeName: proto.String(tc.message)},
						{Name: proto.String("entries"), Number: proto.Int32(2), Label: &repeated, Type: &message, TypeName: proto.String(".formregression.Container.EntriesEntry")},
					},
					NestedType: []*descriptorpb.DescriptorProto{{Name: proto.String("EntriesEntry"), Options: &descriptorpb.MessageOptions{MapEntry: proto.Bool(true)}, Field: []*descriptorpb.FieldDescriptorProto{
						{Name: proto.String("key"), Number: proto.Int32(1), Label: &optional, Type: &str},
						{Name: proto.String("value"), Number: proto.Int32(2), Label: &optional, Type: &message, TypeName: proto.String(tc.message)},
					}}},
				}},
			}, protoregistry.GlobalFiles)
			if err != nil {
				t.Fatal(err)
			}
			md := fd.Messages().ByName("Container")
			for _, field := range []string{"items", "entries"} {
				t.Run(field, func(t *testing.T) {
					msg := dynamicpb.NewMessage(md)
					// Catch a baseline panic so all four independent assignment paths run.
					defer func() {
						if r := recover(); r != nil {
							t.Errorf("DecodeValues panicked: %v", r)
						}
					}()
					var values url.Values
					if field == "items" {
						values = url.Values{"items": {"null", tc.valid, "null"}}
					} else {
						values = url.Values{"entries[absent]": {"null"}, "entries[kept]": {tc.valid}}
					}
					if err := DecodeValues(msg, values); err != nil {
						t.Fatal(err)
					}
					descriptor := md.Fields().ByName(protoreflect.Name(field))
					if field == "items" {
						list := msg.Get(descriptor).List()
						if list.Len() != 1 {
							t.Fatalf("want only valid list element, got %d", list.Len())
						}
						expected, err := parseMessage(descriptor.Message(), tc.valid)
						if err != nil {
							t.Fatal(err)
						}
						if !proto.Equal(list.Get(0).Message().Interface(), expected.Message().Interface()) {
							t.Fatal("valid element changed")
						}
					} else {
						entries := msg.Get(descriptor).Map()
						if entries.Len() != 1 || entries.Has(protoreflect.ValueOfString("absent").MapKey()) {
							t.Fatalf("null map entry retained: %v", entries)
						}
						got := entries.Get(protoreflect.ValueOfString("kept").MapKey())
						expected, err := parseMessage(descriptor.MapValue().Message(), tc.valid)
						if err != nil {
							t.Fatal(err)
						}
						if !got.IsValid() || !proto.Equal(got.Message().Interface(), expected.Message().Interface()) {
							t.Fatal("valid map entry changed")
						}
					}
				})
			}
		})
	}
}
