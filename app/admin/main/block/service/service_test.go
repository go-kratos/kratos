package service

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/block/conf"
	"go-common/app/admin/main/block/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.TODO()
)

func TestMain(m *testing.M) {
	defer os.Exit(0)
	flag.Set("conf", "../cmd/block-admin-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New()
	defer s.Close()

	m.Run()
}

func TestService(t *testing.T) {
	Convey("", t, func() {
		s.Ping(c)
		s.Close()
	})
}

func TestBlock(t *testing.T) {
	Convey("block", t, func() {
		var (
			p = &model.ParamBatchBlock{
				MIDs:      []int64{1, 2, 3, 4},
				AdminID:   233,
				AdminName: "233",
				Source:    1,
				Area:      model.BlockAreaNone,
				Reason:    "test",
				Comment:   "test",
				Action:    model.BlockActionLimit,
				Duration:  1,
				Notify:    false,
			}
			pm = &model.ParamBatchRemove{
				MIDs:      []int64{1, 2, 3, 4},
				AdminID:   233,
				AdminName: "233",
				Comment:   "test",
				Notify:    false,
			}
		)
		err := s.BatchBlock(c, p)
		So(err, ShouldBeNil)
		err = s.BatchRemove(c, pm)
		So(err, ShouldBeNil)
	})
}
