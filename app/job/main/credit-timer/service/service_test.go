package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/credit-timer/conf"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}
func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

// TestService_loadConf
func Test_loadConf(t *testing.T) {
	Convey("should return err be nil", t, func() {
		s.loadConf(context.TODO())
		So(s.c.Judge.CaseEndVoteTotal, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

// TestService_ComputePoint
func Test_ComputePoint(t *testing.T) {
	Convey("should return err be nil", t, func() {
		r, err := s.ComputePoint(context.TODO(), 88889017)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
		// So(r, ShouldResemble,model.KpiPoint{})
	})
}

func TestService_sort(t *testing.T) {
	var (
		res = []int64{1, 3, 4, 5, 5, 6, 7, 9, 9, 10, 11, 12}
		ps  []int64
	)
	for _, k := range res {
		if len(ps) == 0 {
			ps = append(ps, k)
			continue
		}
		if ps[len(ps)-1] == k {
			continue
		}
		ps = append(ps, k)
	}
	t.Logf("%v", ps)
	for _, k := range res {
		for i, r := range ps {
			if r == k {
				t.Logf("%d,%d,%d", k, i+1, (i+1)*100/len(ps))
				break
			}
		}
	}
}

func TestService_point(t *testing.T) {
	var (
		point     float64
		voteTotal int64
		voteRight int64
		//投准率
		vr float64
		//投准率系数
		vf float64
	)
	voteTotal = 12
	voteRight = 5
	vr = float64(voteRight) / float64(voteTotal)
	vf = float64(1.0)
	point = float64(voteTotal) * vr * vf
	t.Logf("%f", point)
}

func initConf(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
}

func Test_Time(t *testing.T) {
	Convey("should return err be nil", t, func() {
		d := time.Now().AddDate(0, 0, 1)
		ts1 := time.Until(time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 1, 0, time.Local))
		ts2 := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 1, 0, time.Local).Sub(time.Now())
		t.Errorf("%#v %#v", ts1, ts2)
		So(ts1, ShouldEqual, ts2)
	})
}

func Test_FixKPI(t *testing.T) {
	Convey("should return err be nil", t, func() {
		res, err := s.FixKPI(context.TODO(), 2018, 2, 10, 7584862)
		fmt.Printf("res:%+v \n", res)
		So(err, ShouldBeNil)
	})
}
