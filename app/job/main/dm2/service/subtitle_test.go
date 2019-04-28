package service

import (
	"context"
	"testing"

	"go-common/app/job/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubtitle(t *testing.T) {
	var (
		oid        int64 = 10109227
		subtitleID int64 = 1
	)
	Convey("", t, func() {
		err := svr.SubtitleFilter(context.Background(), oid, subtitleID)
		So(err, ShouldBeNil)
	})
}

func TestSubtitleFilter(t *testing.T) {
	body := &model.SubtitleBody{
		Bodys: []*model.SubtitleItem{
			{
				From:    0,
				To:      10,
				Content: "习近平",
			},
			{
				From:    0,
				To:      10,
				Content: "习大大",
			},
			{
				From:    0,
				To:      10,
				Content: "不要哇",
			},
			{
				From:    0,
				To:      10,
				Content: "呀咩爹",
			},
		},
	}
	Convey("subtitle filter", t, func() {
		hits, err := svr.checkFilter(context.Background(), body)
		So(err, ShouldBeNil)
		t.Logf("hits:%v", hits)
	})
}
