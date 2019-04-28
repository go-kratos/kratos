package dao

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"go-common/app/interface/main/upload/conf"
	xtime "go-common/library/time"
)

var (
	b        *Bfs
	testData []byte
	once     sync.Once
)

func TestMain(m *testing.M) {
	initData()
	once.Do(initBFS)
	os.Exit(m.Run())
}

func initBFS() {
	b = NewBfs(&conf.Config{
		Bfs: &conf.Bfs{
			BfsURL:          "uat-bfs.bilibili.co",
			WaterMarkURL:    "http://172.18.33.121:8090/imageserver/watermark/gen",
			ImageGenURL:     "http://172.18.33.121:8090/imageserver/image/gen",
			TimeOut:         xtime.Duration(time.Second * 5),
			WmTimeOut:       xtime.Duration(time.Second * 5),
			ImageGenTimeOut: xtime.Duration(time.Second * 5),
		},
	})
}

func initData() {
	client := &http.Client{}
	resp, err := client.Get("http://uat-i0.hdslb.com/bfs/archive/fc7cd08beb29f93c596426557cf1aa11a08e9730.jpg")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if testData, err = ioutil.ReadAll(resp.Body); err != nil {
		panic(err)
	}
}
