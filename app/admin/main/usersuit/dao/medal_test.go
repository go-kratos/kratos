package dao

import (
	"context"
	"go-common/app/admin/main/usersuit/model"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_hit(t *testing.T) {
	Convey("return someting", t, func() {
		re := d.hit(1)
		So(re, ShouldEqual, "1")
	})
}

func TestDao_Medal(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.Medal(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDao_MedalByID(t *testing.T) {
	Convey("return someting", t, func() {
		re, err := d.MedalByID(context.Background(), 1)
		So(err, ShouldBeNil)
		So(re, ShouldNotBeNil)
	})
}

func TestDao_AddMedal(t *testing.T) {
	Convey("return someting", t, func() {
		pg := &model.Medal{
			Name:        "知名偶像",
			Description: "红白出道，拯救高校",
			Condition: "所有自制视频总播放数>=100万	",
			GID:        4,
			Level:      int8(3),
			Sort:       3,
			LevelRank:  "100万",
			IsOnline:   1,
			Image:      "/bfs/face/27a952195555e64508310e366b3e38bd4cd143fc.png",
			ImageSmall: "/bfs/face/0497be49e08357bf05bca56e33a0637a273a7610.png",
		}
		id, err := d.AddMedal(context.Background(), pg)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeNil)
	})
}

func TestDao_UpMedal(t *testing.T) {
	Convey("return someting", t, func() {
		pg := &model.Medal{
			Name:        "test",
			Description: "Description",
			Condition:   "Condition",
			GID:         1,
			Level:       int8(3),
			Sort:        1,
			LevelRank:   "LevelRank",
			IsOnline:    1,
			Image:       "Image",
			ImageSmall:  "ImageSmall",
		}
		id, err := d.UpMedal(context.Background(), 1, pg)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeNil)
	})
}

func TestDao_MedalGroup(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.MedalGroup(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDao_MedalGroupInfo(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.MedalGroupInfo(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDao_MedalGroupParent(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.MedalGroupParent(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDao_MedalGroupByID(t *testing.T) {
	Convey("return someting", t, func() {
		re, err := d.MedalGroupByID(context.Background(), 2)
		So(err, ShouldBeNil)
		So(re, ShouldNotBeNil)
	})
}

func TestDao_MedalGroupAdd(t *testing.T) {
	Convey("return someting", t, func() {
		pg := &model.MedalGroup{
			Name:     "test",
			PID:      1,
			Rank:     int8(1),
			IsOnline: int8(1),
		}
		id, err := d.MedalGroupAdd(context.Background(), pg)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeNil)
	})
}

func TestDao_MedalGroupUp(t *testing.T) {
	Convey("return someting", t, func() {
		pg := &model.MedalGroup{
			Name:     "test111",
			PID:      2,
			Rank:     2,
			IsOnline: 0,
		}
		id, err := d.MedalGroupUp(context.Background(), 37, pg)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeNil)
	})
}

func Test_MedalOwnerAdd(t *testing.T) {
	mid := int64(rand.Int31())
	Convey("return someting", t, func() {
		id, err := d.MedalOwnerAdd(context.Background(), mid, 1)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeNil)
	})
	Convey("return someting", t, func() {
		_, err := d.MedalOwner(context.Background(), mid)
		So(err, ShouldBeNil)
	})
}

func TestDao_MedalAddList(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.MedalAddList(context.Background(), 111)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDao_MedalOwnerUpActivated(t *testing.T) {
	Convey("return someting", t, func() {
		id, err := d.MedalOwnerUpActivated(context.Background(), 1, 1)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeNil)
	})
}
func TestDao_MedalOwnerUpNotActivated(t *testing.T) {
	Convey("return someting", t, func() {
		id, err := d.MedalOwnerUpNotActivated(context.Background(), 1, 1)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeNil)
	})
}

func TestDao_MedalOwnerDel(t *testing.T) {
	Convey("return someting", t, func() {
		id, err := d.MedalOwnerDel(context.Background(), 1, 1, 1)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeNil)
	})
}
