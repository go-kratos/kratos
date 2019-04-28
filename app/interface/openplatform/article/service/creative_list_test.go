package service

import (
	"context"
	"go-common/app/interface/openplatform/article/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_CreativeListAllArticles(t *testing.T) {
	Convey("get articles", t, WithCleanCache(func() {
		list, res, err := s.CreativeListAllArticles(context.TODO(), 100, 5)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeEmpty)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_CreativeCanAddArticles(t *testing.T) {
	Convey("get articles", t, WithCleanCache(func() {
		res, err := s.CreativeCanAddArticles(context.TODO(), 88888929)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_calculateListArtPosition(t *testing.T) {
	a := &model.ListArtMeta{ID: 1, Position: 1000}
	b := &model.ListArtMeta{ID: 2, Position: 2000}
	c := &model.ListArtMeta{ID: 3, Position: 3000}
	d := &model.ListArtMeta{ID: 4, Position: 4000}

	// a b c d  =>  d a b c  update d 移动1
	// a b c d  =>  c d b a update a b  移动2
	// a b c    =>  a b c d  update d  末尾增加1
	// a b c    =>  a d b c update d  中间增加1
	// a b c    =>  d a b c 前面增加1
	// []       =>  a b 新增加2
	// a b c d  =>  a b remove 2
	// 可以考虑对删一 加一 移1 进行特别优化 其他全修改

	Convey("a b c  =>  a b c d", t, func() {
		update, delete := calculateListArtPosition([]*model.ListArtMeta{a, b, c}, []int64{a.ID, b.ID, c.ID, d.ID})
		So(delete, ShouldBeNil)
		So(update, ShouldResemble, []*model.ListArtMeta{
			&model.ListArtMeta{ID: d.ID, Position: 4000},
		})
	})

	Convey("[]  =>  a b", t, func() {
		update, delete := calculateListArtPosition(nil, []int64{a.ID, b.ID})
		So(delete, ShouldBeNil)
		So(update, ShouldResemble, []*model.ListArtMeta{a, b})
	})

	Convey("a b c d  =>  a b c ", t, func() {
		update, delete := calculateListArtPosition([]*model.ListArtMeta{a, b, c, d}, []int64{a.ID, b.ID, c.ID})
		So(delete, ShouldResemble, []int64{d.ID})
		So(update, ShouldBeNil)
	})

	Convey("a b c d  =>  d a b c", t, func() {
		update, delete := calculateListArtPosition([]*model.ListArtMeta{a, b, c, d}, []int64{d.ID, a.ID, b.ID, c.ID})
		So(delete, ShouldBeNil)
		So(update, ShouldResemble, []*model.ListArtMeta{
			&model.ListArtMeta{ID: d.ID, Position: 1000},
			&model.ListArtMeta{ID: a.ID, Position: 2000},
			&model.ListArtMeta{ID: b.ID, Position: 3000},
			&model.ListArtMeta{ID: c.ID, Position: 4000},
		})
	})
}

func Test_CreativeUpdateListArticles(t *testing.T) {
	Convey("update articles", t, WithCleanCache(func() {
		mid := int64(88888929)
		aids := []int64{929, 830, 1000}
		list, err := s.CreativeUpdateListArticles(context.TODO(), 8, "newName", "", "summary", false, mid, aids)
		So(err, ShouldBeNil)
		list, _ = s.dao.RawList(context.TODO(), 8)
		metas, _ := s.dao.CreativeListArticles(context.TODO(), 8)
		So(list.Name, ShouldEqual, "newName")
		So(len(metas), ShouldEqual, 2)
		So(metas[0].ID, ShouldEqual, 929)
		So(metas[1].ID, ShouldEqual, 830)
	}))
}

func Test_CreativeUpdateArticleList(t *testing.T) {
	c := context.TODO()
	aid := int64(821)
	listid := int64(8)
	mid := int64(88888929)
	Convey("set list", t, WithCleanCache(func() {
		err := s.CreativeUpdateArticleList(c, mid, aid, listid, false)
		So(err, ShouldBeNil)
		arts, err := s.dao.RawArtsListID(c, []int64{aid})
		So(err, ShouldBeNil)
		So(arts[aid], ShouldEqual, listid)
		// time.Sleep(time.Second)
		// r, err := s.dao.ArticlesListCache(c, []int64{aid})
		// So(err, ShouldBeNil)
		// So(r[aid], ShouldEqual, listid)
		// r2, err := s.dao.ListArtsCacheMap(c, listid)
		// So(err, ShouldBeNil)
		// So(r2[aid], ShouldNotBeEmpty)
		Convey("remove", func() {
			err := s.CreativeUpdateArticleList(c, mid, aid, 0, false)
			So(err, ShouldBeNil)
			arts, err := s.dao.RawArtsListID(c, []int64{aid})
			So(err, ShouldBeNil)
			So(arts[aid], ShouldEqual, 0)
		})
		Convey("update", func() {
			err := s.CreativeUpdateArticleList(c, mid, aid, 17, false)
			So(err, ShouldBeNil)
			arts, err := s.dao.RawArtsListID(c, []int64{aid})
			So(err, ShouldBeNil)
			So(arts[aid], ShouldEqual, 17)
		})
	}))
}
