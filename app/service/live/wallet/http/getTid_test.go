package http

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"net/url"
	"sync"
	"testing"
	"time"
)

type TidRes struct {
	Code int            `json:"code"`
	Resp *model.TidResp `json:"data"`
}

type TestTidParams struct {
	Biz  string
	Time int64
}

func getTestTidParams(biz string) *TestTidParams {
	return &TestTidParams{
		Biz:  biz,
		Time: time.Now().Unix(),
	}
}

func getTestRandServiceType() int32 {
	return r.Int31n(4)
}

func testWith(f func()) func() {
	once.Do(startHTTP)
	return f
}

func getTestParamsJson() string {
	params := getTestTidParams("gift")
	paramsBytes, _ := json.Marshal(params)
	paramsJson := string(paramsBytes[:])
	return paramsJson
}

func queryGetTid(t *testing.T, serviceType int32, tidParams string) *TidRes {
	params := url.Values{}
	params.Set("type", fmt.Sprintf("%d", serviceType))
	params.Set("params", tidParams)
	req, _ := client.NewRequest("POST", _getTidURL, "127.0.0.1", params)

	var res TidRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func TestGetTid(t *testing.T) {

	Convey("normal", t, testWith(func() {
		serviceType := getTestRandServiceType()
		paramsJson := getTestParamsJson()

		res := queryGetTid(t, serviceType, paramsJson)
		So(res.Code, ShouldEqual, 0)
		So(res.Resp.TransactionId, ShouldNotEqual, "")
	}))

	Convey("Twice Same params 同样的参数调用getTid　得到的tid应该不一样", t, testWith(func() {
		serviceType := getTestRandServiceType()
		paramsJson := getTestParamsJson()
		res := queryGetTid(t, serviceType, paramsJson)
		So(res.Code, ShouldEqual, 0)
		So(res.Resp.TransactionId, ShouldNotEqual, "")

		res1 := queryGetTid(t, serviceType, paramsJson)
		So(res1.Code, ShouldEqual, 0)
		So(res1.Resp.TransactionId, ShouldNotEqual, res.Resp.TransactionId)
	}))

	Convey("Test multi Same params", t, testWith(func() {
		serviceType := getTestRandServiceType()
		paramsJson := getTestParamsJson()

		resMap := make(map[int]*TidRes)
		wg := sync.WaitGroup{}

		var mutex sync.Mutex

		times := 10

		for i := 0; i < times; i++ {
			wg.Add(1)
			go func(i int) {
				mutex.Lock()
				resMap[i] = queryGetTid(t, serviceType, paramsJson)
				mutex.Unlock()
				wg.Done()
			}(i)
		}

		wg.Wait()

		tidMap := map[string]bool{}
		for _, res := range resMap {
			So(res.Code, ShouldEqual, 0)
			So(res.Resp.TransactionId, ShouldNotEqual, "")
			_, ok := tidMap[res.Resp.TransactionId]
			So(ok, ShouldEqual, false) // 应该找不到　因为同样的参数应该生成不同的tid
			tidMap[res.Resp.TransactionId] = true
		}

	}))

	Convey("params error", t, func() {
		Convey("serviceType 应该为　0 1 2 3 ", testWith(func() {
			wrongServiceTypes := []int32{-1, 5}
			for _, v := range wrongServiceTypes {
				paramsJson := getTestParamsJson()
				res := queryGetTid(t, v, paramsJson)
				So(res.Code, ShouldEqual, ecode.RequestErr.Code())
			}
		}))

		Convey("param不为空", testWith(func() {
			st := getTestRandServiceType()
			params := ""
			res := queryGetTid(t, st, params)
			So(res.Code, ShouldEqual, ecode.RequestErr.Code())
		}))
	})
}
