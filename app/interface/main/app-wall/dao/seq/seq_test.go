package seq

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-wall/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-wall-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestSeqID(t *testing.T) {
	Convey("SeqID", t, func() {
		res, err := d.SeqID(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
