package metadata

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

//go:generate protoc --proto_path=. --proto_path=../../third_party --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --go-http_out=paths=source_relative:. metadata.proto

// Server is api meta server
type Server struct {
	UnimplementedMetadataServer

	srv      *grpc.Server
	lock     sync.Mutex
	services map[string]*dpb.FileDescriptorSet
	methods  map[string][]string
}

// NewServer create server instance
func NewServer(srv *grpc.Server) *Server {
	return &Server{
		srv:      srv,
		services: make(map[string]*dpb.FileDescriptorSet),
		methods:  make(map[string][]string),
	}
}

func (s *Server) load() error {
	if len(s.services) > 0 {
		return nil
	}
	if s.srv != nil {
		for name, info := range s.srv.GetServiceInfo() {
			fd, err := parseMetadata(info.Metadata)
			if err != nil {
				return fmt.Errorf("invalid service %s metadata err:%v", name, err)
			}
			protoSet, err := allDependency(fd)
			if err != nil {
				return err
			}
			s.services[name] = &dpb.FileDescriptorSet{File: protoSet}
			for _, method := range info.Methods {
				s.methods[name] = append(s.methods[name], method.Name)
			}
		}
		return nil
	}
	var err error
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if fd.Services() != nil {
			for i := 0; i < fd.Services().Len(); i++ {
				svc := fd.Services().Get(i)
				fdp, e := fileDescriptorProto(fd.Path())
				if e != nil {
					err = e
					return false
				}
				fdps, e := allDependency(fdp)
				if e != nil {
					err = e
					return false
				}
				s.services[string(svc.FullName())] = &dpb.FileDescriptorSet{File: fdps}
				if svc.Methods() != nil {
					for j := 0; j < svc.Methods().Len(); j++ {
						method := svc.Methods().Get(j)
						s.methods[string(svc.FullName())] = append(s.methods[string(svc.FullName())], string(method.Name()))
					}
				}
			}
		}
		return true
	})
	return err
}

// ListServices return all services
func (s *Server) ListServices(ctx context.Context, in *ListServicesRequest) (*ListServicesReply, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if err := s.load(); err != nil {
		return nil, err
	}
	reply := new(ListServicesReply)
	for name := range s.services {
		reply.Services = append(reply.Services, name)
	}
	for name, methods := range s.methods {
		for _, method := range methods {
			reply.Methods = append(reply.Methods, fmt.Sprintf("/%s/%s", name, method))
		}
	}
	return reply, nil
}

// GetServiceDesc return service meta by name
func (s *Server) GetServiceDesc(ctx context.Context, in *GetServiceDescRequest) (*GetServiceDescReply, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if err := s.load(); err != nil {
		return nil, err
	}
	fds, ok := s.services[in.Name]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "service %s not found", in.Name)
	}
	return &GetServiceDescReply{FileDescSet: fds}, nil
}

// parseMetadata finds the file descriptor bytes specified meta.
// For SupportPackageIsVersion4, m is the name of the proto file, we
// call proto.FileDescriptor to get the byte slice.
// For SupportPackageIsVersion3, m is a byte slice itself.
func parseMetadata(meta interface{}) (*dpb.FileDescriptorProto, error) {
	// Check if meta is the file name.
	if fileNameForMeta, ok := meta.(string); ok {
		return fileDescriptorProto(fileNameForMeta)
	}
	// Check if meta is the byte slice.
	if enc, ok := meta.([]byte); ok {
		fd, err := decodeFileDesc(enc)
		if err != nil {
			return nil, err
		}
		return fd, nil
	}
	return nil, fmt.Errorf("proto not sumpport metadata: %v", meta)
}

// decodeFileDesc does decompression and unmarshalling on the given
// file descriptor byte slice.
func decodeFileDesc(enc []byte) (*dpb.FileDescriptorProto, error) {
	raw, err := decompress(enc)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress enc: %v", err)
	}
	fd := new(dpb.FileDescriptorProto)
	if err := proto.Unmarshal(raw, fd); err != nil {
		return nil, fmt.Errorf("bad descriptor: %v", err)
	}
	return fd, nil
}

func allDependency(fd *dpb.FileDescriptorProto) ([]*dpb.FileDescriptorProto, error) {
	var files []*dpb.FileDescriptorProto
	for _, dep := range fd.Dependency {
		fdDep, err := fileDescriptorProto(dep)
		if err != nil {
			return nil, err
		}
		temp, err := allDependency(fdDep)
		if err != nil {
			return nil, err
		}
		files = append(files, temp...)
	}
	files = append(files, fd)
	return files, nil
}

// decompress does gzip decompression.
func decompress(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	return out, nil
}

func fileDescriptorProto(path string) (*dpb.FileDescriptorProto, error) {
	fd, err := protoregistry.GlobalFiles.FindFileByPath(path)
	if err != nil {
		return nil, err
	}
	fdpb := protodesc.ToFileDescriptorProto(fd)
	return fdpb, nil
}
