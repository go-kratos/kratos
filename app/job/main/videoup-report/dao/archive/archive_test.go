package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"time"
)

func Test_Archive(t *testing.T) {
	Convey("ArchiveByAid", t, func() {
		archive, err := d.ArchiveByAid(context.TODO(), 1)
		So(err, ShouldBeNil)
		Println(archive)
	})
}

func Test_ExcitationArchivesByTime(t *testing.T) {
	Convey("ExcitationArchivesByTime", t, func() {
		now := time.Now()
		st := now.Add(-1680000 * time.Hour)
		archives, err := d.ExcitationArchivesByTime(context.TODO(), 27515256, st, now)
		So(err, ShouldBeNil)
		Println(archives)
	})
}
