package cms

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsLoadArcsMediaMap(t *testing.T) {
	var ctx = context.Background()
	convey.Convey("LoadArcsMediaMap", t, func(c convey.C) {
		aids, errPick := pickIDs(d.db, _pickAids)
		if errPick != nil || len(aids) == 0 {
			fmt.Println("Empty aids ", errPick)
			return
		}
		resMetas, err := d.LoadArcsMediaMap(ctx, aids)
		c.Convey("Then err should be nil.resMetas should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(resMetas, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadVideosMeta(t *testing.T) {
	var (
		ctx  = context.Background()
		cids = []int64{}
	)
	convey.Convey("LoadVideosMeta", t, func(c convey.C) {
		resMetas, err := d.LoadVideosMeta(ctx, cids)
		c.Convey("Then err should be nil.resMetas should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(resMetas, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadArcsMedia(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("LoadArcsMedia", t, func(c convey.C) {
		aids, errPick := pickIDs(d.db, _pickAids)
		if errPick != nil || len(aids) == 0 {
			fmt.Println("Empty aids ", errPick)
			return
		}
		arcs, err := d.LoadArcsMedia(ctx, aids)
		c.Convey("Then err should be nil.arcs should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(arcs, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadArcMeta(t *testing.T) {
	var (
		ctx = context.Background()
		aid = int64(0)
	)
	convey.Convey("LoadArcMeta", t, func(c convey.C) {
		aids, errPick := pickIDs(d.db, _pickAids)
		if errPick != nil || len(aids) == 0 {
			fmt.Println("Empty aids ", errPick)
			return
		}
		aid = aids[0]
		arcMeta, err := d.LoadArcMeta(ctx, aid)
		c.Convey("Then err should be nil.arcMeta should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(arcMeta, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadVideoMeta(t *testing.T) {
	var (
		ctx = context.Background()
		cid = int64(0)
	)
	convey.Convey("LoadVideoMeta", t, func(c convey.C) {
		aids, errPick := pickIDs(d.db, _pickCids)
		if errPick != nil || len(aids) == 0 {
			fmt.Println("Empty aids ", errPick)
			return
		}
		cid = aids[0]
		videoMeta, err := d.LoadVideoMeta(ctx, cid)
		c.Convey("Then err should be nil.videoMeta should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(videoMeta, convey.ShouldNotBeNil)
		})
	})
}
