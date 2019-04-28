package service

import (
	"context"
	"testing"

	artmdl "go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	cats = []*artmdl.Category{
		&artmdl.Category{Name: "游戏", ID: 1},
		&artmdl.Category{Name: "动漫", ID: 2},
	}
	draft = artmdl.Draft{
		Article: &artmdl.Article{
			Meta: &artmdl.Meta{
				Category:    cats[0],
				Title:       "隐藏于时区记忆中的,是希望还是绝望!",
				Summary:     "说起日本校服,第一个浮现在我们脑海中的必然是那象征着青春阳光 蓝白色相称的水手服啦. 拉色短裙配上洁白的直袜",
				BannerURL:   "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg",
				TemplateID:  1,
				State:       artmdl.StatePending,
				Author:      &artmdl.Author{Mid: 8167601, Name: "爱蜜莉雅", Face: "http://i1.hdslb.com/bfs/face/5c6109964e78a84021299cdf71739e21cd7bc208.jpg"},
				Reprint:     0,
				ImageURLs:   []string{"http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg"},
				PublishTime: 1495784507,
				Stats:       &artmdl.Stats{Favorite: 100, Like: 10, View: 500, Dislike: 1, Share: 99},
			},
			Content: "test content",
		},
		Tags: []string{"tag1", "tag2"},
	}
)

func Test_Draft(t *testing.T) {
	var (
		err error
		aid int64
		c   = context.TODO()
	)
	Convey("creation draft", t, WithService(func(s *Service) {
		Convey("AddArtDraft", func() {
			aid, err = s.AddArtDraft(c, &draft)
			t.Logf("aid: %d", aid)
			So(err, ShouldBeNil)
			So(aid, ShouldBeGreaterThan, 0)
			// t.Logf("result: %+v", aid)

			Convey("ArtDraft", func() {
				res, err := s.ArtDraft(c, aid, draft.Author.Mid)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeEmpty)
				t.Logf("result: %+v", res.Title)
				// t.Logf("result: %+v", res.Content)

				Convey("DelArtDraft", func() {
					err = s.DelArtDraft(c, aid, draft.Author.Mid)
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("UpperDrafts", func() {
			mid := art.Author.Mid
			pn := 1
			ps := 10
			res2, err := s.UpperDrafts(c, mid, pn, ps)
			So(err, ShouldBeNil)
			So(res2, ShouldNotBeNil)
			// fmt.Println("res2", res2.Page)
			// fmt.Println("res2", len(res2.Drafts))
			// fmt.Printf("meta %+v:", res2.Drafts[0].Meta)
			// fmt.Printf("category %+v", res2.Drafts[0].Category)
		})
	}))
}
