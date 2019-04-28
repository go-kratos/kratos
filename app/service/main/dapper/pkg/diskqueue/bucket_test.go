package diskqueue

import (
	"crypto/rand"
	"io"
	"os"
	"reflect"
	"testing"
)

func Test_membucket(t *testing.T) {
	cap := int32(16)
	data := make([]byte, _blockByte*cap)
	mb := &memBucket{
		cap:  cap,
		data: data,
	}
	t.Run("test push & pop small data", func(t *testing.T) {
		p := []byte("hello world")
		err := mb.push(p)
		if err != nil {
			t.Error(err)
		}
		ret, err := mb.pop()
		if err != nil {
			t.Error(err)
		} else {
			if !reflect.DeepEqual(ret, p) {
				t.Errorf("%s not equal %s", ret, p)
			}
		}
	})
	t.Run("test push & pop big data", func(t *testing.T) {
		p := make([]byte, 1890)
		rand.Read(p)
		err := mb.push(p)
		if err != nil {
			t.Error(err)
		}
		ret, err := mb.pop()
		if err != nil {
			t.Error(err)
		} else {
			if !reflect.DeepEqual(ret, p) {
				t.Logf("buf: %v", mb.data)
				t.Errorf("%v not equal %v", ret, p)
			}
		}
	})
	t.Run("push big data", func(t *testing.T) {
		p := make([]byte, _blockByte*cap*2)
		err := mb.push(p)
		if err != errBucketFull {
			t.Errorf("expect err == errBucketFull get: %v", err)
		}
	})
	t.Run("pop io.EOF", func(t *testing.T) {
		_, err := mb.pop()
		if err != io.EOF {
			t.Errorf("expect err == io.EOF get: %v", err)
		}
	})
}

func Test_fileBucket(t *testing.T) {
	fpath := "bucket.bin"
	defer os.RemoveAll(fpath)
	cap := int32(16)
	data := make([]byte, _blockByte*cap)
	mb := &memBucket{
		cap:  cap,
		data: data,
	}
	d1 := []byte("hello world")
	for i := 0; i < 10; i++ {
		mb.push(d1)
	}
	fp, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	mb.dump(fp)
	fp.Close()
	fb, err := newFileBucket(fpath)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for {
		ret, err := fb.pop()
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
		count++
		if !reflect.DeepEqual(ret, d1) {
			t.Errorf("%v not equal %v", ret, d1)
		}
	}
	if count != 10 {
		t.Errorf("expect 10 data get %d", count)
	}
}
