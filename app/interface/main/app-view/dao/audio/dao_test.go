package audio

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-view/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-view-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func ctx() context.Context {
	return context.Background()
}

func TestAudioByCids(t *testing.T) {
	Convey("get AudioByCids all", t, func() {
		res, err := d.AudioByCids(ctx(), []int64{1})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
