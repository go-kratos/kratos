package aids

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestSave(t *testing.T) {
	Convey("Save aids", t, WithService(func(s *Service) {
		err := s.Save(context.TODO(), "1,2,3,4,5")
		So(err, ShouldBeNil)
	}))
}
