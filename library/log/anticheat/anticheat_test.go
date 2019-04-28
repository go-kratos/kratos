package anticheat

import (
	"sync"
	"testing"

	"go-common/library/log/infoc"
)

var (
	once sync.Once
	a    *AntiCheat
)

func onceInit() {
	a = New(&infoc.Config{
		TaskID:   "000146",
		Addr:     "172.16.0.204:514",
		Proto:    "tcp",
		ChanSize: 1,
	})
}

// go test  -test.v -test.bench Benchmark_InfoAntiCheat
// func Benchmark_InfoAntiCheat(b *testing.B) {
// 	once.Do(onceInit)
// 	client := httpx.NewClient(&httpx.ClientConfig{
// 		App: &conf.App{
// 			Key:    "appKey",
// 			Secret: "appSecret",
// 		},
// 		Timeout: 1,
// 	})
// 	params := url.Values{}
// 	params.Set("access_key", "infoc_access_key")
// 	params.Set("platform", "android")
// 	params.Set("build", "1111111")
// 	req, err := client.NewRequest("GET", "foo-api", "127.1.1.1", params)
// 	if err != nil {
// 		b.FailNow()
// 	}
// 	c := wctx.NewContext(ctx, req, nil, time.Millisecond*100)
// 	for j := 0; j < b.N; j++ {
// 		a.InfoAntiCheat(c, "infoc-test", "ip-address", "mid", "4", "5", "6", "7")
// 	}
// }

// go test  -test.v -test.bench Benchmark_ServiceAntiCheat
func Benchmark_ServiceAntiCheat(b *testing.B) {
	once.Do(onceInit)
	ac := map[string]string{
		"itemType": infoc.ItemTypeAv,
		"action":   infoc.ActionShare,
		"ip":       "remoteIP",
		"mid":      "mid",
		"fid":      "fid",
		"aid":      "aid",
		"sid":      "sid",
		"ua":       "ua",
		"buvid":    "buvid",
		"refer":    "refer",
		"url":      "infoc-test",
	}
	for j := 0; j < b.N; j++ {
		a.ServiceAntiCheat(ac)
	}
}
