package service

import (
	"testing"
	"time"

	"go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Hotsports(t *testing.T) {
	_hotspotArtTime = time.Hour * 24 * 365 * 10
	Convey("gen data", t, func() {
		err := s.UpdateHotspots(true)
		So(err, ShouldBeNil)
		res, err := s.dao.CacheHotspots(c)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		hot := res[0]
		Convey("add art", func() {
			err = s.AddCacheHotspotArt(c, &model.SearchArt{StatsView: 100, ID: 100, Tags: []string{hot.Tag}, PublishTime: time.Now().Unix()})
			So(err, ShouldBeNil)
			var hot2 *model.Hotspot
			hot2, err = s.dao.CacheHotspot(c, hot.ID)
			So(hot2.Stats.Read-hot.Stats.Read, ShouldEqual, 100)
		})
		Convey("get art", func() {
			var arts []*model.MetaWithLike
			var hotspot *model.Hotspot
			hotspot, arts, err = s.HotspotArts(c, hot.ID, 1, 100, nil, 0, 0)
			So(err, ShouldBeNil)
			So(arts, ShouldNotBeEmpty)
			So(hotspot, ShouldNotBeEmpty)
		})
	})
}
