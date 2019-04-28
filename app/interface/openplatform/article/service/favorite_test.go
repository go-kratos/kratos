package service

// import (
// 	"context"
// 	"testing"

// 	. "github.com/smartystreets/goconvey/convey"
// )

// func Test_Favorite(t *testing.T) {
// 	var (
// 		aid = int64(9999999)
// 		mid = int64(88888929)
// 		fid = int64(0)
// 		pn  = 1
// 		ps  = 10
// 	)
// 	Convey("AddFavorite", t, WithService(func(s *Service) {
// 		err := s.AddFavorite(context.TODO(), mid, aid, fid, "")
// 		// 11201 means you have added it before.
// 		So(err, ShouldBeNil)

// 		Convey("Favorites", func() {
// 			res, page, err := s.Favs(context.TODO(), mid, fid, pn, ps, "")
// 			So(err, ShouldBeNil)
// 			So(res, ShouldNotBeEmpty)
// 			So(page, ShouldNotBeEmpty)
// 			// t.Logf("result: %+v", res)
// 			// t.Logf("page: %+v", page)

// 			Convey("DelFavorite", func() {
// 				err := s.DelFavorite(context.TODO(), mid, aid, fid, "")
// 				So(err, ShouldBeNil)
// 			})
// 		})
// 	}))
// }
