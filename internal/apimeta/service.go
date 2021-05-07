package apimeta

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func allDependency(fd *dpb.FileDescriptorProto) ([]*dpb.FileDescriptorProto, error) {
	var files []*dpb.FileDescriptorProto
	files = append(files, fd)
	for _, dep := range fd.Dependency {
		fdenc := proto.FileDescriptor(dep)
		fdDep, err := decodeFileDesc(fdenc)
		if err != nil {
			return nil, err
		}
		temp, err := allDependency(fdDep)
		if err != nil {
			return nil, err
		}
		files = append(files, temp...)
	}
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
