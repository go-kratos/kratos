package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/admin/main/activity/model"

	xtime "go-common/library/time"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SubjectList(t *testing.T) {
	Convey("service test", t, WithService(func(s *Service) {
		p := &model.ListSub{
			Page:     1,
			PageSize: 15,
			Keyword:  "layang123",
			States:   []int{1},
			Types:    []int{18},
			Sctime:   1534835169,
			Ectime:   1546272001,
		}
		list, err := s.SubjectList(context.Background(), p)
		So(err, ShouldBeNil)
		for _, v := range list.List {
			fmt.Printf("%+v", v)
		}
	}))
}

func TestService_VideoList(t *testing.T) {
	Convey("service test", t, WithService(func(s *Service) {
		list, err := s.VideoList(context.Background())
		So(err, ShouldBeNil)
		for _, v := range list {
			fmt.Printf("%+v  %+v", v.ActSubject, v.Aids)
		}
	}))
}

func TestService_AddActSubject(t *testing.T) {
	Convey("service test", t, WithService(func(s *Service) {
		p := &model.AddList{
			ActSubject: model.ActSubject{
				Oid:        11,
				Type:       9,
				State:      1,
				Level:      5,
				Rank:       100,
				Stime:      xtime.Time(time.Now().Unix()),
				Etime:      xtime.Time(time.Now().Unix()),
				Ctime:      xtime.Time(time.Now().Unix()),
				Mtime:      xtime.Time(time.Now().Unix()),
				Lstime:     xtime.Time(time.Now().Unix()),
				Letime:     xtime.Time(time.Now().Unix()),
				Uetime:     xtime.Time(time.Now().Unix()),
				Ustime:     xtime.Time(time.Now().Unix()),
				Name:       "test one",
				Author:     "layang",
				ActURL:     "http://www.baidu.com/",
				Cover:      "cover",
				Flag:       128,
				Dic:        "dif",
				H5Cover:    "H5Cover",
				LikeLimit:  5,
				AndroidURL: "AndroidURL",
				IosURL:     "IosURL",
			},
			Protocol:  "Protocol",
			Types:     "1,2,3",
			Pubtime:   xtime.Time(time.Now().Unix()),
			Deltime:   xtime.Time(time.Now().Unix()),
			Editime:   xtime.Time(time.Now().Unix()),
			Tags:      "由三",
			Interval:  1,
			Tlimit:    123,
			Ltime:     124,
			Hot:       1,
			BgmID:     3,
			PasterID:  4,
			Oids:      "5,7,8",
			ScreenSet: 1,
		}
		res, err := s.AddActSubject(context.Background(), p)
		So(err, ShouldBeNil)
		fmt.Printf("%d", res)
	}))
}

func TestService_UpActSubject(t *testing.T) {
	Convey("service test", t, WithService(func(s *Service) {
		p := &model.AddList{
			ActSubject: model.ActSubject{
				Oid:        12,
				Type:       9,
				State:      0,
				Level:      6,
				Rank:       101,
				Stime:      xtime.Time(time.Now().Unix()),
				Etime:      xtime.Time(time.Now().Unix()),
				Ctime:      xtime.Time(time.Now().Unix()),
				Mtime:      xtime.Time(time.Now().Unix()),
				Lstime:     xtime.Time(time.Now().Unix()),
				Letime:     xtime.Time(time.Now().Unix()),
				Uetime:     xtime.Time(time.Now().Unix()),
				Ustime:     xtime.Time(time.Now().Unix()),
				Name:       "test two",
				Author:     "layang2",
				ActURL:     "http://www.baidu.com/2",
				Cover:      "cover2",
				Flag:       129,
				Dic:        "dif2",
				H5Cover:    "H5Cover2",
				LikeLimit:  6,
				AndroidURL: "AndroidURL2",
				IosURL:     "IosURL2",
			},
			Protocol:  "Protocol2",
			Types:     "1,2,3,4",
			Pubtime:   xtime.Time(time.Now().Unix()),
			Deltime:   xtime.Time(time.Now().Unix()),
			Editime:   xtime.Time(time.Now().Unix()),
			Tags:      "由三2",
			Interval:  2,
			Tlimit:    124,
			Ltime:     125,
			Hot:       0,
			BgmID:     4,
			PasterID:  8,
			Oids:      "5,7,8.9",
			ScreenSet: 2,
		}
		res, err := s.UpActSubject(context.Background(), p, 10298)
		So(err, ShouldBeNil)
		fmt.Printf("%+v", res)
	}))
}

func TestService_SubProtocol(t *testing.T) {
	Convey("sub protovol ", t, WithService(func(s *Service) {
		list, err := s.SubProtocol(context.Background(), 10256)
		So(err, ShouldBeNil)
		fmt.Printf("%+v", list)
	}))
}

func TestService_TimeConf(t *testing.T) {
	Convey("sub TimeConf ", t, WithService(func(s *Service) {
		list, err := s.TimeConf(context.Background(), 10298)
		So(err, ShouldBeNil)
		fmt.Printf("%+v", list)
	}))
}

func TestService_GetArticleMetas(t *testing.T) {
	Convey("sub TimeConf ", t, WithService(func(s *Service) {
		list, err := s.GetArticleMetas(context.Background(), []int64{1412})
		So(err, ShouldBeNil)
		fmt.Printf("%+v", list)
	}))
}
