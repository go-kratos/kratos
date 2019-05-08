package generator

import (
	"os"
	"os/exec"
	"testing"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func TestGenerateParseCommandLineParamsError(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		g := &bm{}
		g.Generate(&plugin.CodeGeneratorRequest{
			Parameter: proto.String("invalid"),
		})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGenerateParseCommandLineParamsError")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
