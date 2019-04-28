package http

import (
	"context"
	"net/url"
	"testing"
	"time"

	"go-common/app/admin/main/reply/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

const (
	_monitorState = "http://127.0.0.1:6711/x/internal/replyadmin/monitor/state"
)

func TestHttp(t *testing.T) {
	var (
		err error
	)
	if err = conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	client := bm.NewClient(&bm.ClientConfig{
		Dial:    xtime.Duration(time.Second),
		Timeout: xtime.Duration(time.Second),
	})
	Init(conf.Conf)
	// test
	testMonitorState(client, t)
}

func testMonitorState(client *bm.Client, t *testing.T) {
	var err error
	params := url.Values{}
	params.Set("type", "1")
	params.Set("oid", "1")
	params.Set("state", "1")
	params.Set("adid", "11")
	params.Set("remark", "test")
	// send
	res := map[string]interface{}{}
	if err = client.Post(context.Background(), _monitorState, "", params, &res); err != nil {
		t.Errorf("client.Post() error(%v)", err)
		t.FailNow()
	}
	validRes(_monitorState, params, res, t)
}

func validRes(url string, params url.Values, res map[string]interface{}, t *testing.T) {
	if code, ok := res["code"]; ok && code.(float64) == 0 {
		t.Logf("\nurl:%s\nparams:%s\nres:%v", url, params.Encode(), res)
	} else {
		t.Errorf("\nurl:%s\nparams:%s\nres:%v", url, params.Encode(), res)
	}
}
