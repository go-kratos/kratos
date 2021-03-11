package server

import (
	"bytes"
	"html/template"
)

var serviceTemplate = `
package service

import(
	"context"

	pb "{{.Package}}"
)

type {{.Service}}Service struct {
	pb.Unimplemented{{.Service}}Server
}

func New{{.Service}}Service() pb.{{.Service}}Server {
	return &{{.Service}}Service{}
}
{{ range .Methods }}
func (s *{{.Service}}Service) {{.Name}}(ctx context.Context, req *pb.{{.Request}}) (*pb.{{.Reply}}, error) {
	return &pb.{{.Reply}}{}, nil
}
{{- end }}
`

// Service is a proto service.
type Service struct {
	Package string
	Service string
	Methods []*Method
}

// Method is a proto method.
type Method struct {
	Service string
	Name    string
	Request string
	Reply   string
}

func (s *Service) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
