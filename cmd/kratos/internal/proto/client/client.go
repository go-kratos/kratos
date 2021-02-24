package client

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/base"
	"github.com/spf13/cobra"
)

var (
	// CmdClient represents the source command.
	CmdClient = &cobra.Command{
		Use:   "client",
		Short: "Generate the proto client code",
		Long:  "Generate the proto client code. Example: kratos proto client helloworld.proto",
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	var (
		err   error
		proto = strings.TrimSpace(args[0])
	)
	if _, err = exec.LookPath("protoc-gen-go-http"); err != nil {
		// update the kratos plugins
		if err := exec.Command("kratos", "upgrade").Run(); err != nil {
			fmt.Println(err)
		}
	}
	if strings.HasSuffix(proto, ".proto") {
		err = generate(proto)
	} else {
		err = walk(proto)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func walk(dir string) error {
	if dir == "" {
		dir = "."
	}
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" {
			return nil
		}
		return generate(path)
	})
}

// generate is used to execute the generate command for the specified proto file
func generate(proto string) error {
	path, name := filepath.Split(proto)
	fd := exec.Command("protoc", []string{
		"--proto_path=.",
		"--proto_path=" + filepath.Join(base.KratosMod(), "api"),
		"--proto_path=" + filepath.Join(base.KratosMod(), "third_party"),
		"--proto_path=" + filepath.Join(os.Getenv("GOPATH"), "src"),
		"--go_out=paths=source_relative:.",
		"--go-grpc_out=paths=source_relative:.",
		"--go-http_out=paths=source_relative:.",
		"--go-errors_out=paths=source_relative:.",
		name,
	}...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = path
	if err := fd.Run(); err != nil {
		return err
	}
	fmt.Printf("proto: %s\n", proto)
	return nil
}
