package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SimpleArchive(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("SimpleArchive", t, WithService(func(s *Service) {
		data, err := svr.SimpleArchive(c, 222, 1)
		fmt.Println(data)
		So(err, ShouldBeNil)
	}))
}

func TestService_SimpleVideos(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("SimpleVideos", t, WithService(func(s *Service) {
		_, err := svr.SimpleVideos(c, 1)
		So(err, ShouldNotBeNil)
	}))
}

func TestService_Archive(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("Archive", t, WithService(func(s *Service) {
		_, err := svr.Archive(c, 222)
		So(err, ShouldBeNil)
	}))
}

func TestService_ArchivePOI(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("Archive poi", t, WithService(func(s *Service) {
		_, err := svr.ArchivePOI(c, 222)
		So(err, ShouldBeNil)
	}))
}

func TestService_Archives(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("Archives", t, WithService(func(s *Service) {
		_, err := svr.Archives(c, []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func TestService_Flows(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("Flows", t, WithService(func(s *Service) {
		data := svr.Flows(c)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_UpsForbid(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("UpsForbid", t, WithService(func(s *Service) {
		data := svr.UpsForbid(c)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_ArcHistorys(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("ArcHistorys", t, WithService(func(s *Service) {
		data := svr.ArcHistorys(c, 1)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_AppFeedAids(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("AppFeedAids", t, WithService(func(s *Service) {
		data, _ := svr.AppFeedAids(c)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_DescFormats(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("DescFormats", t, WithService(func(s *Service) {
		data, _ := svr.DescFormats(c)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_VideoJamLevel(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("VideoJamLevel", t, WithService(func(s *Service) {
		data, _ := svr.VideoJamLevel(c)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_Recos(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("Recos", t, WithService(func(s *Service) {
		_, err := svr.Recos(c, 1)
		So(err, ShouldBeNil)
	}))
}

func Test_RejectedArchives(t *testing.T) {
	var (
		mid      int64 = 2089809
		state    int32 = -4
		pn       int32
		ps       int32
		start, _ = time.Parse("20060102", "20100101")
	)
	Convey("RejectedArchives", t, func() {
		as, count, err := svr.RejectedArchives(context.TODO(), mid, state, pn, ps, &start)
		So(err, ShouldBeNil)
		for _, a := range as {
			ShouldEqual(a.Mid, mid)
			ShouldNotBeEmpty(a.RejectReason)
			ShouldNotEqual(count, 0)
		}
	})
}
