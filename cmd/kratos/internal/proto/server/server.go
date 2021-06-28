package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/emicklei/proto"
	"github.com/spf13/cobra"
)

// CmdServer the service command.
var CmdServer = &cobra.Command{
	Use:   "server",
	Short: "Generate the proto Server implementations",
	Long:  "Generate the proto Server implementations. Example: kratos proto server api/xxx.proto -target-dir=internal/service",
	Run:   run,
}
var targetDir string

func init() {
	CmdServer.Flags().StringVarP(&targetDir, "target-dir", "t", "internal/service", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the proto file. Example: kratos proto server api/xxx.proto")
		return
	}
	reader, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
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
				Service: s.Name,
			}
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if ok {
					cs.Methods = append(cs.Methods, &Method{Service: s.Name, Name: r.Name, Request: r.RequestType, Reply: r.ReturnsType})
				}
			}
			res = append(res, cs)
		}),
	)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exsits\n", targetDir)
		return
	}
	for _, s := range res {
		to := path.Join(targetDir, strings.ToLower(s.Service)+".go")
		if _, err := os.Stat(to); !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s already exists: %s\n", s.Service, to)
			continue
		}
		b, err := s.execute()
		if err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(to, b, 0644); err != nil {
			log.Fatal(err)
		}
		fmt.Println(to)
	}
}
