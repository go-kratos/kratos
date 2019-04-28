package relation

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
	time.Sleep(5 * time.Second)
}

func TestPrompt(t *testing.T) {
	Convey("Prompt", t, func() {
		acc, err := d.Prompt(context.TODO(), 1, 1, 1)
		So(err, ShouldBeNil)
		Println(acc)
	})
}

func TestStat(t *testing.T) {
	Convey("Stat", t, func() {
		acc, err := d.Stat(context.TODO(), 1)
		So(err, ShouldBeNil)
		Println(acc)
	})
}
