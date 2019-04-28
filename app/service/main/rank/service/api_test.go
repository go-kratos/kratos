package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/service/main/rank/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Sort(t *testing.T) {
	Convey("Sort", t, func() {
		s.rmap = make(map[int][]*model.Field)
		for i := 1; i < 80; i++ {
			f := new(model.Field)
			f.Flag = true
			f.Oid = int64(i)
			f.Pid = int16(i) % 10
			f.Click = i * i
			f.Pubtime = xtime.Time(time.Now().Unix())
			s.setField(int64(i), f)
		}
		fmt.Println(s.rmap)
		arg := &model.SortReq{
			Business: "archive",
			Field:    "click",
			Order:    "desc",
			Pn:       1,
			Ps:       30,
		}
		for i := 1; i < 60; i++ {
			arg.Oids = append(arg.Oids, int64(i))
		}
		res, err := s.Sort(context.Background(), arg)
		fmt.Println(res.Page)
		t.Logf("res:%+v", res)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func Test_Group(t *testing.T) {
	Convey("Group", t, func() {
		s.rmap = make(map[int][]*model.Field)
		for i := 1; i < 30; i++ {
			f := new(model.Field)
			f.Flag = true
			f.Oid = int64(i)
			f.Pid = int16(i) % 10
			f.Click = i * i
			f.Pubtime = xtime.Time(time.Now().Unix())
			s.setField(int64(i), f)
		}
		arg := &model.GroupReq{
			Business: "archive",
			Oids:     []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		}
		res, err := s.Group(context.Background(), arg)
		t.Logf("res:%+v", res)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
