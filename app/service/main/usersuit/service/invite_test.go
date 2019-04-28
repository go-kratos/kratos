package service

import (
	"context"
	"testing"
	"time"

	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_BuyInvite(t *testing.T) {
	time.Sleep(time.Second * 2)
	Convey("buy invite", t, func() {
		Convey("when cookie is nil the return err should be NoLogin", func() {
			mid := int64(88888970)
			num := int64(2)
			_, err := s.BuyInvite(context.Background(), mid, num, "127.0.0.1")
			So(ecode.Cause(err), ShouldEqual, ecode.LackOfCoins)
		})
	})
}

func TestService_RangeMonth(t *testing.T) {
	const (
		_format = "2006-01-02 15:04:05"
	)
	type Range struct {
		start string
		end   string
	}
	m := map[string]*Range{
		"2017-12-01 00:00:00": {
			"2017-12-01 00:00:00",
			"2017-12-31 23:59:59",
		},
		"2017-12-24 12:05:00": {
			"2017-12-01 00:00:00",
			"2017-12-31 23:59:59",
		},
		"2017-12-24 23:59:59": {
			"2017-12-01 00:00:00",
			"2017-12-31 23:59:59",
		},
	}
	Convey("Range month", t, func() {
		for k, v := range m {
			t, err := time.Parse(_format, k)
			So(err, ShouldBeNil)
			start, end := rangeMonth(t)
			So(start.Format(_format), ShouldEqual, v.start)
			So(end.Format(_format), ShouldEqual, v.end)
		}
	})
}

func TestService_GeneInviteCode(t *testing.T) {
	Convey("Generate 10 codes in batch", t, func() {
		mid := int64(88888970)
		num := int64(1000)
		ts := time.Now().Unix()
		m := make(map[string]struct{})
		for i := int64(0); i < num; i++ {
			code := geneInviteCode(mid, ts)
			m[code] = struct{}{}
		}
		So(len(m), ShouldEqual, num)
	})
}
