package diskqueue

import (
	"bytes"
	"crypto/rand"
	"io"
	mrand "math/rand"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"
)

func init() {
	mrand.Seed(time.Now().UnixNano())
}

func TestDiskQueuePushPopMem(t *testing.T) {
	dirname := "testdata/d1"
	defer os.RemoveAll(dirname)
	queue, err := New(dirname)
	if err != nil {
		t.Fatal(err)
	}
	N := 10
	p := []byte("hello world")
	for i := 0; i < N; i++ {
		if err := queue.Push(p); err != nil {
			t.Error(err)
		}
	}
	count := 0
	for {
		data, err := queue.Pop()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(data, p) {
			t.Errorf("invalid data: %s", data)
		}
		count++
	}
	if count != N {
		t.Errorf("wrong count %d", count)
	}
}

func TestDiskQueueDisk(t *testing.T) {
	data := make([]byte, 2233)
	rand.Read(data)
	count := 1024 * 256
	dirname := "testdata/d2"
	defer os.RemoveAll(dirname)
	t.Run("test write disk", func(t *testing.T) {
		queue, err := New(dirname)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < count; i++ {
			if err := queue.Push(data); err != nil {
				time.Sleep(time.Second)
				if err := queue.Push(data); err != nil {
					t.Error(err)
				}
			}
		}
		queue.Close()
	})
	t.Run("test read disk", func(t *testing.T) {
		n := 0
		queue, err := New(dirname)
		if err != nil {
			t.Fatal(err)
		}
		for {
			ret, err := queue.Pop()
			if err == io.EOF {
				break
			}
			if !bytes.Equal(data, ret) {
				t.Errorf("invalid data unequal")
			}
			n++
		}
		if n != count {
			t.Errorf("want %d get %d", count, n)
		}
	})
}

func TestDiskQueueTrans(t *testing.T) {
	dirname := "testdata/d3"
	defer os.RemoveAll(dirname)
	queue, err := New(dirname)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 1890)
	rand.Read(data)
	cycles := 512
	var wg sync.WaitGroup
	wg.Add(2)
	done := false
	writed := 0
	readed := 0
	go func() {
		defer wg.Done()
		for i := 0; i < cycles; i++ {
			ms := mrand.Intn(40) + 10
			time.Sleep(time.Duration(ms) * time.Millisecond)
			for i := 0; i < 128; i++ {
				if err := queue.Push(data); err != nil {
					t.Error(err)
				} else {
					writed++
				}
			}
		}
		done = true
	}()
	go func() {
		defer wg.Done()
		for {
			ret, err := queue.Pop()
			if err == io.EOF && done {
				break
			}
			if err == io.EOF {
				ms := mrand.Intn(10)
				time.Sleep(time.Duration(ms) * time.Millisecond)
				continue
			}
			if !bytes.Equal(ret, data) {
				t.Fatalf("invalid data, data length: %d, want: %d, data: %v, want: %v", len(ret), len(data), ret, data)
			}
			readed++
		}
	}()
	wg.Wait()
	os.RemoveAll(dirname)
	if writed != readed {
		t.Errorf("readed: %d != writed: %d", readed, writed)
	}
}

func TestEmpty(t *testing.T) {
	dirname := "testdata/d4"
	defer os.RemoveAll(dirname)
	queue, err := New(dirname)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		_, err := queue.Pop()
		if err != io.EOF {
			t.Errorf("expect err == io.EOF, get %v", err)
		}
	}
}

func TestEmptyCache(t *testing.T) {
	datadir := "testdata/emptycache"
	dirname := "testdata/de"
	if err := exec.Command("cp", "-r", datadir, dirname).Run(); err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirname)
	queue, err := New(dirname)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		_, err := queue.Pop()
		if err != io.EOF {
			t.Errorf("expect err == io.EOF, get %v", err)
		}
	}
}

func BenchmarkDiskQueue(b *testing.B) {
	queue, err := New("testdata/d5")
	if err != nil {
		b.Fatal(err)
	}
	done := make(chan bool, 1)
	go func() {
		for {
			if _, err := queue.Pop(); err != nil {
				if err == io.EOF {
					break
				}
			}
		}
		done <- true
	}()
	data := make([]byte, 768)
	rand.Read(data)
	for i := 0; i < b.N; i++ {
		queue.Push(data)
	}
	<-done
}
