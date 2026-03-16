package base

import (
	"strings"
	"testing"
)

func TestReplaceTemplateContentPreservesProtobufRawDesc(t *testing.T) {
	const oldMod = "github.com/example/template"
	const newMod = "github.com/example/service"

	input := strings.Join([]string{
		"package conf",
		"",
		"import dep \"" + oldMod + "/api/greeter/v1\"",
		"",
		"const file_conf_conf_proto_rawDesc = \"\" +",
		"\t\"\\n\" +",
		"\t\"\\x0fconf/conf.proto\\x12\\n\" +",
		"\t\"kratos.apiB7Z5" + oldMod + "/internal/conf;confb\\x06proto3\"",
		"",
		"var wireImport = \"" + oldMod + "/internal/server\"",
		"",
	}, "\n")

	got := string(replaceTemplateContent([]byte(input), []string{oldMod, newMod}))

	if !strings.Contains(got, `import dep "`+newMod+`/api/greeter/v1"`) {
		t.Fatalf("expected import path to be replaced, got:\n%s", got)
	}
	if !strings.Contains(got, `var wireImport = "`+newMod+`/internal/server"`) {
		t.Fatalf("expected regular string to be replaced, got:\n%s", got)
	}
	if !strings.Contains(got, `B7Z5`+oldMod+`/internal/conf;confb\x06proto3"`) {
		t.Fatalf("expected protobuf raw descriptor to stay untouched, got:\n%s", got)
	}
	if strings.Contains(got, `B7Z5`+newMod+`/internal/conf;confb\x06proto3"`) {
		t.Fatalf("protobuf raw descriptor was unexpectedly rewritten:\n%s", got)
	}
}

func TestReplaceTemplateContentWithoutRawDescReplacesEverywhere(t *testing.T) {
	const oldMod = "github.com/example/template"
	const newMod = "github.com/example/service"

	input := []byte(`package main

import _ "` + oldMod + `/internal/biz"

var modulePath = "` + oldMod + `"
`)

	got := string(replaceTemplateContent(input, []string{oldMod, newMod}))

	if strings.Contains(got, oldMod) {
		t.Fatalf("expected all occurrences to be replaced, got:\n%s", got)
	}
	if !strings.Contains(got, newMod) {
		t.Fatalf("expected replacement to be applied, got:\n%s", got)
	}
}

func TestReplaceTemplateContentReplacesLegacyProtobufByteDesc(t *testing.T) {
	const oldMod = "github.com/example/template"
	const newMod = "github.com/example/service"

	input := strings.Join([]string{
		"package conf",
		"",
		"var file_conf_conf_proto_rawDesc = []byte{",
		"\t0x42, 0x37, 0x5a, 0x35, // B7Z5",
		"}",
		"",
		"const legacyPath = \"" + oldMod + "/internal/conf;conf\"",
		"",
	}, "\n")

	got := string(replaceTemplateContent([]byte(input), []string{oldMod, newMod}))

	if !strings.Contains(got, `const legacyPath = "`+newMod+`/internal/conf;conf"`) {
		t.Fatalf("expected legacy protobuf file content to still be replaced, got:\n%s", got)
	}
	if strings.Contains(got, `const legacyPath = "`+oldMod+`/internal/conf;conf"`) {
		t.Fatalf("legacy protobuf replacement did not occur, got:\n%s", got)
	}
}
