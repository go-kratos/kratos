package generator

import (
	"io"

	"go-common/app/tool/gengo/namer"
	"go-common/app/tool/gengo/types"
)

// consts
const (
	GolangFileType = "golang"
)

// DefaultGen implements a do-nothing Generator.
//
// It can be used to implement static content files.
type DefaultGen struct {
	// OptionalName, if present, will be used for the generator's name, and
	// the filename (with ".go" appended).
	OptionalName string

	// OptionalBody, if present, will be used as the return from the "Init"
	// method. This causes it to be static content for the entire file if
	// no other generator touches the file.
	OptionalBody []byte
}

func (d DefaultGen) Name() string                                        { return d.OptionalName }
func (d DefaultGen) Filter(*Context, *types.Type) bool                   { return true }
func (d DefaultGen) Namers(*Context) namer.NameSystems                   { return nil }
func (d DefaultGen) Imports(*Context) []string                           { return []string{} }
func (d DefaultGen) PackageVars(*Context) []string                       { return []string{} }
func (d DefaultGen) PackageConsts(*Context) []string                     { return []string{} }
func (d DefaultGen) GenerateType(*Context, *types.Type, io.Writer) error { return nil }
func (d DefaultGen) Filename() string                                    { return d.OptionalName + ".go" }
func (d DefaultGen) FileType() string                                    { return GolangFileType }
func (d DefaultGen) Finalize(*Context, io.Writer) error                  { return nil }

func (d DefaultGen) Init(c *Context, w io.Writer) error {
	_, err := w.Write(d.OptionalBody)
	return err
}

var (
	_ = Generator(DefaultGen{})
)
