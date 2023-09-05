package project

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/cmd/kratos/v2/internal/base"
)

// TestCmdNew tests the `kratos new` command.
func TestCmdNew(t *testing.T) {
	cwd := changeCurrentDir(t)
	projectName := "helloworld"

	// create a new project
	CmdNew.SetArgs([]string{projectName})
	if err := CmdNew.Execute(); err != nil {
		t.Fatalf("executing command: %v", err)
	}

	// check that the expected files were created
	for _, file := range []string{
		"go.mod",
		"go.sum",
		"README.md",
		"cmd/helloworld/main.go",
	} {
		if _, err := os.Stat(filepath.Join(cwd, projectName, file)); err != nil {
			t.Errorf("expected file %s to exist", file)
		}
	}

	// check that the go.mod file contains the expected module name
	assertGoMod(t, filepath.Join(cwd, projectName, "go.mod"), projectName)

	assertImportsInclude(t, filepath.Join(cwd, projectName, "cmd", projectName, "wire.go"), fmt.Sprintf(`"%s/internal/biz"`, projectName))
}

// TestCmdNewNoMod tests the `kratos new` command with the --nomod flag.
func TestCmdNewNoMod(t *testing.T) {
	cwd := changeCurrentDir(t)

	// create a new project
	CmdNew.SetArgs([]string{"project"})
	if err := CmdNew.Execute(); err != nil {
		t.Fatalf("executing command: %v", err)
	}

	// add new app with --nomod flag
	CmdNew.SetArgs([]string{"--nomod", "project/app/user"})
	if err := CmdNew.Execute(); err != nil {
		t.Fatalf("executing command: %v", err)
	}

	// check that the expected files were created
	for _, file := range []string{
		"go.mod",
		"go.sum",
		"README.md",
		"cmd/project/main.go",
		"app/user/cmd/user/main.go",
	} {
		if _, err := os.Stat(filepath.Join(cwd, "project", file)); err != nil {
			t.Errorf("expected file %s to exist", file)
		}
	}

	assertImportsInclude(t, filepath.Join(cwd, "project/app/user/cmd/user/wire.go"), `"project/app/user/internal/biz"`)
}

// assertImportsInclude checks that the file at path contains the expected import.
func assertImportsInclude(t *testing.T, path, expected string) {
	t.Helper()

	got, err := imports(path)
	if err != nil {
		t.Fatalf("getting imports: %v", err)
	}

	for _, imp := range got {
		if imp == expected {
			return
		}
	}

	t.Errorf("expected imports to include %s, got %v", expected, got)
}

// imports returns the imports in the file at path.
func imports(path string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	imports := make([]string, 0, len(f.Imports))
	for _, s := range f.Imports {
		imports = append(imports, s.Path.Value)
	}

	return imports, nil
}

// assertGoMod checks that the go.mod file contains the expected module name.
func assertGoMod(t *testing.T, path, expected string) {
	t.Helper()

	got, err := base.ModulePath(path)
	if err != nil {
		t.Fatalf("getting module path: %v", err)
	}

	if got != expected {
		t.Errorf("expected module name %s, got %s", expected, got)
	}
}

// change the working directory to the tempdir
func changeCurrentDir(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()

	oldCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting working directory: %v", err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("changing working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldCWD); err != nil {
			t.Fatalf("restoring working directory: %v", err)
		}
	})

	return tmp
}
