package service

import (
	"context"
	"testing"
	"time"

	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Generate(t *testing.T) {
	Convey("Generate 10 codes in batch", t, func() {
		mid := int64(88888970)
		num := int64(10)
		expireDay := int64(30)
		res, err := s.Generate(context.Background(), mid, num, expireDay)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, num)
	})
}

func TestService_List(t *testing.T) {
	Convey("List when account's invite codes is not empty", t, func() {
		mid := int64(88888970)
		now := time.Now().Unix()
		start, end := now-86400, now+86400
		res, err := s.List(context.Background(), mid, start, end)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	})
}

func TestService_ConcurrentGeneInviteCode(t *testing.T) {
	Convey("Generate 1000 codes in concurrency", t, func() {
		num := 1000
		mid := int64(88888970)
		ts := time.Now().Unix()
		m, err := concurrentGenerateCode(mid, ts, num, _geneSubCount)
		So(err, ShouldBeNil)
		So(len(m), ShouldEqual, num)
	})
}

func TestService_FetchMultiInfo(t *testing.T) {
	time.Sleep(time.Second * 2)
	Convey("Fetch multi info", t, func() {
		mids := []int64{88888970}
		Convey("when not timeout", func() {
			res, err := s.fetchInfos(context.Background(), mids, time.Second)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, len(mids))
		})
		Convey("when timeout", func() {
			_, err := s.fetchInfos(context.Background(), mids, time.Millisecond)
			So(err, ShouldEqual, ecode.Deadline.Error())
		})
	})
}
