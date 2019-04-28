package http

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/library/ecode"
	"net/url"
	"testing"
)

type DelCacheRes struct {
	Code int `json:"code"`
}

func queryDelCache(t *testing.T, uid int64) *DelCacheRes {
	params := url.Values{}
	params.Set("uid", fmt.Sprintf("%d", uid))
	req, _ := client.NewRequest("GET", _delCacheURL, "127.0.0.1", params)

	var res DelCacheRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func TestDelCache(t *testing.T) {
	once.Do(startHTTP)
	Convey("Del Cache", t, func() {
		var uid int64 = 1
		queryGet(t, uid, "pc")

		r := queryDelCache(t, uid)
		So(r.Code, ShouldEqual, 0)

		r = queryDelCache(t, uid)
		So(r.Code, ShouldEqual, ecode.NothingFound)
	})
}
