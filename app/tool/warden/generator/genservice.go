package generator

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	assets "go-common/app/tool/warden/generator/templates"
	"go-common/app/tool/warden/types"
)

const (
	protoTemplateName = "service.tmpl"
	contextType       = "context.Context"
)

// ProtoMessage ProtoMessage
type ProtoMessage struct {
	Name   string
	Fields []ProtoField
}

// ProtoField ProtoField
type ProtoField struct {
	FieldID   int
	FieldType string
	FieldName string
}

// ProtoMethod method info
type ProtoMethod struct {
	Comments []string
	Name     string
	Req      string
	Reply    string
}

// ProtoValue proto template render value
type ProtoValue struct {
	Package   string
	Name      string
	GoPackage string
	Imports   map[string]bool
	Messages  map[string]ProtoMessage
	Methods   []ProtoMethod

	options *ServiceProtoOptions
}

// ServiceProtoOptions ...
type ServiceProtoOptions struct {
	GoPackage    string
	ProtoPackage string
	IgnoreType   bool
	ImportPaths  []string
}

func readProtoPackage(protoFile string) (string, error) {
	fp, err := os.Open(protoFile)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	buf := bufio.NewReader(fp)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "package") {
			continue
		}
		return strings.TrimSpace(strings.TrimRight(line[len("package"):], ";")), nil
	}
	return "", fmt.Errorf("proto %s miss package define", protoFile)
}

func underscore(s string) string {
	cc := []byte(s)
	us := make([]byte, 0, len(cc)+3)
	pervUp := true
	for _, b := range cc {
		if 65 <= b && b <= 90 {
			if pervUp {
				us = append(us, b+32)
			} else {
				us = append(us, '_', b+32)
			}
			pervUp = true
		} else {
			pervUp = false
			us = append(us, b)
		}
	}
	return string(us)
}

func (p *ProtoValue) convertType(t types.Typer) (string, error) {
	switch v := t.(type) {
	case *types.BasicType:
		return convertBasicType(v.String())
	case *types.ArrayType:
		if v.EltType.String() == "byte" {
			return "bytes", nil
		}
		elt, err := p.convertType(v.EltType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("repeated %s", elt), nil
	case *types.MapType:
		kt, err := p.convertType(v.KeyType)
		if err != nil {
			return "", err
		}
		vt, err := p.convertType(v.ValueType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map<%s, %s>", kt, vt), nil
	case *types.StructType:
		if v.ProtoFile == "" {
			messageName := fmt.Sprintf("%s%s", strings.Title(v.Package), v.IdentName)
			err := p.renderMessage(messageName, v.Fields)
			if err != nil {
				return "", err
			}
			return messageName, nil
		}
		protoPackage, err := readProtoPackage(v.ProtoFile)
		if err != nil {
			return "", err
		}
		p.importPackage(v.ProtoFile)
		if p.Package == protoPackage {
			return v.IdentName, nil
		}
		return fmt.Sprintf(".%s.%s", protoPackage, v.IdentName), nil
	}
	return "", fmt.Errorf("unsupport type %s", t)
}

func convertBasicType(gt string) (string, error) {
	switch gt {
	case "float64":
		return "double", nil
	case "float32":
		return "float", nil
	case "int", "int8", "uint8", "int16", "uint16":
		return "int32", nil
	case "int64", "int32", "uint32", "uint64", "string", "bool":
		return gt, nil
	}
	return "", fmt.Errorf("unsupport basic type %s", gt)
}

func (p *ProtoValue) render(spec *types.ServiceSpec, options *ServiceProtoOptions) (*ProtoValue, error) {
	p.options = options
	p.Name = spec.Name
	p.GoPackage = options.GoPackage
	p.Package = options.ProtoPackage
	p.Imports = make(map[string]bool)
	p.Messages = make(map[string]ProtoMessage)
	return p, p.renderMethods(spec.Methods)
}

func (p *ProtoValue) renderMethods(methods []*types.Method) error {
	for _, method := range methods {
		protoMethod := ProtoMethod{
			Comments: method.Comments,
			Name:     method.Name,
		}
		//if len(method.Parameters) == 0 || (len(method.Parameters) == 1 && method.Parameters[0].Type.String() == contextType) {
		//	p.importPackage(emptyProtoFile)
		//	protoMethod.Req = emptyProtoMsg
		//} else {
		//	protoMethod.Req = fmt.Sprintf("%sReq", method.Name)
		//	if err := p.renderMessage(protoMethod.Req, method.Parameters); err != nil {
		//		return err
		//	}
		//}

		//if len(method.Results) == 0 || (len(method.Results) == 1 && method.Results[0].Type.String() == "error") {
		//	p.importPackage(emptyProtoFile)
		//	protoMethod.Reply = emptyProtoMsg
		//} else {
		//	protoMethod.Reply = fmt.Sprintf("%sReply", method.Name)
		//	if err := p.renderMessage(protoMethod.Reply, method.Results); err != nil {
		//		return err
		//	}
		//}
		protoMethod.Req = fmt.Sprintf("%sReq", method.Name)
		if err := p.renderMessage(protoMethod.Req, method.Parameters); err != nil {
			return err
		}
		protoMethod.Reply = fmt.Sprintf("%sReply", method.Name)
		if err := p.renderMessage(protoMethod.Reply, method.Results); err != nil {
			return err
		}
		p.Methods = append(p.Methods, protoMethod)
	}
	return nil
}

func (p *ProtoValue) importPackage(imp string) {
	for _, importPath := range p.options.ImportPaths {
		if strings.HasPrefix(imp, importPath) {
			p.Imports[strings.TrimLeft(imp[len(importPath):], "/")] = true
			return
		}
	}
	p.Imports[imp] = true
}

func (p *ProtoValue) renderMessage(name string, fields []*types.Field) error {
	if _, ok := p.Messages[name]; ok {
		return nil
	}
	message := ProtoMessage{
		Name: name,
	}
	for i, field := range fields {
		if field.Type.String() == "error" || field.Type.String() == contextType {
			continue
		}
		fieldName := underscore(field.Name)
		if fieldName == "" {
			fieldName = fmt.Sprintf("data_%d", i)
		}
		pField := ProtoField{
			FieldID:   i + 1,
			FieldName: fieldName,
		}
		ptype, err := p.convertType(field.Type)
		if err != nil {
			if p.options.IgnoreType {
				log.Printf("warning convert type fail %s", err)
				ptype = fmt.Sprintf("//FIXME type %s", field.Type)
			} else {
				return err
			}
		}
		pField.FieldType = ptype
		message.Fields = append(message.Fields, pField)
	}
	p.Messages[name] = message
	return nil
}

func renderProtoValue(spec *types.ServiceSpec, options *ServiceProtoOptions) (*ProtoValue, error) {
	v := &ProtoValue{}
	return v.render(spec, options)
}

// GenServiceProto generator proto service by service spec
func GenServiceProto(out io.Writer, spec *types.ServiceSpec, options *ServiceProtoOptions) error {
	value, err := renderProtoValue(spec, options)
	if err != nil {
		return err
	}
	assets.MustAsset(protoTemplateName)
	t, err := template.New(protoTemplateName).Parse(string(assets.MustAsset(protoTemplateName)))
	if err != nil {
		return err
	}
	return t.Execute(out, value)
}
