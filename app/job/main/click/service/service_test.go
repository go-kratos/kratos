package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/click/conf"
	"go-common/app/job/main/click/model"
	"go-common/library/sync/errgroup"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/click-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_Archive(t *testing.T) {
	Convey("isReplay test", t, func() {
		b := s.isReplay(context.TODO(), 1, 1, "", 1)
		So(b, ShouldBeFalse)
	})
	return
}

func Test_SetSpecial(t *testing.T) {
	Convey("SetSpecial", t, func() {
		s.SetSpecial(context.TODO(), 1, 10, "")
	})
}

func Test_ArcDuration(t *testing.T) {
	Convey("ArcDuration", t, func() {
		eg, _ := errgroup.WithContext(context.TODO())
		var aids = []int64{10099744, 10099730, 10099729, 10099728, 10099726, 10099670, 10099669, 10099668, 10099667, 10099660, 10099653, 10099652, 10099651}
		for _, aid := range aids {
			id := aid
			eg.Go(func() (err error) {
				Println(s.ArcDuration(context.TODO(), id))
				return
			})
		}
		eg.Wait()
		Println(s.arcDurWithMutex.Durations)
		for aid, val := range s.arcDurWithMutex.Durations {
			Printf("aid(%d) dur(%d) gotTime(%d)\n", aid, val.Duration, val.GotTime)
		}
	})
}

func Test_SendStat(t *testing.T) {
	Convey("SendStat", t, func() {
		vmsg := &model.StatViewMsg{Type: "archive", ID: 123456, Count: 123456, Ts: time.Now().Unix()}
		err := s.statViewPub.Send(context.TODO(), "test", vmsg)
		So(err, ShouldBeNil)
	})
}
