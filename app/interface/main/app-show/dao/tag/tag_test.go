package tag

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestTagInfo(t *testing.T) {
	Convey("TagInfo", t, func() {
		res, err := d.TagInfo(ctx(), 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestHots(t *testing.T) {
	Convey("Hots", t, func() {
		res, err := d.Hots(ctx(), 1, 1, 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestNewArcs(t *testing.T) {
	Convey("NewArcs", t, func() {
		res, err := d.NewArcs(ctx(), 1, 1, 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestSimilarTag(t *testing.T) {
	Convey("SimilarTag", t, func() {
		res, err := d.SimilarTag(ctx(), 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestTagHotsId(t *testing.T) {
	Convey("TagHotsId", t, func() {
		res, err := d.TagHotsId(ctx(), 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestSimilarTagChange(t *testing.T) {
	Convey("SimilarTagChange", t, func() {
		res, err := d.SimilarTagChange(ctx(), 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestDetail(t *testing.T) {
	Convey("Detail", t, func() {
		res, err := d.Detail(ctx(), 1, 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestDetailRanking(t *testing.T) {
	Convey("DetailRanking", t, func() {
		res, err := d.DetailRanking(ctx(), 1, 1, 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
