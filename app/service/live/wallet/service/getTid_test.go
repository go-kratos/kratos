package service

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"sync"
	"testing"
	"time"
)

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

type ServiceResForTest struct {
	Err error
	V   interface{}
}

var tidLen = 36

func getTestParamsJson() string {
	params := getTestTidParams("test-service")
	paramsBytes, _ := json.Marshal(params)
	paramsJson := string(paramsBytes[:])
	return paramsJson
}

func getTestRandServiceType() int32 {
	return r.Int31n(4)
}

func TestService_GetTid(t *testing.T) {
	Convey("normal", t, testWith(func() {
		for i := 0; i < 4; i++ {
			v, err := s.GetTid(ctx, getTestDefaultBasicParam(""), 0, i, getTestParamsJson())
			So(v, ShouldNotBeNil)
			So(err, ShouldBeNil)

			resp := v.(*model.TidResp)
			So(len(resp.TransactionId), ShouldEqual, tidLen)
		}
	}))

	Convey("same params twice", t, testWith(func() {
		var uid int64 = 0
		st := getTestRandServiceType()
		p := getTestDefaultBasicParam("")
		qp := getTestParamsJson()

		v, err := s.GetTid(ctx, p, uid, st, qp)
		So(v, ShouldNotBeNil)
		So(err, ShouldBeNil)
		resp1 := v.(*model.TidResp)
		So(len(resp1.TransactionId), ShouldEqual, tidLen)

		v, err = s.GetTid(ctx, p, uid, st, qp)
		So(v, ShouldNotBeNil)
		So(err, ShouldBeNil)
		resp2 := v.(*model.TidResp)
		So(len(resp2.TransactionId), ShouldEqual, tidLen)
		So(resp1.TransactionId, ShouldNotEqual, resp2.TransactionId)

	}))

	Convey("params", t, testWith(func() {
		Convey("invalid service Type", func() {
			invalidService := []int32{-1, -2, 4, 5}
			for _, st := range invalidService {
				v, err := s.GetTid(ctx, getTestDefaultBasicParam(""), 0, st, getTestParamsJson())
				So(err, ShouldEqual, ecode.RequestErr)
				So(v, ShouldBeNil)

			}
		})

		Convey("invalid query params", func() {
			st := 0
			params := ""
			v, err := s.GetTid(ctx, getTestDefaultBasicParam(""), 0, st, params)

			So(err, ShouldEqual, ecode.RequestErr)
			So(v, ShouldBeNil)

		})

	}))

	Convey("multi", t, testWith(func() {
		var uid int64 = 0
		st := getTestRandServiceType()
		p := getTestDefaultBasicParam("")
		qp := getTestParamsJson()

		var wg sync.WaitGroup
		times := 10

		tidResMap := make(map[int]*ServiceResForTest)
		var lock sync.Mutex
		for i := 0; i < times; i++ {
			wg.Add(1)
			go func(index int) {
				localService := New(conf.Conf)
				v, err := localService.GetTid(ctx, p, uid, st, qp)
				lock.Lock()
				tidResMap[index] = &ServiceResForTest{V: v, Err: err}
				lock.Unlock()
				wg.Done()
			}(i)
		}
		wg.Wait()

		tidMap := make(map[string]bool)

		for _, item := range tidResMap {
			So(item.Err, ShouldBeNil)
			resp := item.V.(*model.TidResp)
			So(len(resp.TransactionId), ShouldEqual, tidLen)
			_, ok := tidMap[resp.TransactionId]
			So(ok, ShouldBeFalse)
			tidMap[resp.TransactionId] = true
		}

		So(len(tidMap), ShouldEqual, times)

	}))

}
