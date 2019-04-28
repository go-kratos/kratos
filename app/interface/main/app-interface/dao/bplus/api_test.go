package bplus

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestDynamicCount dao ut.
func TestDynamicCount(t *testing.T) {
	Convey("get DynamicCount", t, func() {
		_, err := dao.DynamicCount(ctx(), 27515258)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestFavClips dao ut.
func TestFavClips(t *testing.T) {
	Convey("get FavClips", t, func() {
		_, err := dao.FavClips(ctx(), 27515258, "da6863c38e83fc7a035c8f7a7d9b1c11", "appkey", "phone", "iphone", "ios", 8230, 1, 20)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestFavAlbums dao ut.
func TestFavAlbums(t *testing.T) {
	Convey("get FavAlbums", t, func() {
		_, err := dao.FavAlbums(ctx(), 27515258, "da6863c38e83fc7a035c8f7a7d9b1c11", "appkey", "phone", "iphone", "ios", 8230, 1, 20)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestClips dao ut.
func TestClips(t *testing.T) {
	Convey("get Clips", t, func() {
		_, _, _, err := dao.Clips(ctx(), 27515258, 1, 20)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestAlbums dao ut.
func TestAlbums(t *testing.T) {
	Convey("get Albums", t, func() {
		_, _, _, err := dao.Albums(ctx(), 27515258, 1, 20)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestAllClip dao ut.
func TestAllClip(t *testing.T) {
	Convey("get AllClip", t, func() {
		_, _, err := dao.AllClip(ctx(), 27515258, 20)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestAllAlbum dao ut.
func TestAllAlbum(t *testing.T) {
	Convey("get AllAlbum", t, func() {
		_, _, err := dao.AllAlbum(ctx(), 27515258, 20)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestClipDetail dao ut.
func TestClipDetail(t *testing.T) {
	Convey("get ClipDetail", t, func() {
		_, err := dao.ClipDetail(ctx(), []int64{27515258})
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestAlbumDetail dao ut.
func TestAlbumDetail(t *testing.T) {
	Convey("get AlbumDetail", t, func() {
		_, err := dao.AlbumDetail(ctx(), 27515258, []int64{27515258})
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestGroupsCount dao ut.
func TestGroupsCount(t *testing.T) {
	Convey("get GroupsCount", t, func() {
		_, err := dao.GroupsCount(ctx(), 27515258, 27515258)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestDynamic dao ut.
func TestDynamic(t *testing.T) {
	Convey("get Dynamic", t, func() {
		_, err := dao.Dynamic(ctx(), 27515258)
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestDynamicDetails dao ut.
func TestDynamicDetails(t *testing.T) {
	Convey("get DynamicDetails", t, func() {
		_, err := dao.DynamicDetails(ctx(), []int64{27515258}, "search")
		err = nil
		So(err, ShouldBeNil)
	})
}
