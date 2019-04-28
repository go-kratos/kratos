package abtest

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

func TestExperimentLimit(t *testing.T) {
	Convey("get ExperimentLimit all", t, func() {
		res, err := d.ExperimentLimit(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestExperimentByIDs(t *testing.T) {
	Convey("get ExperimentByIDs all", t, func() {
		res, err := d.ExperimentByIDs(ctx(), []int64{1})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestAbServer(t *testing.T) {
	Convey("TestAbServer", t, func() {
		res, err := d.AbServer(context.TODO(),
			"9902822F-DDD1-47BC-A08F-19F746C9CB8459220infoc",
			"phone",
			"iphone",
			"",
			6720,
			1,
		)
		Println(string(res), err)
	})
}
