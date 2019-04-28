package archive

import (
	"context"
	"testing"

	"go-common/app/admin/main/videoup/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFansConf(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		fans, err := d.FansConf(context.Background())
		So(err, ShouldBeNil)
		So(fans, ShouldNotBeNil)
	}))
}

func TestDao_RoundTypeConf(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		r, err := d.RoundTypeConf(context.Background())
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
	}))
}

func TestDao_ThresholdConf(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		r, err := d.ThresholdConf(context.Background())
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
	}))
}

func TestDao_AuditTypesConf(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		r, err := d.AuditTypesConf(context.Background())
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
	}))
}

func TestDao_WeightVC(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		r, err := d.WeightVC(context.Background())
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
	}))
}

func TestDao_SetWeightVC(t *testing.T) {
	Convey("SetWeightVC", t, WithDao(func(d *Dao) {
		_, err := d.SetWeightVC(context.TODO(), &archive.WeightVC{}, "desc")
		So(err, ShouldBeNil)
	}))
}

func TestDao_InWeightVC(t *testing.T) {
	Convey("InWeightVC", t, WithDao(func(d *Dao) {
		_, err := d.InWeightVC(context.TODO(), &archive.WeightVC{}, "desc")
		So(err, ShouldBeNil)
	}))
}
