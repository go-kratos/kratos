package archive

import (
	"context"
	"testing"

	"go-common/app/admin/main/videoup/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetMaxWeight(t *testing.T) {
	Convey("GetMaxWeight", t, WithDao(func(d *Dao) {
		_, err := d.GetMaxWeight(context.Background())
		So(err, ShouldBeNil)
	}))
}

func Test_DelWeightConf(t *testing.T) {
	Convey("DelWeightConf", t, WithDao(func(d *Dao) {
		_, err := d.DelWeightConf(context.Background(), 0)
		So(err, ShouldBeNil)
	}))
}

func Test_ListWeightConf(t *testing.T) {
	Convey("ListWeightConf", t, WithDao(func(d *Dao) {
		_, err := d.ListWeightConf(context.Background(), &archive.Confs{})
		So(err, ShouldBeNil)
	}))
}

func Test_LWConfigHelp(t *testing.T) {
	Convey("LWConfigHelp", t, WithDao(func(d *Dao) {
		_, err := d.LWConfigHelp(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}
