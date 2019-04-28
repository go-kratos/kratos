package generator_test

import (
	"bytes"
	"strings"
	"testing"

	"go-common/app/tool/gengo/generator"
	"go-common/app/tool/gengo/namer"
	"go-common/app/tool/gengo/parser"
)

func construct(t *testing.T, files map[string]string) *generator.Context {
	b := parser.New()
	for name, src := range files {
		if err := b.AddFileForTest("/tmp/"+name, name, []byte(src)); err != nil {
			t.Fatal(err)
		}
	}
	c, err := generator.NewContext(b, namer.NameSystems{
		"public":  namer.NewPublicNamer(0),
		"private": namer.NewPrivateNamer(0),
	}, "public")
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSnippetWriter(t *testing.T) {
	var structTest = map[string]string{
		"base/foo/proto/foo.go": `
package foo

// Blah is a test.
// A test, I tell you.
type Blah struct {
	// A is the first field.
	A int64 ` + "`" + `json:"a"` + "`" + `

	// B is the second field.
	// Multiline comments work.
	B string ` + "`" + `json:"b"` + "`" + `
}
`,
	}

	c := construct(t, structTest)
	b := &bytes.Buffer{}
	err := generator.NewSnippetWriter(b, c, "$", "$").
		Do("$.|public$$.|private$", c.Order[0]).
		Error()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if e, a := "Blahblah", b.String(); e != a {
		t.Errorf("Expected %q, got %q", e, a)
	}

	err = generator.NewSnippetWriter(b, c, "$", "$").
		Do("$.|public", c.Order[0]).
		Error()
	if err == nil {
		t.Errorf("expected error on invalid template")
	} else {
		// Dear reader, I apologize for making the worst change
		// detection test in the history of ever.
		if e, a := "snippet_writer_test.go", err.Error(); !strings.Contains(a, e) {
			t.Errorf("Expected %q but didn't find it in %q", e, a)
		}
	}
}
