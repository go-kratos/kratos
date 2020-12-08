package service

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

// CmdService represents the service command.
var CmdService = &cobra.Command{
	Use:   "service",
	Short: "Generate the proto Service implementations",
	Long:  "Generate the proto Service implementations. Example: kratos proto service api/xxx.proto -target-dir=internal/service",
	Run:   run,
}
var targetDir string

func init() {
	CmdService.Flags().StringVarP(&targetDir, "-target-dir", "t", ".", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
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
	for _, s := range res {
		to := path.Join(targetDir, strings.ToLower(s.Service)+".go")
		if _, err := os.Stat(to); !os.IsNotExist(err) {
			fmt.Printf("%s already exists\n", s.Service)
			continue
		}
		b, err := s.execute()
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(to, b, 0644)
	}
}
