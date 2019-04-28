package goparser

import (
	"os"
	"testing"
)

var dpath = "/Users/weicheng/Go/src/go-common/app/service/account/service"

//var dpath = "/Users/weicheng/Go/src/playground/testgen/service"

func TestParse(t *testing.T) {
	spec, err := Parse("account", dpath, "Service", dpath)
	if err != nil {
		t.Fatal(err)
	}
	for _, method := range spec.Methods {
		t.Logf("method %s", method.Name)
		for _, param := range method.Parameters {
			t.Logf(">> param %s", param)
		}
		for _, result := range method.Results {
			t.Logf("<< result %s", result)
		}
	}
}

func TestExtractProtoFile(t *testing.T) {
	comment := "// source: article.proto\n"
	protoFile := extractProtoFile(comment)
	if protoFile != "article.proto" {
		t.Errorf("expect %s get %s", "article.proto", protoFile)
	}
}

func TestGoPackage(t *testing.T) {
	os.Setenv("GOPATH", "/go:/go1:/go3")
	type args struct {
		dpath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{"/go/src/hello/hello.go"},
			want: "hello",
		},
		{
			name: "test2",
			args: args{"/go3/src/hello/foo/hello.go"},
			want: "hello/foo",
		},
		{
			name:    "test3",
			args:    args{"/g/src/hello/foo/hello.go"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoPackage(tt.args.dpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GoPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}
