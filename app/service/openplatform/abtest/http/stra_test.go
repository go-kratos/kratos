package http

import (
	"context"
	"net/url"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_getListURL      = "http://localhost:8801/openplatform/admin/abtest/list"
	_getVersionIDURL = "http://localhost:8801/openplatform/internal/abtest/versionid"
	_getVersionURL   = "http://localhost:8801/openplatform/internal/abtest/version"
	_AddURL          = "http://localhost:8801/openplatform/admin/abtest/add"
	_DelURL          = "http://localhost:8801/openplatform/admin/abtest/delete"
	_UpdateURL       = "http://localhost:8801/openplatform/admin/abtest/update"
	_UpdateStatusURL = "http://localhost:8801/openplatform/admin/abtest/status"
)

type TestData map[string]string
type Shoulds []interface{}

type TestCase struct {
	tag      string
	testData TestData
	should   Shoulds
}

var gvcs = []TestCase{
	TestCase{tag: "TestGetVersionID: valid parameters", testData: TestData{"group": "1"}, should: Shoulds{0}},
	TestCase{tag: "TestGetVersionID: empty parameters", testData: TestData{"group": ""}, should: Shoulds{-400}},
	TestCase{tag: "TestGetVersionID: invalid parameters", testData: TestData{"group": "asd"}, should: Shoulds{-400}},
	TestCase{tag: "TestGetVersionID: no parameters", testData: TestData{}, should: Shoulds{-400}},
}

func TestGetVersionID(t *testing.T) {
	for _, td := range gvcs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _getVersionIDURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}
			So(res.Code, ShouldEqual, td.should[0])
		})
	}
}

var guscs = []TestCase{
	TestCase{tag: "TestGetVersion: valid parameters", testData: TestData{"group": "1", "key": "23232", "version": "{}"}, should: Shoulds{0}},
	TestCase{tag: "TestGetVersion: no version", testData: TestData{"group": "1", "key": ""}, should: Shoulds{0}},
	TestCase{tag: "TestGetVersion: no key", testData: TestData{"group": "1", "version": "{}"}, should: Shoulds{0}},
	TestCase{tag: "TestGetVersion: no group", testData: TestData{"group": "1", "version": "{}"}, should: Shoulds{0}},
}

func TestGetVersion(t *testing.T) {
	for _, td := range guscs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _getVersionURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
				Data struct {
					V int         `json:"v"`
					D interface{} `json:"d"`
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

var glcs = []TestCase{
	TestCase{tag: "TestGetListAb: valid parameters", testData: TestData{"pn": "1", "ps": "20", "mstatus": "1,2,0"}, should: Shoulds{0}},
	TestCase{tag: "TestGetListAb: no pn", testData: TestData{"ps": "1", "mstatus": "1,2,0"}, should: Shoulds{-400}},
	TestCase{tag: "TestGetListAb: no ps", testData: TestData{"pn": "1", "mstatus": "1,2,0"}, should: Shoulds{-400}},
	TestCase{tag: "TestGetListAb: no mstatus", testData: TestData{"pn": "1", "ps": "10"}, should: Shoulds{-400}},
	TestCase{tag: "TestGetListAb: invalid pn", testData: TestData{"pn": "a", "ps": "20", "mstatus": "1,2,0"}, should: Shoulds{-400}},
	TestCase{tag: "TestGetListAb: invalid ps", testData: TestData{"pn": "1", "ps": "a", "mstatus": "1,2,0"}, should: Shoulds{-400}},
	TestCase{tag: "TestGetListAb: invalid mstatus", testData: TestData{"pn": "1", "ps": "20", "mstatus": "a"}, should: Shoulds{-400}},
}

func TestGetListAb(t *testing.T) {
	for _, td := range glcs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _getListURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
				Data struct {
					Result interface{} `json:"result"`
					Total  interface{} `json:"total"`
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

var testID int
var adcs = []TestCase{
	TestCase{tag: "TestAddAb: valid json", testData: TestData{"data": `{"name":"test1","desc":"test","stra":{"precision":100,"ratio":[20,80]},"result":1,"status":0,"group":1,"author":"test","modifer":"test"}`}, should: Shoulds{0}},
	TestCase{tag: "TestAddAb: no permission", testData: TestData{"data": `{"name":"test1","desc":"test","stra":{"precision":100,"ratio":[20,80]},"result":1,"status":0,"group":1,"author":"test","modifer":"test"}`, "group": "2"}, should: Shoulds{-400}},
	TestCase{tag: "TestAddAb: not json", testData: TestData{"data": `{"name":"test2","desc""test""stra":"{"precision":100,"ratio":[20,80]},"result":1,"status":0,"group":1,"author":"test","modifer":"test"}`}, should: Shoulds{-400}},
	TestCase{tag: "TestAddAb: invalid stra", testData: TestData{"data": `{"name":"test3","desc":"test","stra":{"precision":100,"ratio":[20,70]},"result":1,"status":0,"group":1,"author":"test","modifer":"test"}`}, should: Shoulds{-400}},
	TestCase{tag: "TestAddAb: no data", testData: TestData{}, should: Shoulds{-400}},
}

func TestAddAb(t *testing.T) {
	for _, td := range adcs {
		Convey(td.tag, t, func() {
			var res struct {
				Code int `json:"code"`
				Data struct {
					Newid int `json:"newid"`
				} `json:"data"`
			}

			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _AddURL, "127.0.0.1", params)
			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}

			So(res.Code, ShouldEqual, td.should[0])
			if res.Code == 0 {
				testID = res.Data.Newid
			}
		})
	}
}

func TestUpdateAb(t *testing.T) {
	var upcs = []TestCase{
		TestCase{tag: "TestUpdateAb: valid params", testData: TestData{"id": strconv.Itoa(testID), "data": `{"name":"test","desc":"update","stra":{"precision":100,"ratio":[20,80]},"result":1,"status":0,"group":1,"author":"test","modifer":"test"}`}, should: Shoulds{0}},
		TestCase{tag: "TestUpdateAb: no permission", testData: TestData{"id": strconv.Itoa(testID), "data": `{"name":"test1","desc":"test","stra":{"precision":100,"ratio":[20,80]},"result":1,"status":0,"group":2,"author":"test","modifer":"test"}`, "group": "2"}, should: Shoulds{-500}},
		TestCase{tag: "TestUpdateAb: invalid params", testData: TestData{}, should: Shoulds{-400}},
		TestCase{tag: "TestUpdateAb: invalid id", testData: TestData{"id": "11111", "data": "aa"}, should: Shoulds{-500}},
		TestCase{tag: "TestUpdateAb: invalid data", testData: TestData{"id": strconv.Itoa(testID), "data": "aa"}, should: Shoulds{-400}},
		TestCase{tag: "TestUpdateAb: valid stra", testData: TestData{"id": strconv.Itoa(testID), "data": `{"name":"test","desc":"update","stra":{"precision":100,"ratio":[20,81]},"result":1,"status":0,"group":1,"author":"test","modifer":"test"}`}, should: Shoulds{-400}},
		TestCase{tag: "TestUpdateAb: valid params group 0", testData: TestData{"id": strconv.Itoa(testID), "data": `{"name":"test","desc":"update2","stra":{"precision":100,"ratio":[20,80]},"result":1,"status":0,"group":0,"author":"test","modifer":"test"}`}, should: Shoulds{0}},
	}
	for _, td := range upcs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _UpdateURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}
			So(res.Code, ShouldEqual, td.should[0])
		})
	}
}

func TestUpdateStatusAb(t *testing.T) {
	var upscs = []TestCase{
		TestCase{tag: "TestUpdateStatusAb: valid params", testData: TestData{"id": strconv.Itoa(testID), "status": "1", "modifier": "test2"}, should: Shoulds{0}},
		TestCase{tag: "TestUpdateStatusAb: no permission", testData: TestData{"id": strconv.Itoa(testID), "status": "1", "modifier": "test2", "group": "2"}, should: Shoulds{-500}},
		TestCase{tag: "TestUpdateStatusAb: invalid params", testData: TestData{}, should: Shoulds{-400}},
		TestCase{tag: "TestUpdateStatusAb: invalid id", testData: TestData{"id": "11111", "data": "aa"}, should: Shoulds{-400}},
		TestCase{tag: "TestUpdateStatusAb: invalid status", testData: TestData{"id": strconv.Itoa(testID), "status": "4", "modifier": "test2"}, should: Shoulds{-400}},
		TestCase{tag: "TestUpdateStatusAb: valid params", testData: TestData{"id": strconv.Itoa(testID), "status": "3", "modifier": "test2"}, should: Shoulds{0}},
	}
	for _, td := range upscs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _UpdateStatusURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}
			So(res.Code, ShouldEqual, td.should[0])
		})
	}
}

func TestDelAb(t *testing.T) {
	var dacs = []TestCase{
		TestCase{tag: "TestDelAb: no permission", testData: TestData{"id": strconv.Itoa(testID), "group": "2"}, should: Shoulds{-500}},
		TestCase{tag: "TestDelAb: valid id", testData: TestData{"id": strconv.Itoa(testID)}, should: Shoulds{0}},
		TestCase{tag: "TestDelAb: invalid id", testData: TestData{"id": "x"}, should: Shoulds{-400}},
	}
	for _, td := range dacs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _DelURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}
			So(res.Code, ShouldEqual, td.should[0])
		})
	}
}
