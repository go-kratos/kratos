package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type result struct {
	Code    int    `gorm:"column:code" json:"code"`
	Message string `gorm:"column:message" json:"message"`
}

//http client
func httpDo(method, uri, cookie string, params url.Values) (data *result, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, uri, strings.NewReader(params.Encode()))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data = &result{}
	json.Unmarshal(body, data)
	return
}

func TestHttpDo(t *testing.T) {
	uri := "http://127.0.0.1:7331/x/admin/apm/canal/apply/approval/process"
	params := url.Values{"id": {"20"}, "state": {"3"}}
	cookie := "username=hedan; _AJSESSIONID=d9dde52c35e2e0389fd8c345873c7d22; sven-apm=429fca48067ef8234335cddda8a7d64ba46662e3d9448d38df6184f63a37e165"

	t.Error("print log")
	if res, err := httpDo("POST", uri, cookie, params); err != nil {
		fmt.Println("There are some error~")
		fmt.Printf("%s", err)
	} else {
		fmt.Println("code:", res.Code)
		fmt.Println("message:", res.Message)
	}
}

func TestProcess(t *testing.T) {

	Convey("if apply.State == 1, state =[2,3,4]; expect: code ==0", t, func() {
		uri := "http://127.0.0.1:7331/x/admin/apm/canal/apply/approval/process"
		cookie := "username=hedan; _AJSESSIONID=d9dde52c35e2e0389fd8c345873c7d22; sven-apm=f78e8fadc7f41048258acf43f5e9abee59f7b7a5bb6312b81fd31a21945fd21e"
		params := url.Values{"id": {"7"}, "state": {"2"}}
		res, err := httpDo("POST", uri, cookie, params)
		if err != nil {
			t.Error(err)
		} else {
			So(res.Code, ShouldEqual, 0)
		}
	})
}

func TestProcess2(t *testing.T) {

	Convey("if apply.State == 2, state =[2,3,4]; expect: code ==0", t, func() {
		uri := "http://127.0.0.1:7331/x/admin/apm/canal/apply/approval/process"
		cookie := "username=hedan; _AJSESSIONID=d9dde52c35e2e0389fd8c345873c7d22; sven-apm=f78e8fadc7f41048258acf43f5e9abee59f7b7a5bb6312b81fd31a21945fd21e"
		params := url.Values{"id": {"8"}, "state": {"3"}}
		res, err := httpDo("POST", uri, cookie, params)
		if err != nil {
			t.Error(err)
		} else {
			So(res.Code, ShouldEqual, 0)
		}
	})
}

func TestProcess3(t *testing.T) {

	Convey("if apply.State == 3, state =[2,3,4]; expect: code == -400", t, func() {
		uri := "http://127.0.0.1:7331/x/admin/apm/canal/apply/approval/process"
		cookie := "username=hedan; _AJSESSIONID=d9dde52c35e2e0389fd8c345873c7d22; sven-apm=f78e8fadc7f41048258acf43f5e9abee59f7b7a5bb6312b81fd31a21945fd21e"
		params := url.Values{"id": {"9"}, "state": {"2"}}
		res, err := httpDo("POST", uri, cookie, params)
		if err != nil {
			t.Error(err)
		} else {
			So(res.Code, ShouldEqual, -400)
			So(res.Message, ShouldEqual, "只有申请中和打回才可审核")
		}
	})
}

func TestProcess4(t *testing.T) {

	Convey("if apply.State == 4, state =[2,3,4]; expect: code == -400", t, func() {
		uri := "http://127.0.0.1:7331/x/admin/apm/canal/apply/approval/process"
		cookie := "username=hedan; _AJSESSIONID=d9dde52c35e2e0389fd8c345873c7d22; sven-apm=f78e8fadc7f41048258acf43f5e9abee59f7b7a5bb6312b81fd31a21945fd21e"
		params := url.Values{"id": {"10"}, "state": {"2"}}
		res, err := httpDo("POST", uri, cookie, params)
		if err != nil {
			t.Error(err)
		} else {
			So(res.Code, ShouldEqual, -400)
			So(res.Message, ShouldEqual, "只有申请中和打回才可审核")
		}
	})
}

func TestProcess5(t *testing.T) {

	Convey("if apply.State == 1, state !=[2,3,4]; expect: code == -400", t, func() {
		uri := "http://127.0.0.1:7331/x/admin/apm/canal/apply/approval/process"
		cookie := "username=hedan; _AJSESSIONID=d9dde52c35e2e0389fd8c345873c7d22; sven-apm=f78e8fadc7f41048258acf43f5e9abee59f7b7a5bb6312b81fd31a21945fd21e"
		params := url.Values{"id": {"11"}, "state": {"5"}}
		res, err := httpDo("POST", uri, cookie, params)
		if err != nil {
			t.Error(err)
		} else {
			So(res.Code, ShouldEqual, -400)
			So(res.Message, ShouldEqual, "state值范围2,3,4")
		}
	})
}
