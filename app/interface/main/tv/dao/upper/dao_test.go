package upper

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/tv/conf"

	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx = context.TODO()
	d   *Dao
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}

func TestDao_LoadUpMeta(t *testing.T) {
	Convey("TestDao_LoadUpMeta", t, WithDao(func(d *Dao) {
		res, err := d.LoadUpMeta(ctx, 27515256)
		So(err, ShouldBeNil)
		data, _ := json.Marshal(res)
		fmt.Println(string(data))
	}))
}
