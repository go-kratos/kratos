package discovery

import (
	"context"
	"os"
	"testing"

	"go-common/app/service/main/bns/agent/backend"
	"go-common/library/log"
)

func init() {
	log.Init(&log.Config{
		Stdout: true,
	})
}

var (
	// test discovery
	testURL    = "http://api.bilibili.co"
	testSecret = "b370880d1aca7d3a289b9b9a7f4d6812"
	testAppKey = "0c4b8fe3ff35a4b6"

	// test app
	testAppID          = "middleware.databus"
	testApplicationEnv = "uat"
	testZone           = "sh001"
	testRegion         = "sh"
)

var dis *discovery

func TestMain(m *testing.M) {
	config := map[string]interface{}{
		"url":    testURL,
		"secret": testSecret,
		"appKey": testAppKey,
	}
	backend, err := New(config)
	if err != nil {
		log.Error("new discovery error %s", err)
		os.Exit(1)
	}
	dis = backend.(*discovery)
	os.Exit(m.Run())
}

func TestNodes(t *testing.T) {
	ctx := context.Background()
	nodes, err := dis.Nodes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", nodes)
}

func TestQuery(t *testing.T) {
	ctx := context.Background()
	appID, sel, err := backend.ParseName(testAppID, backend.Selector{Env: testApplicationEnv, Region: testRegion, Zone: testZone})
	if err != nil {
		t.Fatal(err)
	}
	instances, err := dis.Query(ctx, appID, sel, backend.Metadata{
		ClientHost: "locahost",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", instances)
}

func BenchmarkQuery(b *testing.B) {
	ctx := context.Background()
	appID, sel, err := backend.ParseName(testAppID, backend.Selector{Env: testApplicationEnv, Region: testRegion, Zone: testZone})
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		data, err := dis.Query(ctx, appID, sel, backend.Metadata{ClientHost: "locahost"})
		if err != nil {
			b.Error(err)
		}
		if len(data) == 0 {
			b.Error("not data found")
		}
	}
}
