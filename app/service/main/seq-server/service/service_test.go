package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/service/main/seq-server/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/seq-server-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_Close(t *testing.T) {
	Convey("Close", t, func() {
		s.Close()
	})
}

func Test_ID(t *testing.T) {
	Convey("ID", t, func() {
		fmt.Println(s.ID(context.TODO(), 10, "Nf9phmDdzjTMW9M5V8YQuLpVTwhvn5IO"))
	})
}
