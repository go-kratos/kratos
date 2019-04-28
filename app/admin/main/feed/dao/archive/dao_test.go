package archive

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/feed/conf"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/feed-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(5 * time.Second)
}

func Test_Archive3(t *testing.T) {
	Convey("Card", t, func() {
		arc, err := d.Archive3(context.TODO(), 10099667)
		if err != nil {
			if err == ecode.NothingFound {
				fmt.Println("NothingFound")
				return
			}
		}
		So(err, ShouldBeNil)
		v, err := json.Marshal(arc)
		if err != nil {
			fmt.Println(err)
			return
		}
		if v == nil {
			fmt.Println("empty value")
		} else {
			fmt.Println(string(v))
		}
	})
}
