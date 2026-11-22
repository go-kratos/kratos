package main

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestSupportedFeaturesIncludeEditions(t *testing.T) {
	gen, err := protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{})
	if err != nil {
		t.Fatal(err)
	}
	setSupportedFeatures(gen)
	resp := gen.Response()

	want := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL |
		pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	if resp.GetSupportedFeatures() != want {
		t.Fatalf("supported features = %b, want %b", resp.GetSupportedFeatures(), want)
	}
	if got := descriptorpb.Edition(resp.GetMinimumEdition()); got != descriptorpb.Edition_EDITION_PROTO2 {
		t.Fatalf("minimum edition = %v, want EDITION_PROTO2", got)
	}
	if got := descriptorpb.Edition(resp.GetMaximumEdition()); got != descriptorpb.Edition_EDITION_2024 {
		t.Fatalf("maximum edition = %v, want EDITION_2024", got)
	}
}
