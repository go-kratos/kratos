package jump

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
)

//Hash get result by hash
func Hash(key uint64, numBuckets int) int32 {
	var b int64 = -1
	var j int64

	for j < int64(numBuckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}

	return int32(b)
}

//Md5 get result by Md5
func Md5(key string) uint64 {
	var x uint64
	s := md5.Sum([]byte(key))
	b := bytes.NewBuffer(s[:])
	binary.Read(b, binary.BigEndian, &x)
	return x
}
