package http

import (
	"context"
	"net/url"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_getGroupListURL = "http://localhost:8801/openplatform/admin/abtest/group/list"
	_AddGroupURL     = "http://localhost:8801/openplatform/admin/abtest/group/add"
	_DelGroupURL     = "http://localhost:8801/openplatform/admin/abtest/group/delete"
	_UpdateGroupURL  = "http://localhost:8801/openplatform/admin/abtest/group/update"
)

var agcs = []TestCase{
	TestCase{tag: "TestAddGroup: valid parameters", testData: TestData{"name": "test", "desc": "test add"}, should: Shoulds{0}},
	TestCase{tag: "TestAddGroup: empty parameters", testData: TestData{"name": "", "desc": ""}, should: Shoulds{-400}},
	TestCase{tag: "TestAddGroup: no parameters", testData: TestData{}, should: Shoulds{-400}},
}

func TestAddGroup(t *testing.T) {
	for _, td := range agcs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _AddGroupURL, "127.0.0.1", params)
			var res struct {
				Code int `json:"code"`
				Data int `json:"data"`
			}

			if err := client.Do(context.TODO(), req, &res); err != nil {
				t.Errorf("client.Do() error(%v)", err)
				t.FailNow()
			}
			So(res.Code, ShouldEqual, td.should[0])
			if res.Code == 0 {
				testID = res.Data
			}
		})
	}
}

func TestUpdateGroup(t *testing.T) {
	var ugcs = []TestCase{
		TestCase{tag: "TestUpdateGroup: valid parameters", testData: TestData{"id": strconv.Itoa(testID), "name": "test", "desc": "test update"}, should: Shoulds{0}},
		TestCase{tag: "TestUpdateGroup: empty parameters", testData: TestData{"id": "0", "name": "", "desc": ""}, should: Shoulds{-400}},
		TestCase{tag: "TestUpdateGroup: no parameters", testData: TestData{}, should: Shoulds{-400}},
	}

	for _, td := range ugcs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _UpdateGroupURL, "127.0.0.1", params)
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

func TestDeleteGroup(t *testing.T) {
	var dgcs = []TestCase{
		TestCase{tag: "TestDeleteGroup: valid parameters", testData: TestData{"id": strconv.Itoa(testID), "name": "test", "desc": "test delete"}, should: Shoulds{0}},
		TestCase{tag: "TestDeleteGroup: empty parameters", testData: TestData{"id": "0", "name": "", "desc": ""}, should: Shoulds{-400}},
		TestCase{tag: "TestDeleteGroup: no parameters", testData: TestData{}, should: Shoulds{-400}},
	}

	for _, td := range dgcs {
		Convey(td.tag, t, func() {
			params := url.Values{}
			for k, v := range td.testData {
				params.Set(k, v)
			}
			req, _ := client.NewRequest("GET", _DelGroupURL, "127.0.0.1", params)
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

func TestListGroup(t *testing.T) {
	Convey("TestListGroup: ", t, func() {
		params := url.Values{}
		req, _ := client.NewRequest("GET", _getGroupListURL, "127.0.0.1", params)
		var res struct {
			Code int `json:"code"`
		}

		if err := client.Do(context.TODO(), req, &res); err != nil {
			t.Errorf("client.Do() error(%v)", err)
			t.FailNow()
		}
		So(res.Code, ShouldEqual, 0)
	})
}
