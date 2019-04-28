package account

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

func Test_Card(t *testing.T) {
	Convey("Card", t, func() {
		acc, err := d.Card3(context.TODO(), 400062)
		if err != nil {
			if err.Error() == ecode.MemberNotExist.Error() {
				fmt.Println("MemberNotExist")
				return
			}
		}
		So(err, ShouldBeNil)
		v, err := json.Marshal(acc)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		if v == nil {
			fmt.Println("empty value")
		} else {
			fmt.Println(string(v))
		}
	})
}
