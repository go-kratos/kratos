package http

// import (
// 	"context"
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"testing"
// 	"time"

// 	"go-common/app/admin/main/apm/conf"
// 	"go-common/app/admin/main/apm/service"
// 	bm "go-common/library/net/http/context"
// )

// var path = "http://127.0.0.1:80/xx"

// type resp struct {
// 	http.ResponseWriter
// }

// func TestNotifyApply(t *testing.T) {
// 	flag.Set("conf", "../cmd/apm-admin-test.toml")
// 	conf.Init()
// 	apmSvc = service.New(conf.Conf)
// 	params := url.Values{}
// 	params.Set("cluster", "test_kafka_9092-266")
// 	params.Set("topic_name", "AccLabour-T")
// 	params.Set("remark", "test")
// 	params.Set("project", "main.web-svr")
// 	params.Set("topic_remark", "test")
// 	params.Set("offset", "new")
// 	params.Set("filter", "1")
// 	params.Set("concurrent", "10")
// 	params.Set("callback", `{"aa":1}`)
// 	params.Set("filters", `[{"field":"svv","condition":1,"value":"v"}]`)
// 	req, _ := http.NewRequest("POST", path, strings.NewReader(params.Encode()))
// 	req.Header.Set("X-BACKEND-BILI-REAL-IP", "xxx")
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.ParseForm()
// 	ctx := bm.NewContext(context.TODO(), req, resp{}, time.Second)
// 	ctx.Set("username", "haoguanwei")
// 	ctx.Request().ParseForm()
// 	//databusNotifyApplyAdd(ctx)
// 	fmt.Println(ctx.Result())
// }

// func TestApplyList(t *testing.T) {
// 	flag.Set("conf", "../cmd/apm-admin-test.toml")
// 	conf.Init()
// 	apmSvc = service.New(conf.Conf)
// 	params := url.Values{}
// 	//	params.Set("cluster", "test_kafka_9092-266")
// 	//	params.Set("topic_name", "AccAnswer-T")
// 	//	params.Set("project", "main.web-svr")
// 	req, _ := http.NewRequest("POST", path, strings.NewReader(params.Encode()))
// 	req.Header.Set("X-BACKEND-BILI-REAL-IP", "xxx")
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.ParseForm()
// 	ctx := bm.NewContext(context.TODO(), req, resp{}, time.Second)
// 	ctx.Set("username", "haoguanwei")
// 	ctx.Request().ParseForm()
// 	//databusApplyList(ctx)

// 	data, err := json.Marshal(ctx.Result())
// 	fmt.Println(string(data), err)
// }

// func TestNotifyList(t *testing.T) {
// 	flag.Set("conf", "../cmd/apm-admin-test.toml")
// 	conf.Init()
// 	apmSvc = service.New(conf.Conf)
// 	params := url.Values{}
// 	params.Set("cluster", "test_kafka_9092-266")
// 	params.Set("topic_name", "AccAnswer-T")
// 	params.Set("remark", "test")
// 	//params.Set("project", "main.web-svr")
// 	req, _ := http.NewRequest("POST", path, strings.NewReader(params.Encode()))
// 	req.Header.Set("X-BACKEND-BILI-REAL-IP", "xxx")
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.ParseForm()
// 	ctx := bm.NewContext(context.TODO(), req, resp{}, time.Second)
// 	ctx.Set("username", "haoguanwei")
// 	ctx.Request().ParseForm()
// 	//databusNotifyList(ctx)

// 	data, err := json.Marshal(ctx.Result())
// 	fmt.Println(string(data), err)
// }
// func TestNotifyEdit(t *testing.T) {
// 	flag.Set("conf", "../cmd/apm-admin-test.toml")
// 	conf.Init()
// 	apmSvc = service.New(conf.Conf)
// 	params := url.Values{}
// 	params.Set("n_id", "159")
// 	params.Set("offset", "new")
// 	params.Set("filter", "0")
// 	params.Set("state", "0")
// 	params.Set("concurrent", "10")
// 	params.Set("callback", `{"aaq":1}`)
// 	params.Set("filters", `[{"id":4,"field":"test","condition":1,"value":"vvv"}]`)
// 	req, _ := http.NewRequest("POST", path, strings.NewReader(params.Encode()))
// 	req.Header.Set("X-BACKEND-BILI-REAL-IP", "xxx")
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.ParseForm()
// 	ctx := bm.NewContext(context.TODO(), req, resp{}, time.Second)
// 	ctx.Set("username", "haoguanwei")
// 	ctx.Request().ParseForm()
// 	//databusNotifyEdit(ctx)
// 	data, err := json.Marshal(ctx.Result())
// 	fmt.Println(string(data), err)
// }

// func TestApplyProcess(t *testing.T) {
// 	flag.Set("conf", "../cmd/apm-admin-test.toml")
// 	conf.Init()
// 	apmSvc = service.New(conf.Conf)
// 	params := url.Values{}
// 	params.Set("id", "165")
// 	params.Set("state", "3")

// 	req, _ := http.NewRequest("POST", path, strings.NewReader(params.Encode()))
// 	req.Header.Set("X-BACKEND-BILI-REAL-IP", "xxx")
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.ParseForm()
// 	ctx := bm.NewContext(context.TODO(), req, resp{}, time.Second)
// 	ctx.Set("username", "haoguanwei")
// 	ctx.Request().ParseForm()
// 	//databusApplyApprovalProcess(ctx)
// 	fmt.Println(ctx.Result())
// }
