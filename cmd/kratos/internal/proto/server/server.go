package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CmdServer the service command.
var CmdServer = &cobra.Command{
	Use:   "server",
	Short: "Generate the proto server implementations",
	Long:  "Generate the proto server implementations. Example: kratos proto server api/xxx.proto --target-dir=internal/service",
	Run:   run,
}

var (
	targetDir    string
	templatePath string
)

func init() {
	CmdServer.Flags().StringVarP(&targetDir, "target-dir", "t", "internal/service", "generate target directory")
	CmdServer.Flags().StringVarP(&templatePath, "template-file", "m", "", "specify custom template file")
}

func run(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the proto file. Example: kratos proto server api/xxx.proto")
		return
	}
	reader, err := os.Open(args[0])
	if err != nil {
		log.Fatalf("Failed to open proto file %q: %v", args[0], err)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Fatalf("Failed to parse proto file %q: %v", args[0], err)
	}

	var (
		pkg string
		res []*Service
	)
	proto.Walk(definition,
		proto.WithOption(func(o *proto.Option) {
			if o.Name == "go_package" {
				pkg = strings.Split(o.Constant.Source, ";")[0]
			}
		}),
		proto.WithService(func(s *proto.Service) {
			cs := &Service{
				Package: pkg,
				Service: serviceName(s.Name),
			}
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if !ok {
					continue
				}
				cs.Methods = append(cs.Methods, &Method{
					Service: serviceName(s.Name), Name: rpcName(r.Name), Request: parametersName(r.RequestType),
					Reply: parametersName(r.ReturnsType), Type: getMethodType(r.StreamsRequest, r.StreamsReturns),
				})
			}
			res = append(res, cs)
		}),
	)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exist\n", targetDir)
		return
	}
	for _, s := range res {
		to := filepath.Join(targetDir, strings.ToLower(s.Service)+".go")
		if _, err := os.Stat(to); !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s already exists: %s\n", s.Service, to)
			continue
		}
		b, err := s.execute(templatePath)
		if err != nil {
			log.Fatalf("Failed to generate server implementation for %s: %v", s.Service, err)
		}
		if err := os.WriteFile(to, b, 0o644); err != nil {
			log.Fatalf("Failed to write server file %q: %v", to, err)
		}
		fmt.Println(to)
	}
}

func getMethodType(streamsRequest, streamsReturns bool) MethodType {
	if !streamsRequest && !streamsReturns {
		return unaryType
	} else if streamsRequest && streamsReturns {
		return twoWayStreamsType
	} else if streamsRequest {
		return requestStreamsType
	} else if streamsReturns {
		return returnsStreamsType
	}
	return unaryType
}

func parametersName(name string) string {
	return strings.ReplaceAll(name, ".", "_")
}

func serviceName(name string) string {
	return strings.TrimSuffix(toUpperCamelCase(strings.Split(name, ".")[0]), "Service")
}

func rpcName(name string) string {
	return toUpperCamelCase(strings.Split(name, ".")[0])
}

func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}
