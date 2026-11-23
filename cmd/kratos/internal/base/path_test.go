package base

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFileSkipsReplacementInGeneratedProtobufFiles(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()
	src := filepath.Join(srcDir, "conf.pb.go")
	dst := filepath.Join(dstDir, "conf.pb.go")
	content := []byte("Z5github.com/go-kratos/kratos-layout/internal/conf;conf")
	if err := os.WriteFile(src, content, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(src, dst, []string{"github.com/go-kratos/kratos-layout", "helloworld2"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(got, []byte("helloworld2")) {
		t.Fatalf("generated protobuf file was rewritten, got %q", got)
	}
	if !bytes.Contains(got, content) {
		t.Fatalf("generated protobuf file content changed, got %q", got)
	}
}

func TestCopyFileReplacesModulePathInGoSourceFiles(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()
	src := filepath.Join(srcDir, "main.go")
	dst := filepath.Join(dstDir, "main.go")
	content := []byte("module github.com/go-kratos/kratos-layout")
	if err := os.WriteFile(src, content, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(src, dst, []string{"github.com/go-kratos/kratos-layout", "helloworld2"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(got, []byte("module helloworld2")) {
		t.Fatalf("go source file was not rewritten, got %q", got)
	}
}
