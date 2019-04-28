package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSearchSubtitle(t *testing.T) {
	var (
		arg = &model.SubtitleSearchArg{
			Oid: int64(10131981),
		}
	)
	Convey("search subtitles", t, func() {
		res, err := testDao.SearchSubtitle(context.Background(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("res:%+v", res)
		for _, rpt := range res.Result {
			t.Logf("======%+v", rpt)
		}
	})
}
