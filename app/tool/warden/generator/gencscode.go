package generator

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	assets "go-common/app/tool/warden/generator/templates"
	"go-common/app/tool/warden/types"
)

// GenCSCodeOptions options
type GenCSCodeOptions struct {
	PbPackage   string
	RecvPackage string
	RecvName    string
}

// CSValue ...
type CSValue struct {
	options       *GenCSCodeOptions
	Name          string
	PbPackage     string
	RecvName      string
	RecvPackage   string
	Imports       map[string]struct{}
	ClientImports map[string]struct{}
	Methods       []CSMethod
}

// CSMethod ...
type CSMethod struct {
	Name         string
	Comments     []string
	ParamBlock   string
	ReturnBlock  string
	ParamPbBlock string
}

func (c *CSValue) render(spec *types.ServiceSpec) error {
	c.PbPackage = c.options.PbPackage
	c.Name = spec.Name
	c.RecvName = c.options.RecvName
	c.RecvPackage = c.options.RecvPackage
	c.Imports = map[string]struct{}{"context": struct{}{}}
	c.ClientImports = make(map[string]struct{})
	return c.renderMethods(spec.Methods)
}

func (c *CSValue) renderMethods(methods []*types.Method) error {
	for _, method := range methods {
		csMethod := CSMethod{
			Name:        method.Name,
			Comments:    method.Comments,
			ParamBlock:  c.formatField(method.Parameters),
			ReturnBlock: c.formatField(method.Results),
		}
		c.Methods = append(c.Methods, csMethod)
	}
	return nil
}

func (c *CSValue) formatField(fields []*types.Field) string {
	var ss []string
	clientImps := make(map[string]struct{})
	for _, field := range fields {
		if field.Name == "" {
			ss = append(ss, field.Type.String())
		} else {
			ss = append(ss, fmt.Sprintf("%s %s", field.Name, field.Type))
		}
		importType(clientImps, field.Type)
	}
	for k := range clientImps {
		if _, ok := c.Imports[k]; !ok {
			c.ClientImports[k] = struct{}{}
		}
	}
	return strings.Join(ss, ", ")
}

func importType(m map[string]struct{}, t types.Typer) {
	if m == nil {
		panic("map is nil")
	}
	switch v := t.(type) {
	case *types.StructType:
		m[v.ImportPath] = struct{}{}
		for _, f := range v.Fields {
			importType(m, f.Type)
		}
	case *types.ArrayType:
		importType(m, v.EltType)
	case *types.InterfaceType:
		m[v.ImportPath] = struct{}{}
	}
}

func renderCSValue(spec *types.ServiceSpec, options *GenCSCodeOptions) (*CSValue, error) {
	value := &CSValue{
		options: options,
	}
	return value, value.render(spec)
}

// GenCSCode generator client, server code
func GenCSCode(csdir string, spec *types.ServiceSpec, options *GenCSCodeOptions) error {
	value, err := renderCSValue(spec, options)
	if err != nil {
		return err
	}
	return genCode(value, "server", csdir)
}

func genCode(value *CSValue, name, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	fp, err := os.OpenFile(path.Join(dir, fmt.Sprintf("%s.go", name)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	templateName := fmt.Sprintf("%s.tmpl", name)
	t, err := template.New(name).Parse(string(assets.MustAsset(templateName)))
	if err != nil {
		return err
	}
	return t.Execute(fp, value)
}
