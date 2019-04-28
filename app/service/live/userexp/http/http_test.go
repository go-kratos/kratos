package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"go-common/app/service/live/userexp/conf"
	"go-common/app/service/live/userexp/model"
	"go-common/app/service/live/userexp/service"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_getLevelURL      = "http://localhost:8801/x/internal/liveexp/level/get"
	_multiGetLevelURL = "http://localhost:8801/x/internal/liveexp/level/mulGet"
	_addUexpURL       = "http://localhost:8801/x/internal/liveexp/level/addUexp"
	_addRexpURL       = "http://localhost:8801/x/internal/liveexp/level/addRexp"
)

var (
	once   sync.Once
	client *httpx.Client
)

func startHTTP() {
	if err := conf.Init(); err != nil {
		panic(fmt.Errorf("conf.Init() error(%v)", err))
	}
	svr := service.New(conf.Conf)
	client = httpx.NewClient(conf.Conf.HTTPClient.Read)
	Init(conf.Conf, svr)
}

func TestGetLevel(t *testing.T) {
	once.Do(startHTTP)
	params := url.Values{}
	params.Set("uid", "10001")
	var err error
	var req *http.Request
	if req, err = client.NewRequest("GET", _getLevelURL, "127.0.0.1", params); err != nil {
		t.Errorf("http.NewRequest(GET, %s) error(%v)", _getLevelURL, err)
		t.FailNow()
	}
	reqURI := req.URL.String()
	t.Logf("req uri: %s\n", reqURI)
	var res struct {
		Code int          `json:"code"`
		Data *model.Level `json:"data"`
	}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	str, _ := json.Marshal(res)
	t.Logf("res: %+v", string(str))
}

func TestMultiGetLevel(t *testing.T) {
	once.Do(startHTTP)
	params := url.Values{}
	params.Set("uids", "10001,10002,10003")
	var err error
	var req *http.Request
	if req, err = client.NewRequest("GET", _multiGetLevelURL, "127.0.0.1", params); err != nil {
		t.Errorf("http.NewRequest(GET, %s) error(%v)", _multiGetLevelURL, err)
		t.FailNow()
	}
	reqURI := req.URL.String()
	t.Logf("req uri: %s\n", reqURI)
	var res struct {
		Code int                     `json:"code"`
		Data map[string]*model.Level `json:"data"`
	}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	str, _ := json.Marshal(res)
	t.Logf("res: %+v", string(str))
}

func TestAddUexp(t *testing.T) {
	once.Do(startHTTP)
	params := url.Values{}
	params.Set("uid", "10001")
	params.Set("uexp", "1234")
	var err error
	var req *http.Request
	if req, err = client.NewRequest("GET", _addUexpURL, "127.0.0.1", params); err != nil {
		t.Errorf("http.NewRequest(GET, %s) error(%v)", _addUexpURL, err)
		t.FailNow()
	}
	reqURI := req.URL.String()
	t.Logf("req uri: %s\n", reqURI)
	var res struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	str, _ := json.Marshal(res)
	t.Logf("res: %+v", string(str))
}

func TestAddRexp(t *testing.T) {
	once.Do(startHTTP)
	params := url.Values{}
	params.Set("uid", "10001")
	params.Set("rexp", "4321")
	var err error
	var req *http.Request
	if req, err = client.NewRequest("GET", _addRexpURL, "127.0.0.1", params); err != nil {
		t.Errorf("http.NewRequest(GET, %s) error(%v)", _addRexpURL, err)
		t.FailNow()
	}
	reqURI := req.URL.String()
	t.Logf("req uri: %s\n", reqURI)
	var res struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	str, _ := json.Marshal(res)
	t.Logf("res: %+v", string(str))
}
