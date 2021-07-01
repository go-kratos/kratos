package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/base"

	"github.com/spf13/cobra"
)

var (
	// CmdClient represents the source command.
	CmdClient = &cobra.Command{
		Use:                "client",
		Short:              "Generate the proto client code",
		Long:               "Generate the proto client code. Example: kratos proto client helloworld.proto",
		DisableFlagParsing: true,
		Run:                run,
	}
)

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the proto file or directory")
		return
	}
	var (
		err   error
		proto = strings.TrimSpace(args[0])
	)
	if err = look("protoc-gen-go", "protoc-gen-go-grpc", "protoc-gen-go-http", "protoc-gen-go-errors", "protoc-gen-validate"); err != nil {
		// update the kratos plugins
		cmd := exec.Command("kratos", "upgrade")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	}
	if strings.HasSuffix(proto, ".proto") {
		err = generate(proto, args)
	} else {
		err = walk(proto, args)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}

func walk(dir string, args []string) error {
	if dir == "" {
		dir = "."
	}
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" || strings.HasPrefix(path, "third_party") {
			return nil
		}
		return generate(path, args)
	})
}

// generate is used to execute the generate command for the specified proto file
func generate(proto string, args []string) error {
	input := []string{
		"--proto_path=.",
		"--proto_path=./third_party",
		"--proto_path=" + base.KratosMod(),
		"--proto_path=" + filepath.Join(base.KratosMod(), "third_party"),
		"--go_out=paths=source_relative:.",
		"--go-grpc_out=paths=source_relative:.",
		"--go-http_out=paths=source_relative:.",
		"--go-errors_out=paths=source_relative:.",
	}
	protoBytes, err := ioutil.ReadFile(proto)
	if err == nil && len(protoBytes) > 0 {
		if ok, _ := regexp.Match(`\n[^/]*(import)\s+"validate/validate.proto"`, protoBytes); ok {
			input = append(input, "--validate_out=lang=go,paths=source_relative:.")
		}
	}
	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}
	fmt.Printf("proto: %s\n", proto)
	return nil
}
