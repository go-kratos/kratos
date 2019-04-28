package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

//LittleEndian
func getUint16(b []byte, offset int) uint16 {
	_ = b[offset+1] // early bounds check
	return uint16(b[offset+0]) |
		uint16(b[offset+1])<<8
}

//LittleEndian
func getUint32(b []byte, offset int) uint32 {
	_ = b[offset+3] // early bounds check
	return uint32(b[offset+0]) |
		uint32(b[offset+1])<<8 |
		uint32(b[offset+2])<<16 |
		uint32(b[offset+3])<<24
}

//LittleEndian
func getUint64(b []byte, offset int) uint64 {
	_ = b[offset+7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[offset+0]) |
		uint64(b[offset+1])<<8 |
		uint64(b[offset+2])<<16 |
		uint64(b[offset+3])<<24 |
		uint64(b[offset+4])<<32 |
		uint64(b[offset+5])<<40 |
		uint64(b[offset+6])<<48 |
		uint64(b[offset+7])<<56
}

//LittleEndian
func putUint64(v uint64, b []byte, offset int) {
	_ = b[offset+7] // early bounds check to guarantee safety of writes below
	b[offset+0] = byte(v)
	b[offset+1] = byte(v >> 8)
	b[offset+2] = byte(v >> 16)
	b[offset+3] = byte(v >> 24)
	b[offset+4] = byte(v >> 32)
	b[offset+5] = byte(v >> 40)
	b[offset+6] = byte(v >> 48)
	b[offset+7] = byte(v >> 56)
}

//LittleEndian
func putUint32(v uint32, b []byte, offset int) {
	_ = b[offset+3]
	b[offset+0] = byte(v)
	b[offset+1] = byte(v >> 8)
	b[offset+2] = byte(v >> 16)
	b[offset+3] = byte(v >> 24)
}

func copyBytes(src []byte, srcStart int, dst []byte, dstStart int, count int) (int, error) {
	if len(src) < srcStart+count || len(dst) < dstStart+count {
		return -1, fmt.Errorf("Array index out of bounds")
	}
	for i := 0; i < count; i++ {
		dst[dstStart+i] = src[srcStart+i]
	}
	return count, nil
}

func fileNameAndExt(path string) (string, string) {
	name := filepath.Base(path)
	for i := len(name) - 1; i >= 0 && !os.IsPathSeparator(name[i]); i-- {
		if name[i] == '.' {
			return name[:i], name[i:]
		}
	}
	return name, ""
}
