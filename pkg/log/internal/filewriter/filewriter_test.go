package filewriter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const logdir = "testlog"

func touch(dir, name string) {
	os.MkdirAll(dir, 0755)
	fp, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	fp.Close()
}

func TestMain(m *testing.M) {
	ret := m.Run()
	os.RemoveAll(logdir)
	os.Exit(ret)
}

func TestParseRotate(t *testing.T) {
	touch := func(dir, name string) {
		os.MkdirAll(dir, 0755)
		fp, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		fp.Close()
	}
	dir := filepath.Join(logdir, "test-parse-rotate")
	names := []string{"info.log.2018-11-11", "info.log.2018-11-11.001", "info.log.2018-11-11.002", "info.log." + time.Now().Format("2006-01-02") + ".005"}
	for _, name := range names {
		touch(dir, name)
	}
	l, err := parseRotateItem(dir, "info.log", "2006-01-02")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(names), l.Len())

	rt := l.Front().Value.(rotateItem)

	assert.Equal(t, 5, rt.rotateNum)
}

func TestRotateExists(t *testing.T) {
	dir := filepath.Join(logdir, "test-rotate-exists")
	names := []string{"info.log." + time.Now().Format("2006-01-02") + ".005"}
	for _, name := range names {
		touch(dir, name)
	}
	fw, err := New(logdir+"/test-rotate-exists/info.log",
		MaxSize(1024*1024),
		func(opt *option) { opt.RotateInterval = time.Millisecond },
	)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	for i := 0; i < 10; i++ {
		for i := 0; i < 1024; i++ {
			_, err = fw.Write(data)
			if err != nil {
				t.Error(err)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	fw.Close()
	fis, err := ioutil.ReadDir(logdir + "/test-rotate-exists")
	if err != nil {
		t.Fatal(err)
	}
	var fnams []string
	for _, fi := range fis {
		fnams = append(fnams, fi.Name())
	}
	assert.Contains(t, fnams, "info.log."+time.Now().Format("2006-01-02")+".006")
}

func TestSizeRotate(t *testing.T) {
	fw, err := New(logdir+"/test-rotate/info.log",
		MaxSize(1024*1024),
		func(opt *option) { opt.RotateInterval = 1 * time.Millisecond },
	)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	for i := 0; i < 10; i++ {
		for i := 0; i < 1024; i++ {
			_, err = fw.Write(data)
			if err != nil {
				t.Error(err)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	fw.Close()
	fis, err := ioutil.ReadDir(logdir + "/test-rotate")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(fis) > 5, "expect more than 5 file get %d", len(fis))
}

func TestMaxFile(t *testing.T) {
	fw, err := New(logdir+"/test-maxfile/info.log",
		MaxSize(1024*1024),
		MaxFile(1),
		func(opt *option) { opt.RotateInterval = 1 * time.Millisecond },
	)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	for i := 0; i < 10; i++ {
		for i := 0; i < 1024; i++ {
			_, err = fw.Write(data)
			if err != nil {
				t.Error(err)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	fw.Close()
	fis, err := ioutil.ReadDir(logdir + "/test-maxfile")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(fis) <= 2, fmt.Sprintf("expect 2 file get %d", len(fis)))
}

func TestMaxFile2(t *testing.T) {
	files := []string{
		"info.log.2018-12-01",
		"info.log.2018-12-02",
		"info.log.2018-12-03",
		"info.log.2018-12-04",
		"info.log.2018-12-05",
		"info.log.2018-12-05.001",
	}
	for _, file := range files {
		touch(logdir+"/test-maxfile2", file)
	}
	fw, err := New(logdir+"/test-maxfile2/info.log",
		MaxSize(1024*1024),
		MaxFile(3),
		func(opt *option) { opt.RotateInterval = 1 * time.Millisecond },
	)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	for i := 0; i < 10; i++ {
		for i := 0; i < 1024; i++ {
			_, err = fw.Write(data)
			if err != nil {
				t.Error(err)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	fw.Close()
	fis, err := ioutil.ReadDir(logdir + "/test-maxfile2")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(fis) == 4, fmt.Sprintf("expect 4 file get %d", len(fis)))
}

func TestFileWriter(t *testing.T) {
	fw, err := New("testlog/info.log")
	if err != nil {
		t.Fatal(err)
	}
	defer fw.Close()
	_, err = fw.Write([]byte("Hello World!\n"))
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkFileWriter(b *testing.B) {
	fw, err := New("testlog/bench/info.log",
		func(opt *option) { opt.WriteTimeout = time.Second }, MaxSize(1024*1024*8), /*32MB*/
		func(opt *option) { opt.RotateInterval = 10 * time.Millisecond },
	)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_, err = fw.Write([]byte("Hello World!\n"))
		if err != nil {
			b.Error(err)
		}
	}
}
