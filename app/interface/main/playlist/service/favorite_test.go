package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/playlist/conf"
	"go-common/app/interface/main/playlist/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Add(t *testing.T) {
	var (
		mid          = int64(88888929)
		plPublic     = int8(0)
		plName       = "播单名称1"
		plDescripton = "播单描述1"
		cover        = "http://image1.jpg"
	)
	Convey("Add", t, WithService(func(s *Service) {
		res, err := s.Add(context.Background(), mid, plPublic, plName, plDescripton, cover, "", "")
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
	}))
}

func TestService_Del(t *testing.T) {
	var (
		mid = int64(88888929)
		pid = int64(1)
	)
	Convey("Add", t, WithService(func(s *Service) {
		err := s.Del(context.Background(), mid, pid)
		So(err, ShouldBeNil)
	}))
}

func TestService_Update(t *testing.T) {
	var (
		mid          = int64(88888929)
		pid          = int64(1)
		plPublic     = int8(0)
		plName       = "播单名称1"
		plDescripton = "播单描述1"
		cover        = "http://image1.jpg"
	)
	Convey("Add", t, WithService(func(s *Service) {
		err := s.Update(context.Background(), mid, pid, plPublic, plName, plDescripton, cover, "", "")
		So(err, ShouldBeNil)
	}))
}

func TestService_AddVideo(t *testing.T) {
	var (
		mid  = int64(88888929)
		pid  = int64(1)
		aids = []int64{11, 22, 33, 44, 55}
	)
	Convey("Add", t, WithService(func(s *Service) {
		videos, err := s.AddVideo(context.Background(), mid, pid, aids)
		So(videos, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}

func TestService_DelVideo(t *testing.T) {
	var (
		mid  = int64(88888929)
		pid  = int64(1)
		aids = []int64{11, 22, 33, 44, 55}
	)
	Convey("Add", t, WithService(func(s *Service) {
		err := s.DelVideo(context.Background(), mid, pid, aids)
		So(err, ShouldBeNil)
	}))
}

func TestService_SortVideo(t *testing.T) {
	Convey("sort", t, func() {
		var (
			arcs                               []*model.ArcSort
			aidSort, preSort, afSort, orderNum int64
			start, end                         int
			top, bottom                        bool
		)
		aid := int64(6)
		sort := int64(1)
		arcs = []*model.ArcSort{
			{Aid: 1, Sort: 100},
			{Aid: 2, Sort: 200},
			{Aid: 3, Sort: 300},
			{Aid: 4, Sort: 400},
			{Aid: 5, Sort: 500},
			{Aid: 6, Sort: 600},
		}
		if sort == _first {
			top = true
		} else if sort == int64(len(arcs)) {
			bottom = true
		}
		for k, v := range arcs {
			if k == 0 && top {
				afSort = v.Sort
			}
			if k == len(arcs)-1 && bottom {
				preSort = v.Sort
			}
			if aid == v.Aid {
				if sort == int64(k+1) {
					return
				}
				aidSort = v.Sort
			}
			if sort == int64(k+1) {
				if !top && !bottom {
					if aidSort > sort {
						preSort = arcs[k].Sort
						afSort = arcs[k+1].Sort
					} else {
						preSort = arcs[k-1].Sort
						afSort = arcs[k].Sort
					}
				}
			}
		}
		if top {
			println("top")
			orderNum = afSort / 2
		} else if bottom {
			println("bottom")
			orderNum = preSort + int64(conf.Conf.Rule.SortStep)
		} else {
			println("else")
			orderNum = preSort + (afSort-preSort)/2
		}
		println(start, end, preSort, afSort)
		for _, v := range arcs {
			if v.Aid == aid {
				v.Sort = orderNum
			}
		}
		for _, v := range arcs {
			Printf("%+v", v)
		}

	})
}
