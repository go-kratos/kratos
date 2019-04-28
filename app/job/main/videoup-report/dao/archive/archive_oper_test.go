package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/archive"
)

func TestDao_AddArchiveOper(t *testing.T) {
	Convey("AddArchiveOper", t, func() {
		c := context.TODO()
		a, _ := d.ArchiveByAid(c, 1)
		id, err := d.AddArchiveOper(context.TODO(), a.ID, a.Attribute, a.TypeID, a.State, a.Round, 0, "随意一个变更", "测试啦")
		So(err, ShouldBeNil)
		Println(id)
	})
}

func Test_LastVideoOperUID(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub int64
	)
	Convey("LastVideoOperUID", t, func() {
		sub, err = d.LastVideoOperUID(c, 2333)
		So(err, ShouldNotBeNil)
		So(sub, ShouldBeZeroValue)
	})
}

func Test_LastVideoOper(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.VideoOper
	)
	Convey("LastVideoOper", t, func() {
		sub, err = d.LastVideoOper(c, 2333)
		So(err, ShouldNotBeNil)
		So(sub, ShouldNotBeNil)
	})
}
