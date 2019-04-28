package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	flag.Parse()
	dir, _ := filepath.Abs("../cmd/workflow-admin-develop.toml")
	if err := flag.Set("conf", dir); err != nil {
		panic(err)
	}

	s = New()
}

func TestPing(t *testing.T) {
	convey.Convey("Ping", t, func() {
		err := s.Ping(context.Background())
		convey.So(err, convey.ShouldBeNil)
	})
}
