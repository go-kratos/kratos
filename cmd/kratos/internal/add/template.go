package add

import (
	"bytes"
	"strings"
	"text/template"
)

const protoTemplate = `
syntax = "proto3";

package {{.Package}};

option go_package = "{{.GoPackage}}";
option java_multiple_files = true;
option java_package = "{{.JavaPackage}}";

service {{.Service}} {
    rpc {{.Method}} ({{.Method}}Request) returns ({{.Method}}Reply);
}

message {{.Method}}Request {}

message {{.Method}}Reply {}`

func (p *Proto) execute() (string, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return "", err
	}
	if err := tmpl.Execute(buf, p); err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
