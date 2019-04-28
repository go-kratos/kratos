package block

import (
	"flag"
	"testing"

	"go-common/app/admin/main/block/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/block-admin-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}

	m.Run()
}

func TestTools(t *testing.T) {
	Convey("tools", t, func() {
	})
}
