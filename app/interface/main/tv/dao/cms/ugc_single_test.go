package cms

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/tv/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsArcMetaCache(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(0)
	)
	convey.Convey("ArcMetaCache", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.s should not be nil.", func(ctx convey.C) {
			sids, err := pickIDs(d.db, _pickAids)
			if err != nil || len(sids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			aid = sids[0]
			d.LoadArcMeta(c, aid)
			s, err := d.ArcMetaCache(c, aid)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(s, convey.ShouldNotBeNil)
		})
		ctx.Convey("mc not found Error", func(ctx convey.C) {
			_, err := d.ArcMetaCache(c, 0)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCmsSetArcMetaCache(t *testing.T) {
	var (
		c = context.Background()
		s = &model.ArcCMS{}
	)
	convey.Convey("SetArcMetaCache", t, func(ctx convey.C) {
		err := d.SetArcMetaCache(c, s)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCmsAddArcMetaCache(t *testing.T) {
	var (
		arc = &model.ArcCMS{}
	)
	convey.Convey("AddArcMetaCache", t, func(ctx convey.C) {
		d.AddArcMetaCache(arc)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestCmsVideoMetaCache(t *testing.T) {
	var (
		c   = context.Background()
		cid = int64(0)
	)
	convey.Convey("VideoMetaCache", t, func(ctx convey.C) {
		sids, err := pickIDs(d.db, _pickCids)
		if err != nil || len(sids) == 0 {
			fmt.Println("Empty Sids ", err)
			return
		}
		cid = sids[0]
		d.LoadVideoMeta(c, cid)
		s, err := d.VideoMetaCache(c, cid)
		ctx.Convey("Then err should be nil.s should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(s, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsSetVideoMetaCache(t *testing.T) {
	var (
		c = context.Background()
		s = &model.VideoCMS{}
	)
	convey.Convey("SetVideoMetaCache", t, func(ctx convey.C) {
		err := d.SetVideoMetaCache(c, s)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCmsAddVideoMetaCache(t *testing.T) {
	var (
		video = &model.VideoCMS{}
	)
	convey.Convey("AddVideoMetaCache", t, func(ctx convey.C) {
		d.AddVideoMetaCache(video)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}
