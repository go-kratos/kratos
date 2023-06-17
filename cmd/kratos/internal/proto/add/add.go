package add

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CmdAdd represents the add command.
var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a proto API template",
	Long:  "Add a proto API template. Example: kratos proto add helloworld/v1/hello.proto",
	Run:   run,
}

func run(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the proto file or directory")
		return
	}
	input := args[0]
	n := strings.LastIndex(input, "/")
	if n == -1 {
		fmt.Println("The proto path needs to be hierarchical.")
		return
	}
	path := input[:n]
	fileName := input[n+1:]
	pkgName := strings.ReplaceAll(path, "/", ".")

	p := &Proto{
		Name:        fileName,
		Path:        path,
		Package:     pkgName,
		GoPackage:   goPackage(path),
		JavaPackage: javaPackage(pkgName),
		Service:     serviceName(fileName),
	}
	if err := p.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}

func modName() string {
	modBytes, err := os.ReadFile("go.mod")
	if err != nil {
		if modBytes, err = os.ReadFile("../go.mod"); err != nil {
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}

func goPackage(path string) string {
	s := strings.Split(path, "/")
	return modName() + "/" + path + ";" + s[len(s)-1]
}

func javaPackage(name string) string {
	return name
}

func serviceName(name string) string {
	return toUpperCamelCase(strings.Split(name, ".")[0])
}

func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}
