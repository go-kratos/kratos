package http

import (
	"context"
	"net/url"
	"testing"

	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_qusBankInfoURL = "http://localhost:8801/openplatform/internal/antifraud/qusb/info?qbid=100"
	_qusBanklistURL = "http://localhost:8801/openplatform/internal/antifraud/qusb/list"

	_qslistURL = "http://localhost:8801/openplatform/internal/antifraud/qs/list"
	_qsInfoURL = "http://localhost:8801/openplatform/internal/antifraud/qs/info"
	_qsGetURL  = "http://localhost:8801/openplatform/internal/antifraud/qs/get"
)

type TestData map[string]string
type Shoulds []interface{}

type TestCase struct {
	tag      string
	testData TestData
	should   Shoulds
}

var glcs = []TestCase{
	{tag: "TestQusBankList: valid parameters", testData: TestData{"page": "1", "page_size": "20"}, should: Shoulds{-0}},
	{tag: "TestQusBankList: no page", testData: TestData{"page_size": "1"}, should: Shoulds{-400}},
	{tag: "TestQusBankList: no page_size", testData: TestData{"page": "1"}, should: Shoulds{-400}},
	{tag: "TestQusBankList: no mstatus", testData: TestData{"page": "a", "page_size": "b"}, should: Shoulds{-400}},
	{tag: "TestQusBankList: invalid page", testData: TestData{"page": "a", "page_size": "20"}, should: Shoulds{-400}},
	{tag: "TestQusBankList: invalid page_size", testData: TestData{"page": "1", "page_size": "a"}, should: Shoulds{-400}},
}

func TestQusBankList(t *testing.T) {
	for _, td := range glcs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _qusBanklistURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
				Data struct {
					Result   interface{} `json:"result"`
					Total    interface{} `json:"total"`
					PageNo   interface{} `json:"page_no"`
					PageSize interface{} `json:"page_size"`
					Items    interface{} `json:"items"`
				} `json:"data"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}
			So(res.Code, ShouldEqual, td.should[0])
		})
	}
}

func TestQusList(t *testing.T) {
	for _, td := range glcs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _qslistURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
				Data struct {
					Result   interface{} `json:"result"`
					Total    interface{} `json:"total"`
					PageNo   interface{} `json:"page_no"`
					PageSize interface{} `json:"page_size"`
					Items    interface{} `json:"items"`
				} `json:"data"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}
			So(res.Code, ShouldEqual, td.should[0])
		})
	}
}

var argsBankInfo = []TestCase{
	{tag: "TestQusBankInfo: valid parameters", testData: TestData{"qb_id": "1111"}, should: Shoulds{0}},
	{tag: "TestQusBankInfo: no qb_id", testData: TestData{"qb_id": "1"}, should: Shoulds{0}},
	{tag: "TestQusBankInfo: invalid qb_id", testData: TestData{"qb_id": "a"}, should: Shoulds{-400}},
}

func TestQusBankInfo(t *testing.T) {
	for _, td := range argsBankInfo {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _qusBankInfoURL, "127.0.0.1", params)
			var res struct {
				Code int                `json:"code"`
				Data model.QuestionBank `json:"data"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}

			So(res.Code, ShouldEqual, td.should[0])
		})
	}
}

var argsQusInfo = []TestCase{
	{tag: "TestQusBankInfo: valid parameters", testData: TestData{"qid": "1111"}, should: Shoulds{20001005}},
	{tag: "TestQusBankInfo: invalid qid", testData: TestData{"qid": "a"}, should: Shoulds{-400}},
}

func TestQusInfo(t *testing.T) {
	for _, td := range argsQusInfo {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _qsInfoURL, "127.0.0.1", params)
			var res struct {
				Code int                `json:"code"`
				Data model.QuestionBank `json:"data"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}
			So(res.Code, ShouldEqual, td.should[0])
		})
	}
}

var argsGetQuestion = []TestCase{
	{tag: "TestQusBankInfo: valid parameters", testData: TestData{"uid": "1111", "target_item": "11111", "target_item_type": "1", "source": "1", "platform": "1", "component_id": "122"},
		should: Shoulds{ecode.BindBankNotFound.Code(), ecode.GetComponentIDErr.Code(), ecode.SetComponentIDErr.Code(), 0}},
	{tag: "TestQusBankInfo: invalid ", testData: TestData{"uid": "a"}, should: Shoulds{-400, -400}},
}

func TestGetQuestion(t *testing.T) {
	for _, td := range argsGetQuestion {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _qsGetURL, "127.0.0.1", params)
			var res struct {
				Code int                `json:"code"`
				Data model.QuestionBank `json:"data"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}

			So(res.Code, ShouldBeIn, td.should...)
		})
	}
}
