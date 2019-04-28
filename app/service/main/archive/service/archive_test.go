package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/service/main/archive/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_MaxAID(t *testing.T) {
	Convey("MaxAID", t, func() {
		var (
			c   = context.TODO()
			err error
			id  int64
		)
		id, err = s.MaxAID(c)
		So(err, ShouldBeNil)
		So(id, ShouldNotEqual, 0)
		Printf("MaxAID = %d\n", id)
	})
}

func Test_Archive(t *testing.T) {
	Convey("Archive", t, func() {
		arc, err := s.Archive3(context.TODO(), 10098500)
		So(err, ShouldBeNil)
		Printf("%+v\n", arc)
		bs, _ := json.Marshal(arc)
		Printf("%s", bs)
	})
}

func Test_ArchiveWithPlayer(t *testing.T) {
	Convey("ArchiveWithPlayer", t, func() {
		arcs, err := s.ArchivesWithPlayer(context.TODO(), &archive.ArgPlayer{
			Aids:     []int64{10111001},
			Qn:       32,
			Platform: "android",
			RealIP:   "121.31.246.238",
			Fnver:    0,
			Fnval:    16,
			Session:  "",
		}, true)
		So(err, ShouldBeNil)
		Printf("%+v\n", arcs[10111001])
	})
}

func Test_Archives3(t *testing.T) {
	Convey("Archives3", t, func() {
		as, err := s.Archives3(context.TODO(), []int64{10098500, 10098501})
		So(err, ShouldBeNil)
		for _, a := range as {
			Printf("%+v\n\n", a)
			bs, _ := json.Marshal(a)
			Printf("%s\n\n", bs)
		}
	})
}
