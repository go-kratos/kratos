package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	categories = []*artmdl.Category{
		&artmdl.Category{Name: "游戏", ID: 1},
		&artmdl.Category{Name: "动漫", ID: 2},
	}

	art = artmdl.Article{
		Meta: &artmdl.Meta{
			Category:        categories[0],
			Title:           "隐藏于时区记忆中的,是希望还是绝望!",
			Summary:         "说起日本校服,第一个浮现在我们脑海中的必然是那象征着青春阳光 蓝白色相称的水手服啦. 拉色短裙配上洁白的直袜",
			BannerURL:       "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg",
			TemplateID:      4,
			State:           artmdl.StatePending,
			Author:          &artmdl.Author{Mid: 8167601, Name: "爱蜜莉雅", Face: "http://i1.hdslb.com/bfs/face/5c6109964e78a84021299cdf71739e21cd7bc208.jpg"},
			Reprint:         0,
			ImageURLs:       []string{"http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg"},
			OriginImageURLs: []string{"http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg"},
			PublishTime:     1495784507,
			Stats:           &artmdl.Stats{Favorite: 100, Like: 10, View: 500, Dislike: 1, Share: 99},
			Attributes:      2,
			Words:           555,
			Dynamic:         "dynamic",
			Tags:            []*artmdl.Tag{&artmdl.Tag{Name: "tag"}},
		},
		Content: "海奉是一个风景优美的地方，但并不在沿海。数量众多的旅行家笔记显示，海奉是一片死火山群。那里坐落着世界上最高的山峰——奈文摩尔峰，峰顶终年积雪。其它沉睡的火山围坐在他的周围，高低不同，错落有致。火山口往往积蓄湖水，形成湖泊，当地人称之为“镜湖”。每到雨季，经过连续的降雨，湖中的水便会溢出，从山顶冲下，形成“水山爆发”的情景。山脚下是海奉人的村落，那里的房子全部以木头搭建，巧妙的避开河水的必经之路。海奉人以木工闻名，无论是精巧的木头机械还是美丽的木雕都不在话下。此外，每一个海奉人都戴着一枚木制的十字架，那是由海奉独有的铁木制成，绝不出售给外人，因而成为海奉人的标志。但是故事并不发生在海奉，这些描写仅是因为主角是海奉人。船还在航行。天色昏暗，雨从来没有停过。船舱紧闭，窗口透出一丝微弱的光。“您是海奉人吗？”山本真奈美借着微弱的灯光盯着他的十字架。",
	}
)

func Test_Creation(t *testing.T) {
	var (
		err error
		aid int64
		c   = context.TODO()
	)
	Convey("creation article", t, WithService(func(s *Service) {
		res := &artmdl.Article{
			Meta: &artmdl.Meta{
				Title:      "title",
				Category:   &artmdl.Category{ID: 39},
				Author:     &artmdl.Author{Mid: art.Author.Mid},
				Tags:       []*artmdl.Tag{&artmdl.Tag{Name: "tag"}},
				TemplateID: 1,
			},
		}
		Convey("AddArticle", func() {
			art.Title = fmt.Sprintf("隐藏于时区记忆中的,是希望还是绝望!_%v", time.Now().UnixNano())
			aid, err = s.AddArticle(c, &art, 0, 0, "")
			So(err, ShouldBeNil)
			So(aid, ShouldBeGreaterThan, 0)

			err = s.SetStat(c, aid, &artmdl.Stats{
				View:  10,
				Reply: 9,
				Like:  5,
			})
			So(err, ShouldBeNil)

			Convey("CreationArticle", func() {
				res, err = s.CreationArticle(c, aid, art.Author.Mid)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeEmpty)
			})

			Convey("UpdateArticleDB", func() {
				art2 := &artmdl.Article{Meta: &artmdl.Meta{}}
				*art2 = art
				art2.ID = aid
				art2.ImageURLs = []string{"http://i2.hdslb.com/bfs/archive/00.jpg"}
				art2.OriginImageURLs = []string{"http://i2.hdslb.com/bfs/archive/01.jpg"}
				art2.Dynamic = "update 2"
				err = s.updateArticleDB(c, art2)
				So(err, ShouldBeNil)
				res, err = s.CreationArticle(c, art2.ID, art2.Author.Mid)
				So(err, ShouldBeNil)
				So(res.Dynamic, ShouldEqual, art2.Dynamic)
				So(res.ImageURLs, ShouldResemble, []string{"https://i0.hdslb.com/bfs/archive/00.jpg"})
				So(res.OriginImageURLs, ShouldResemble, []string{"https://i0.hdslb.com/bfs/archive/01.jpg"})
			})

			Convey("UpdateArticle", func() {
				// res.Content = art.Content
				// res.ID = aid
				// err = s.UpdateArticle(c, res, 0, 0, "")
				// So(err, ShouldBeNil)
				err = s.dao.UpdateArticleState(c, aid, artmdl.StateOpen)
				So(err, ShouldBeNil)

				Convey("DelArticle", func() {
					err = s.DelArticle(c, aid, art.Author.Mid)
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("CreationUpperArticlesMeta", func() {
			mid := art.Author.Mid
			group := 0
			category := 1
			sortType := 1
			pn := 1
			ps := 10
			res2, err := s.CreationUpperArticlesMeta(c, mid, group, category, sortType, pn, ps, "")
			So(err, ShouldBeNil)
			So(res2, ShouldNotBeNil)
		})
	}))
}
