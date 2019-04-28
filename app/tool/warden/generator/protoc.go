package generator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func findVendorDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("getwd error: %s", err)
	}
	for dir := pwd; dir != "/" && dir != "."; dir = path.Dir(dir) {
		vendorDir := path.Join(dir, "vendor")
		if s, err := os.Stat(vendorDir); err == nil && s.IsDir() {
			return vendorDir
		}
	}
	return ""
}

// Protoc run protoc generator go source code
func Protoc(protoFile, protocExec, gen string, paths []string) error {
	if protocExec == "" {
		protocExec = "protoc"
	}
	if gen == "" {
		gen = "gogofast"
	}
	paths = append(paths, ".", os.Getenv("GOPATH"))
	vendorDir := findVendorDir()
	if vendorDir != "" {
		paths = append(paths, vendorDir)
	}
	args := []string{"--proto_path", strings.Join(paths, ":"), fmt.Sprintf("--%s_out=plugins=grpc:.", gen), path.Base(protoFile)}
	log.Printf("run protoc %s", strings.Join(args, " "))
	protoc := exec.Command(protocExec, args...)
	protoc.Stdout = os.Stdout
	protoc.Stderr = os.Stderr
	protoc.Dir = path.Dir(protoFile)
	return protoc.Run()
}
