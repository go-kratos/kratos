package http

import (
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/creative/conf"
	"go-common/app/admin/main/creative/service"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	dir, _ := filepath.Abs("../cmd/creative-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	svc = service.New(conf.Conf)
}

func Test_bindTags(t *testing.T) {
	var (
		oidTIDsMap = map[int64][]int64{
			838: {6, 8, 16, 17, 18, 19},
		}
	)
	Convey("bindTags", t, func() {
		res, err := bindTags(oidTIDsMap)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
