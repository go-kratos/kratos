package medal

import (
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/usersuit/conf"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_Close(t *testing.T) {
	Convey("should return err be nil", t, func() {
		d.Close()
	})
}
