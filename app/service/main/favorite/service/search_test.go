package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Relations(t *testing.T) {
	Convey("Relations", t, func() {
		var (
			pn            = 1
			ps            = 30
			typ     int8  = 1
			mid     int64 = 88888894
			uid     int64 = 88888894
			fid     int64 = 59
			keyword       = "测试一下"
		)
		res, err := s.Relations(context.Background(), typ, mid, uid, fid, 0, 0, pn, ps, keyword, "")
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
