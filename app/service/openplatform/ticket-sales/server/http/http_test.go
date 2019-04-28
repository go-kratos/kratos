package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"go-common/app/service/openplatform/ticket-sales/conf"
	"go-common/app/service/openplatform/ticket-sales/service"
	"go-common/library/conf/paladin"
	httpx "go-common/library/net/http/blademaster"
)

var (
	once         sync.Once
	client       *httpx.Client
	syncOrderURL string
)

func startHTTP() {
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	if err := paladin.Watch("ticket-sales.toml", conf.Conf); err != nil {
		panic(err)
	}
	client = httpx.NewClient(conf.Conf.HTTPClient.Read)
	Init(conf.Conf, service.New(conf.Conf))
}

func TestGetOrder(t *testing.T) {
	syncOrderURL = conf.Conf.UT.DistPrefix + "/distrib/getorder"
	fmt.Println(syncOrderURL)
	once.Do(startHTTP)
	params := url.Values{}
	params.Set("oid", "100000004421700")
	var err error
	var req *http.Request
	if req, err = client.NewRequest("GET", syncOrderURL, "127.0.0.1", params); err != nil {
		t.Errorf("http.NewRequest(GET, %s) error(%v)", syncOrderURL, err)
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
