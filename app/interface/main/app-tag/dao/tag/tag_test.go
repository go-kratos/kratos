package tag

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-tag/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-tag-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestInfoByID(t *testing.T) {
	Convey("get InfoByID all", t, func() {
		res, err := d.InfoByID(ctx(), 0, 1217733)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestTagInfo(t *testing.T) {
	Convey("get TagInfo all", t, func() {
		res, err := d.TagInfo(ctx(), 0, 1217733, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestDetail(t *testing.T) {
	Convey("get Detail all", t, func() {
		res, err := d.Detail(ctx(), 1217733, 1, 20, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestTagHotsId(t *testing.T) {
	Convey("TagHotsId", t, func() {
		res, err := d.TagHotsId(ctx(), 33, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestSimilarTagChange(t *testing.T) {
	Convey("SimilarTagChange", t, func() {
		res, err := d.SimilarTagChange(ctx(), 1217733, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestDetailRanking(t *testing.T) {
	Convey("DetailRanking", t, func() {
		res, err := d.DetailRanking(ctx(), 1, 1217733, 1, 20, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
