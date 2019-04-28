package tcpcollect

import (
	"context"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"go-common/app/service/main/dapper/conf"
	"go-common/app/service/main/dapper/model"
	"go-common/app/service/main/dapper/pkg/process"
)

func TestCollect(t *testing.T) {
	count := 0

	collect := New(&conf.Collect{Network: "tcp", Addr: "127.0.0.1:6190"})
	collect.RegisterProcess(process.MockProcess(func(context.Context, *model.ProtoSpan) error {
		count++
		return nil
	}))
	if err := collect.Start(); err != nil {
		t.Fatal(err)
	}
	fp, err := os.Open("testdata/data.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()
	conn, err := net.Dial("tcp", "127.0.0.1:6190")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = io.Copy(conn, fp)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second)
	if count <= 0 {
		t.Errorf("expect more than one span write")
	}
}
