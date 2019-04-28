package service

import (
	"fmt"
	"go-common/app/service/main/thumbup/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_updateCounts(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		messageID := int64(1)
		bid := int64(2)
		oid := int64(0)
		stat, err := s.dao.Stat(c, bid, oid, messageID)
		So(err, ShouldBeNil)
		args := [][2]int64{
			{1, 1},
			{-1000, -1000},
			{1, 0},
			{0, 1},
			{-1, 0},
			{-1, -1},
			{0, -1},
			{0, 0},
			{100, -10000},
			{100, 20},
			{-10, 20},
		}
		var l, d int64
		for _, x := range args {
			l = x[0]
			d = x[1]
			Convey(fmt.Sprintf("like %v dislike %v", l, d), func() {
				err := s.dao.UpdateCounts(c, bid, oid, messageID, l, d, 0)
				So(err, ShouldBeNil)
				nstat, err := s.dao.Stat(c, bid, oid, messageID)
				So(err, ShouldBeNil)
				likes := stat.Likes + l
				if likes < 0 {
					likes = 0
				}
				dislikes := stat.Dislikes + d
				if dislikes < 0 {
					dislikes = 0
				}
				So(nstat.Likes, ShouldEqual, likes)
				So(nstat.Dislikes, ShouldEqual, dislikes)
			})
		}
	}))
}

func Test_calculateCount(t *testing.T) {
	Convey("get data", t, func() {
		stat := model.Stats{Likes: 1, Dislikes: 1, ID: 1}
		So(calculateCount(stat, model.TypeLike), ShouldResemble, model.Stats{Likes: 2, Dislikes: 1, ID: 1})
		So(calculateCount(stat, model.TypeCancelLike), ShouldResemble, model.Stats{Likes: 0, Dislikes: 1, ID: 1})
		So(calculateCount(stat, model.TypeDislike), ShouldResemble, model.Stats{Likes: 1, Dislikes: 2, ID: 1})
		So(calculateCount(stat, model.TypeCancelDislike), ShouldResemble, model.Stats{Likes: 1, Dislikes: 0, ID: 1})
		So(calculateCount(stat, model.TypeLikeReverse), ShouldResemble, model.Stats{Likes: 0, Dislikes: 2, ID: 1})
		So(calculateCount(stat, model.TypeDislikeReverse), ShouldResemble, model.Stats{Likes: 2, Dislikes: 0, ID: 1})
		stat = model.Stats{Likes: 0, Dislikes: 0, ID: 1}
		So(calculateCount(stat, model.TypeCancelLike), ShouldResemble, model.Stats{Likes: -1, Dislikes: 0, ID: 1})
		So(calculateCount(stat, model.TypeCancelDislike), ShouldResemble, model.Stats{Likes: 0, Dislikes: -1, ID: 1})
	})
}
