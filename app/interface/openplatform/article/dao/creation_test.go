package dao

import (
	"testing"

	"go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Articles(t *testing.T) {
	var (
		c   = ctx()
		aid int64
		art = model.Article{
			Meta: &model.Meta{
				ID:              0,
				Title:           "1",
				Summary:         "2",
				BannerURL:       "https://i0.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg",
				TemplateID:      1,
				State:           0,
				Category:        &model.Category{ID: 1},
				Author:          &model.Author{Mid: 123},
				Reprint:         0,
				ImageURLs:       []string{"https://i0.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "https://i0.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "https://i0.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg"},
				OriginImageURLs: []string{"https://i0.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "https://i0.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "https://i0.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg"},
			},
			Content: "content",
		}
	)
	Convey("creation article operations", t, func() {
		Convey("add article", func() {
			tx, err := d.BeginTran(c)
			So(err, ShouldBeNil)
			var meta = &model.Meta{}
			*meta = *art.Meta
			aid, err = d.TxAddArticleMeta(c, tx, meta, 0)
			So(err, ShouldBeNil)
			err = d.TxAddArticleContent(c, tx, aid, art.Content, []string{})
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)
			Convey("get article", func() {
				res, err1 := d.CreationArticleMeta(c, aid)
				So(err1, ShouldBeNil)
				art.ID = aid
				res.Ctime = 0
				So(res, ShouldResemble, art.Meta)

				content, err2 := d.CreationArticleContent(c, aid)
				So(err2, ShouldBeNil)
				So(content, ShouldEqual, art.Content)
			})
			Convey("list should not be empty", func() {
				res, err1 := d.UpperArticlesMeta(c, art.Author.Mid, 0, 1)
				So(err1, ShouldBeNil)
				So(res, ShouldNotBeEmpty)
			})
			Convey("count should > 0", func() {
				var cnt = &model.CreationArtsType{}
				cnt, err = d.UpperArticlesTypeCount(c, 8167601)
				So(err, ShouldBeNil)
				So(cnt.All, ShouldBeGreaterThan, 0)
			})
			Convey("update state", func() {
				err = d.UpdateArticleState(c, aid, model.StateLock)
				So(err, ShouldBeNil)
				res3, err := d.CreationArticleMeta(c, aid)
				So(err, ShouldBeNil)
				So(res3.State, ShouldEqual, model.StateLock)
			})
			Convey("delete article", func() {
				tx, err := d.BeginTran(c)
				err = d.TxDeleteArticleContent(c, tx, aid)
				So(err, ShouldBeNil)
				err = d.TxDeleteArticleMeta(c, tx, aid)
				So(err, ShouldBeNil)
				err = tx.Commit()
				Convey("article not be present", func() {
					res, err := d.CreationArticleMeta(c, aid)
					So(err, ShouldBeNil)
					So(res, ShouldBeNil)
					content, err := d.CreationArticleContent(c, aid)
					So(err, ShouldBeNil)
					So(content, ShouldBeEmpty)
				})
			})
			Convey("update article", func() {
				art := model.Article{
					Meta: &model.Meta{
						ID:              aid,
						Title:           "new",
						Summary:         "new",
						BannerURL:       "https://i0.hdslb.com/bfs/archive/1.jpg",
						TemplateID:      4,
						State:           2,
						Category:        &model.Category{ID: 2},
						Author:          &model.Author{Mid: 123},
						Reprint:         0,
						ImageURLs:       []string{"https://i0.hdslb.com/bfs/archive/2.jpg"},
						OriginImageURLs: []string{"https://i0.hdslb.com/bfs/archive/3.jpg"},
					},
					Content: "new",
				}
				tx, err := d.BeginTran(c)
				var meta = &model.Meta{}
				*meta = *art.Meta
				err = d.TxUpdateArticleMeta(c, tx, meta)
				So(err, ShouldBeNil)
				err = d.TxUpdateArticleContent(c, tx, aid, art.Content, []string{})
				So(err, ShouldBeNil)
				err = tx.Commit()
				So(err, ShouldBeNil)
				Convey("article should be updated", func() {
					res, err := d.CreationArticleMeta(c, aid)
					So(err, ShouldBeNil)
					art.Ctime = res.Ctime // ignore ctime
					So(res, ShouldResemble, art.Meta)
					content, err := d.CreationArticleContent(c, aid)
					So(err, ShouldBeNil)
					So(content, ShouldEqual, art.Content)
				})
			})
		})
	})
}
