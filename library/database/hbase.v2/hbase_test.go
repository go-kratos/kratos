package hbase

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tsuna/gohbase/hrpc"

	xtime "go-common/library/time"
)

var addrs []string
var client *Client

func TestMain(m *testing.M) {
	addrsStr := os.Getenv("HBASE_TEST_ADDRS")
	if addrsStr == "" {
		println("HBASE_TEST_ADDRS not set skip test !!")
		return
	}
	addrs = strings.Split(addrsStr, ",")
	config := &Config{
		Zookeeper: &ZKConfig{Root: "", Addrs: addrs, Timeout: xtime.Duration(time.Second)},
	}
	client = NewClient(config)
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	if err := client.Ping(context.Background()); err != nil {
		t.Errorf("ping meet err: %v", err)
	}
}

func TestPutGetDelete(t *testing.T) {
	ctx := context.Background()
	values := map[string]map[string][]byte{"name": {"firstname": []byte("hello"), "lastname": []byte("world")}}

	result, err := client.PutStr(ctx, "user", "user1", values)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", result)

	result, err = client.GetStr(ctx, "user", "user1")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Cells) != 2 {
		t.Errorf("unexpect result, expect 2 cell, get %d", len(result.Cells))
	}

	_, err = client.Delete(ctx, "user", "user1", values)
	if err != nil {
		t.Fatal(err)
	}

	result, err = client.GetStr(ctx, "user", "user1")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Cells) > 0 {
		t.Errorf("unexpect result, found cells")
	}
}

func TestScan(t *testing.T) {
	N := 10
	ctx := context.Background()
	values := map[string]map[string][]byte{"name": {"firstname": []byte("hello"), "lastname": []byte("world")}}
	for i := 0; i < N; i++ {
		_, err := client.PutStr(ctx, "user", fmt.Sprintf("scan_%d", i), values)
		if err != nil {
			t.Error(err)
		}
	}

	results, err := client.ScanStrAll(ctx, "user")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != N {
		t.Errorf("unexpect result expect %d result get %d", N, len(results))
	}

	iter, err := client.ScanStr(ctx, "user")
	if err != nil {
		t.Fatal(err)
	}
	defer iter.Close()
	GN := 0
	for {
		_, err := iter.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
		}
		GN++
	}
	if GN != N {
		t.Errorf("unexpect result expect %d result get %d", N, GN)
	}

	for i := 0; i < N; i++ {
		_, err := client.Delete(ctx, "user", fmt.Sprintf("scan_%d", i), nil)
		if err != nil {
			t.Errorf("delete error %s", err)
		}
	}
}

func TestScanRange(t *testing.T) {
	N := 10
	ctx := context.Background()
	values := map[string]map[string][]byte{"name": {"firstname": []byte("hello"), "lastname": []byte("world")}}
	for i := 0; i < N; i++ {
		_, err := client.PutStr(ctx, "user", fmt.Sprintf("scan_%d", i), values)
		if err != nil {
			t.Error(err)
		}
	}

	scanner, err := client.ScanRangeStr(ctx, "user", "scan_0", "scan_3")
	if err != nil {
		t.Fatal(err)
	}
	var results []*hrpc.Result
	for {
		result, err := scanner.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		results = append(results, result)
	}
	if len(results) != 3 {
		t.Errorf("unexpect result expect %d result get %d", N, len(results))
	}

	for i := 0; i < N; i++ {
		_, err := client.Delete(ctx, "user", fmt.Sprintf("scan_%d", i), nil)
		if err != nil {
			t.Errorf("delete error %s", err)
		}
	}
}

func TestClose(t *testing.T) {
	if err := client.Close(); err != nil {
		t.Logf("Close meet error: %v", err)
	}
	if err := client.Ping(context.Background()); err == nil {
		t.Errorf("ping return nil error after being closed")
	}
}
