package metadata

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Service is description service
type Service struct {
	grpcSer  *grpc.Server
	once     sync.Once
	services map[string]*descriptorpb.FileDescriptorSet
}

// NewService create desc Service
func NewService(grpcSrv ...*grpc.Server) *Service {
	s := &Service{}
	if len(grpcSrv) > 0 {
		s.grpcSer = grpcSrv[0]
	}
	return s
}

// ListServices list the full name of all services
func (s *Service) ListServices(ctx context.Context) (services []string, err error) {
	services = []string{}
	s.once.Do(func() {
		err = s.initServices()
	})
	for name := range s.services {
		services = append(services, name)
	}
	return
}

// GetServiceMeta get the full fileDescriptorSet of service
func (s *Service) GetServiceMeta(ctx context.Context, name string) (fds *descriptorpb.FileDescriptorSet, err error) {
	fds = &descriptorpb.FileDescriptorSet{}
	s.once.Do(func() {
		err = s.initServices()
	})
	if temp, ok := s.services[name]; ok {
		fds = temp
	}
	return
}

func (s *Service) initServices() error {
	serviceProto, err := s.listServices()
	if err != nil {
		s.services = make(map[string]*descriptorpb.FileDescriptorSet)
		return err
	}
	s.services = serviceProto
	return nil
}

func (s *Service) listServices() (map[string]*descriptorpb.FileDescriptorSet, error) {
	services := make(map[string]*descriptorpb.FileDescriptorSet, 0)
	if s.grpcSer != nil {
		for svc, info := range s.grpcSer.GetServiceInfo() {
			fdenc, ok := parseMetadata(info.Metadata)
			if !ok {
				return nil, fmt.Errorf("invalid service %s meta", svc)
			}
			fd, err := decodeFileDesc(fdenc)
			if err != nil {
				return nil, err
			}
			protoSet, err := allDependency(fd)
			if err != nil {
				return nil, err
			}
			services[svc] = &dpb.FileDescriptorSet{File: protoSet}
		}
		return services, nil
	}
	var err error
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if fd.Services() != nil && fd.Services().Len() > 0 {
			for i := 0; i < fd.Services().Len(); i++ {
				svc := string(fd.Services().Get(i).FullName())
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
				services[svc] = &dpb.FileDescriptorSet{File: fdps}
			}
		}
		return true
	})
	return services, err
}

func fileDescriptorProto(path string) (*dpb.FileDescriptorProto, error) {
	fdenc := proto.FileDescriptor(path)
	fdDep, err := decodeFileDesc(fdenc)
	if err != nil {
		return nil, err
	}
	return fdDep, nil
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

// parseMetadata finds the file descriptor bytes specified meta.
// For SupportPackageIsVersion4, m is the name of the proto file, we
// call proto.FileDescriptor to get the byte slice.
// For SupportPackageIsVersion3, m is a byte slice itself.
func parseMetadata(meta interface{}) ([]byte, bool) {
	// Check if meta is the file name.
	if fileNameForMeta, ok := meta.(string); ok {
		return proto.FileDescriptor(fileNameForMeta), true
	}

	// Check if meta is the byte slice.
	if enc, ok := meta.([]byte); ok {
		return enc, true
	}

	return nil, false
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
