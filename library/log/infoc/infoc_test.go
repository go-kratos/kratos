package infoc

import (
	"bytes"
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"go-common/library/net/metadata"
)

var (
	once sync.Once
	i1   *Infoc
)

func TestMain(m *testing.M) {
	once.Do(createInfoc)
	defer i1.Close()
	os.Exit(m.Run())
}

func createInfoc() {
	i1 = New(&Config{
		TaskID:   "000146",
		Addr:     "172.16.0.204:514",
		Proto:    "tcp",
		ChanSize: 1,
	})
}

func Test_Infoc(b *testing.T) {

	err := i1.Info("infoc-test", "ip", "mid", 222)
	time.Sleep(2 * time.Second)
	if err != nil {
		b.Fatalf("err %+v", err)
	}
}

func Test_Infocv(b *testing.T) {
	i1.Infov(context.Background(), "infoc-test", "ip", "mid", 222)

	ctx := metadata.NewContext(context.Background(), metadata.MD{metadata.Mirror: true})
	i1.Infov(ctx, "infoc-test", "ip", "mid", 222)

	ctx = metadata.NewContext(context.Background(), metadata.MD{metadata.Mirror: "1"})
	err := i1.Infov(ctx, "infoc-test", "ip", "mid", 222)
	time.Sleep(2 * time.Second)
	if err != nil {
		b.Fatalf("err %+v", err)
	}

	var args []interface{}
	args = append(args, "infoc-test")
	args = append(args, "ip")
	args = append(args, "mid")
	args = append(args, 222)
	err1, buf1 := i1.info(args...)
	err2, buf2 := i1.info(args)

	if bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		b.Fatalf("err %+v,%+v,%+v,%+v", err1, err2, buf1.Bytes(), buf2.Bytes())
	}

	args = append([]interface{}{})
	args = append(args, "infoc-test")

	err1, buf1 = i1.info(args...)
	err2, buf2 = i1.info(args)

	if bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		b.Fatalf("err %+v,%+v,%+v,%+v", err1, err2, buf1.Bytes(), buf2.Bytes())
	}

}

func BenchmarkInfoc(b *testing.B) {
	once.Do(createInfoc)
	b.RunParallel(func(pb *testing.PB) {
		var f float32 = 3.55051
		var i8 int8 = 2
		var u8 uint8 = 2
		for pb.Next() {
			i1.Info("infoc-test", "ip", "mid", i8, u8, f)
		}
	})
}
