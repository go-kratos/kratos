package opensergo

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2"

	v1 "github.com/opensergo/opensergo-go/proto/service_contract/v1"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Option func(*options)

func WithEndpoint(endpoint string) Option {
	return func(o *options) {
		o.Endpoint = endpoint
	}
}

type options struct {
	Endpoint string `json:"endpoint"`
}

func (o *options) ParseJSON(data []byte) error {
	return json.Unmarshal(data, o)
}

type OpenSergo struct {
	mdClient v1.MetadataServiceClient
}

func New(opts ...Option) (*OpenSergo, error) {
	opt := options{
		Endpoint: os.Getenv("OPENSERGO_ENDPOINT"),
	}
	// https://github.com/opensergo/opensergo-specification/blob/main/specification/en/README.md
	if v := os.Getenv("OPENSERGO_BOOTSTRAP"); v != "" {
		if err := opt.ParseJSON([]byte(v)); err != nil {
			return nil, err
		}
	}
	if v := os.Getenv("OPENSERGO_BOOTSTRAP_CONFIG"); v != "" {
		b, err := ioutil.ReadFile(v)
		if err != nil {
			return nil, err
		}
		if err := opt.ParseJSON(b); err != nil {
			return nil, err
		}
	}
	for _, o := range opts {
		o(&opt)
	}
	dialCtx := context.Background()
	dialCtx, cancel := context.WithTimeout(dialCtx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, opt.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &OpenSergo{
		mdClient: v1.NewMetadataServiceClient(conn),
	}, nil
}

func (s *OpenSergo) ReportMetadata(ctx context.Context, app kratos.AppInfo) error {
	services, types, err := listDescriptors()
	if err != nil {
		return err
	}

	serviceMetadata := &v1.ServiceMetadata{
		ServiceContract: &v1.ServiceContract{
			Services: services,
			Types:    types,
		},
	}

	for _, endpoint := range app.Endpoint() {
		u, err := url.Parse(endpoint) //nolint
		if err != nil {
			return err
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}
		portValue, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		serviceMetadata.Protocols = append(serviceMetadata.Protocols, u.Scheme)
		serviceMetadata.ListeningAddresses = append(serviceMetadata.ListeningAddresses, &v1.SocketAddress{
			Address:   host,
			PortValue: uint32(portValue),
		})
	}
	_, err = s.mdClient.ReportMetadata(ctx, &v1.ReportMetadataRequest{
		AppName:         app.Name(),
		ServiceMetadata: []*v1.ServiceMetadata{serviceMetadata},
		// TODO: Node: *v1.Node,
	})
	return err
}

func listDescriptors() (services []*v1.ServiceDescriptor, types []*v1.TypeDescriptor, err error) {
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Services().Len(); i++ {
			var (
				methods []*v1.MethodDescriptor
				sd      = fd.Services().Get(i)
			)
			for j := 0; j < sd.Methods().Len(); j++ {
				md := sd.Methods().Get(j)
				mName := string(md.Name())
				inputType := string(md.Input().FullName())
				outputType := string(md.Output().FullName())
				isClientStreaming := md.IsStreamingClient()
				isServerStreaming := md.IsStreamingServer()
				pattern := proto.GetExtension(md.Options(), annotations.E_Http).(*annotations.HttpRule).GetPattern()
				var httpPath, httpMethod string
				if pattern != nil {
					httpMethod, httpPath = HTTPPatternInfo(pattern)
				}
				methodDesc := v1.MethodDescriptor{
					Name:            mName,
					InputTypes:      []string{inputType},
					OutputTypes:     []string{outputType},
					ClientStreaming: &isClientStreaming,
					ServerStreaming: &isServerStreaming,
					HttpPaths:       []string{httpPath},
					HttpMethods:     []string{httpMethod},
					// TODO: Description: *string,
				}
				methods = append(methods, &methodDesc)
			}
			services = append(services, &v1.ServiceDescriptor{
				Name:    string(sd.Name()),
				Methods: methods,
				// TODO: Description: *string,
			})
		}

		for i := 0; i < fd.Messages().Len(); i++ {
			var (
				fields []*v1.FieldDescriptor
				md     = fd.Messages().Get(i)
			)

			for j := 0; j < md.Fields().Len(); j++ {
				fd := md.Fields().Get(j)
				kind := fd.Kind()
				typeName := kind.String()

				fields = append(fields, &v1.FieldDescriptor{
					Name:     string(fd.Name()),
					Number:   int32(fd.Number()),
					Type:     v1.FieldDescriptor_Type(kind),
					TypeName: &typeName,
					// TODO: Description: *string,
				})
			}

			types = append(types, &v1.TypeDescriptor{
				Name:   string(md.Name()),
				Fields: fields,
			})
		}

		return true
	})
	return
}

func HTTPPatternInfo(pattern interface{}) (method string, path string) {
	switch p := pattern.(type) {
	case *annotations.HttpRule_Get:
		return "GET", p.Get
	case *annotations.HttpRule_Post:
		return "POST", p.Post
	case *annotations.HttpRule_Delete:
		return "DELETE", p.Delete
	case *annotations.HttpRule_Patch:
		return "PATCH", p.Patch
	case *annotations.HttpRule_Put:
		return "PUT", p.Put
	case *annotations.HttpRule_Custom:
		return p.Custom.Kind, p.Custom.Path
	default:
		return "", ""
	}
}
