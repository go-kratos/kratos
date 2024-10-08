package opensergo

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	srvContractPb "github.com/opensergo/opensergo-go/proto/service_contract/v1"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type testMetadataServiceServer struct {
	srvContractPb.UnimplementedMetadataServiceServer
}

func (m *testMetadataServiceServer) ReportMetadata(_ context.Context, _ *srvContractPb.ReportMetadataRequest) (*srvContractPb.ReportMetadataReply, error) {
	return &srvContractPb.ReportMetadataReply{}, nil
}

type testAppInfo struct {
	id       string
	name     string
	version  string
	metaData map[string]string
	endpoint []string
}

func (t testAppInfo) ID() string {
	return t.id
}

func (t testAppInfo) Name() string {
	return t.name
}

func (t testAppInfo) Version() string {
	return t.version
}

func (t testAppInfo) Metadata() map[string]string {
	return t.metaData
}

func (t testAppInfo) Endpoint() []string {
	return t.endpoint
}

func TestWithEndpoint(t *testing.T) {
	o := &options{}
	v := "127.0.0.1:9090"
	WithEndpoint(v)(o)
	if !reflect.DeepEqual(v, o.Endpoint) {
		t.Fatalf("o.Endpoint:%s is not equal to v:%s", o.Endpoint, v)
	}
}

func TestOptionsParseJSON(t *testing.T) {
	want := &options{
		Endpoint: "127.0.0.1:9090",
	}
	o := &options{}
	if err := o.ParseJSON([]byte(`{"endpoint":"127.0.0.1:9090"}`)); err != nil {
		t.Fatalf("o.ParseJSON(v) error:%s", err)
	}
	if !reflect.DeepEqual(o, want) {
		t.Fatalf("o:%v is not equal to want:%v", o, want)
	}
}

func TestListDescriptors(t *testing.T) {
	testPb := &descriptorpb.FileDescriptorProto{
		Syntax:  proto.String("proto3"),
		Name:    proto.String("test.proto"),
		Package: proto.String("test"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("TestMessage"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("id"),
						JsonName: proto.String("id"),
						Number:   proto.Int32(1),
						Type:     descriptorpb.FieldDescriptorProto_Type(pref.Int32Kind).Enum(),
					},
					{
						Name:     proto.String("name"),
						JsonName: proto.String("name"),
						Number:   proto.Int32(2),
						Type:     descriptorpb.FieldDescriptorProto_Type(pref.StringKind).Enum(),
					},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("TestService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       proto.String("Create"),
						InputType:  proto.String("TestMessage"),
						OutputType: proto.String("TestMessage"),
					},
				},
			},
		},
	}

	fd, err := protodesc.NewFile(testPb, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile(pb, nil) error:%s", err)
	}

	protoregistry.GlobalFiles = new(protoregistry.Files)
	err = protoregistry.GlobalFiles.RegisterFile(fd)
	if err != nil {
		t.Fatalf("protoregistry.GlobalFiles.RegisterFile(fd) error:%s", err)
	}

	want := struct {
		services []*srvContractPb.ServiceDescriptor
		types    []*srvContractPb.TypeDescriptor
	}{
		services: []*srvContractPb.ServiceDescriptor{
			{
				Name: "TestService",
				Methods: []*srvContractPb.MethodDescriptor{
					{
						Name:            "Create",
						InputTypes:      []string{"test.TestMessage"},
						OutputTypes:     []string{"test.TestMessage"},
						ClientStreaming: proto.Bool(false),
						ServerStreaming: proto.Bool(false),
						Description:     nil,
						HttpPaths:       []string{""},
						HttpMethods:     []string{""},
					},
				},
			},
		},
		types: []*srvContractPb.TypeDescriptor{
			{
				Name: "TestMessage",
				Fields: []*srvContractPb.FieldDescriptor{
					{
						Name:     "id",
						Number:   int32(1),
						Type:     srvContractPb.FieldDescriptor_TYPE_INT32,
						TypeName: proto.String("int32"),
					},
					{
						Name:     "name",
						Number:   int32(2),
						Type:     srvContractPb.FieldDescriptor_TYPE_STRING,
						TypeName: proto.String("string"),
					},
				},
			},
		},
	}

	services, types, err := listDescriptors()
	if err != nil {
		t.Fatalf("listDescriptors error:%s", err)
	}

	if !reflect.DeepEqual(services, want.services) {
		t.Fatalf("services:%v is not equal to want.services:%v", services, want.services)
	}
	if !reflect.DeepEqual(types, want.types) {
		t.Fatalf("types:%v is not equal to want.types:%v", types, want.types)
	}
}

func TestHTTPPatternInfo(t *testing.T) {
	type args struct {
		pattern interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantMethod string
		wantPath   string
	}{
		{
			name: "get",
			args: args{
				pattern: &annotations.HttpRule_Get{Get: "/foo"},
			},
			wantMethod: http.MethodGet,
			wantPath:   "/foo",
		},
		{
			name: "post",
			args: args{
				pattern: &annotations.HttpRule_Post{Post: "/foo"},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/foo",
		},
		{
			name: "put",
			args: args{
				pattern: &annotations.HttpRule_Put{Put: "/foo"},
			},
			wantMethod: http.MethodPut,
			wantPath:   "/foo",
		},
		{
			name: "delete",
			args: args{
				pattern: &annotations.HttpRule_Delete{Delete: "/foo"},
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/foo",
		},
		{
			name: "patch",
			args: args{
				pattern: &annotations.HttpRule_Patch{Patch: "/foo"},
			},
			wantMethod: http.MethodPatch,
			wantPath:   "/foo",
		},
		{
			name: "custom",
			args: args{
				pattern: &annotations.HttpRule_Custom{
					Custom: &annotations.CustomHttpPattern{
						Kind: "CUSTOM",
						Path: "/foo",
					},
				},
			},
			wantMethod: "CUSTOM",
			wantPath:   "/foo",
		},
		{
			name: "other",
			args: args{
				pattern: nil,
			},
			wantMethod: "",
			wantPath:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotPath := HTTPPatternInfo(tt.args.pattern)
			if gotMethod != tt.wantMethod {
				t.Errorf("HTTPPatternInfo() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Errorf("HTTPPatternInfo() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestOpenSergo(t *testing.T) {
	srv := grpc.NewServer()
	srvContractPb.RegisterMetadataServiceServer(srv, new(testMetadataServiceServer))
	lis, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		t.Fatalf("net.Listen error:%s", err)
	}
	go func() {
		err := srv.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	app := &testAppInfo{
		name:     "testApp",
		endpoint: []string{"//example.com:9090", "//foo.com:9090"},
	}

	type args struct {
		opts []Option
	}
	tests := []struct {
		name      string
		args      args
		preFunc   func(t *testing.T)
		deferFunc func(t *testing.T)
		wantErr   bool
	}{
		{
			name: "test_with_opts",
			args: args{
				opts: []Option{
					WithEndpoint("127.0.0.1:9090"),
				},
			},
			wantErr: false,
		},
		{
			name: "test_with_env_endpoint",
			args: args{
				opts: []Option{},
			},
			preFunc: func(_ *testing.T) {
				err := os.Setenv("OPENSERGO_ENDPOINT", "127.0.0.1:9090")
				if err != nil {
					panic(err)
				}
			},
			wantErr: false,
		},
		{
			name: "test_with_env_config_file",
			args: args{
				opts: []Option{},
			},
			preFunc: func(_ *testing.T) {
				err := os.Setenv("OPENSERGO_BOOTSTRAP", `{"endpoint": "127.0.0.1:9090"}`)
				if err != nil {
					panic(err)
				}
			},
			wantErr: false,
		},
		{
			name: "test_with_env_bootstrap",
			args: args{
				opts: []Option{},
			},
			preFunc: func(t *testing.T) {
				fileContent := `{"endpoint": "127.0.0.1:9090"}`
				err := os.WriteFile("test.json", []byte(fileContent), 0o644)
				if err != nil {
					t.Fatalf("os.WriteFile error:%s", err)
				}
				confPath, err := filepath.Abs("./test.json")
				if err != nil {
					t.Fatalf("filepath.Abs error:%s", err)
				}
				err = os.Setenv("OPENSERGO_BOOTSTRAP_CONFIG", confPath)
				if err != nil {
					panic(err)
				}
			},
			deferFunc: func(t *testing.T) {
				path := os.Getenv("OPENSERGO_BOOTSTRAP_CONFIG")
				if path != "" {
					err := os.Remove(path)
					if err != nil {
						t.Fatalf("os.Remove error:%s", err)
					}
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFunc != nil {
				tt.preFunc(t)
			}
			if tt.deferFunc != nil {
				defer tt.deferFunc(t)
			}
			osServer, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = osServer.ReportMetadata(context.Background(), app)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReportMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
