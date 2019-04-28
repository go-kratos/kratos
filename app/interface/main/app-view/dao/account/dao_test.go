package account

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

func Test_Card(t *testing.T) {
	Convey("Card", t, func() {
		acc, err := d.Card3(context.TODO(), 2)
		So(err, ShouldBeNil)
		Println(acc)
	})
}

func Test_Following(t *testing.T) {
	Convey("Following", t, func() {
		_, err := d.Following3(context.TODO(), 2, 1684013)
		So(err, ShouldBeNil)
	})
}
