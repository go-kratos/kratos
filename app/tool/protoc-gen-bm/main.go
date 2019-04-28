package main

import (
	"flag"
	"os"

	"github.com/golang/glog"

	"go-common/app/tool/protoc-gen-bm/codegenerator"
	"go-common/app/tool/protoc-gen-bm/genbm"
	"go-common/app/tool/protoc-gen-bm/util"
)

var useJSONPB bool

func init() {
	flag.BoolVar(&useJSONPB, "jsonpb", false, "use jsonpb instead of std library, NOTE: jsonpb very slow")
}

func main() {
	flag.Parse()
	req, err := codegenerator.ParseRequest(os.Stdin)
	if err != nil {
		glog.Fatal(err)
	}
	if err = util.ParseParamSetFlag(req.GetParameter(), flag.CommandLine); err != nil {
		glog.Fatal(err)
	}
	g := genbm.New(useJSONPB)
	resp, err := g.Generate(req)
	if err = codegenerator.WriteResponse(os.Stdout, resp, err); err != nil {
		glog.Fatal(err)
	}
}
