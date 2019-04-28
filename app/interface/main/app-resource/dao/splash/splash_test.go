package splash

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-resource-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestActiveAll(t *testing.T) {
	Convey("get ActiveAll all", t, func() {
		res, err := d.ActiveAll(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestActiveBirth(t *testing.T) {
	Convey("get ActiveBirth all", t, func() {
		res, err := d.ActiveBirth(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestActiveVip(t *testing.T) {
	Convey("get ActiveVip all", t, func() {
		res, err := d.ActiveVip(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
