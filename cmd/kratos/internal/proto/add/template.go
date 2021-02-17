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
    rpc Create{{.Service}} (Create{{.Service}}Request) returns (Create{{.Service}}Reply);
    rpc Update{{.Service}} (Update{{.Service}}Request) returns (Update{{.Service}}Reply);
    rpc Delete{{.Service}} (Delete{{.Service}}Request) returns (Delete{{.Service}}Reply);
    rpc Get{{.Service}} (Get{{.Service}}Request) returns (Get{{.Service}}Reply);
    rpc List{{.Service}} (List{{.Service}}Request) returns (List{{.Service}}Reply);
}

message Create{{.Service}}Request {}
message Create{{.Service}}Reply {}

message Update{{.Service}}Request {}
message Update{{.Service}}Reply {}

message Delete{{.Service}}Request {}
message Delete{{.Service}}Reply {}

message Get{{.Service}}Request {}
message Get{{.Service}}Reply {}

message List{{.Service}}Request {}
message List{{.Service}}Reply {}
`

func (p *Proto) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
