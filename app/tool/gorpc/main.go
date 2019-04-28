// A commandline tool for generating rpc code by service methods.
//
// This tool can generate rpc client code ,rpc arg model for specific Go project dir.
// Usage :
//  $gorpc [options]
// Available options:
//
//  -client   generate rpc client code.
//
//  -d        specific go project service dir (default ./service/)
//
//  -model    generate rpc arg model code.
//
//  -o        specific rpc client code output file. (default ./rpc/client/client.go)
//
//  -m        specific rpc arg model output file. (default ./model/rpc.go)
//
//  -s        print code to stdout.
// Example:
//  $ cd $GOPATH/relation
//  $ gorpc
//  such command will generate rpc client code by functions of ./service/* and write code to $GOPATH/relation/rpc/client/client.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/importer"
	"io"
	"io/ioutil"
	"log"
	"path"

	"go-common/app/tool/gorpc/goparser"
	"go-common/app/tool/gorpc/input"
	"go-common/app/tool/gorpc/model"

	"golang.org/x/tools/imports"
)

const (
	tpl = `var (
	_noRes = &struct{}{}
)
type Service struct {
	client *rpc.Client2
}
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(c)
	return
}`
)

type output struct {
	rpcClient *bytes.Buffer
	methodStr *bytes.Buffer
	methods   *bytes.Buffer
	model     *bytes.Buffer
}

var out = &output{
	rpcClient: new(bytes.Buffer),
	methodStr: new(bytes.Buffer),
	methods:   new(bytes.Buffer),
	model:     new(bytes.Buffer),
}
var (
	dir        = flag.String("d", "./service/", "source code dir")
	argModel   = flag.Bool("model", false, "use -model to generate rpc arg model")
	client     = flag.Bool("client", true, "use -client to generate rpc client code")
	clientfile = flag.String("o", "./rpc/client/client.go", "out file name")
	modelfile  = flag.String("m", "./model/rpc.go", "generate rpc arg model")
	std        = flag.Bool("s", false, "use -s to print code to stdout")
)

func main() {
	flag.Parse()
	p := &goparser.Parser{Importer: importer.Default()}
	files, err := input.Files(path.Dir(*dir))
	if err != nil {
		panic(err)
	}
	if *argModel {
		out.model.WriteString("package model\n")
	}
	for _, f := range files {
		rs, err := p.Parse(string(f), files)
		if err != nil {
			panic(err)
		}
		for _, f := range rs.Funcs {
			if f.IsExported && len(f.Results) <= 1 && f.Receiver != nil && f.ReturnsError {
				if *client {
					out.generateMeStr(f)
					out.generateMethod(f)
				}
				if *argModel {
					out.generateModel(f)
				}
			}
		}
	}

	if *argModel {
		m, err := imports.Process("model.go", out.model.Bytes(), &imports.Options{TabWidth: 4})
		if err != nil {
			log.Printf("gopimorts err %v", err)
			return
		}
		if *std {
			fmt.Printf("%s", m)
		} else {
			ioutil.WriteFile(*modelfile, m, 0666)
		}
	}
	if *client {
		out.rpcClient.WriteString("package client \n")
		out.rpcClient.WriteString(tpl)
		out.rpcClient.WriteString("\nconst ( \n")
		io.Copy(out.rpcClient, out.methodStr)
		out.rpcClient.WriteString("\n)\n")
		io.Copy(out.rpcClient, out.methods)
		c, err := imports.Process("client.go", out.rpcClient.Bytes(), &imports.Options{TabWidth: 4})
		if err != nil {
			log.Printf("gopimorts err %v", err)
			return
		}
		if *std {
			fmt.Printf("%s", c)
		} else {
			ioutil.WriteFile(*clientfile, c, 0666)
		}
	}
}

func (o *output) generateMeStr(f *model.Function) {
	b := o.methodStr
	fmt.Fprintf(b, "_%s = \"RPC.%s\"\n", f.Name, f.Name)
}

func (o *output) generateModel(f *model.Function) {
	b := o.model
	for _, p := range f.Parameters {
		if !p.IsBasicType() {
			if p.Type.String() == "context.Context" {
				continue
			}
			return
		}
	}
	fmt.Fprintf(b, fmt.Sprintf("type Arg%s struct{\n", f.Name))
	for _, p := range f.Parameters {
		fmt.Fprintln(b, string(bytes.ToUpper([]byte(p.Name))), p.Type)
	}
	fmt.Fprintln(b, "}")
}

func (o *output) generateMethod(f *model.Function) {
	b := o.methods
	name := f.Name
	for _, p := range f.Parameters {
		if !p.IsBasicType() {
			if p.Type.String() == "context.Context" {
				continue
			}
			return
		}
	}
	fmt.Fprintf(b, "func(s %s)%s(", f.Receiver.Type, f.Name)
	fmt.Fprintf(b, "c context.Context, arg *model.Arg%s)(", f.Name)
	if f.ReturnsError {
		if len(f.Results) > 0 {
			fmt.Fprintf(b, "res %s", f.Results[0].Type)
			fmt.Fprint(b, ",")
		}
		fmt.Fprintf(b, "err error)")
	} else {
		fmt.Fprint(b, ")")
	}
	fmt.Fprint(b, " {\n")
	if len(f.Results) == 0 {
		fmt.Fprintf(b, "err = s.client.Call(c, _%s, arg, &_noRes)", name)
	} else if f.Results[0].Type.IsStar {
		fmt.Fprintf(b, "res = new(%s)\n", f.Results[0].Type.String()[1:])
		fmt.Fprintf(b, "err = s.client.Call(c, _%s, arg, res)", name)
	} else {
		fmt.Fprintf(b, "err = s.client.Call(c, _%s, arg, &res)", name)
	}
	fmt.Fprintf(b, "\nreturn\n}\n")
}
