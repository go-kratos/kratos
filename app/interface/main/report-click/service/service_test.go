package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/report-click/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var svr *Service

func init() {
	flag.Parse()
	dir, _ := filepath.Abs("../cmd/report-click.toml")
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
}

func TestReport(t *testing.T) {
	var (
		c          = context.Background()
		err        error
		aid        = "11159485"
		cid        = "18464413"
		mid        = "35152246"
		playedTime = "0"
		realtime   = "0"
		tp         = "3"
		dt         = "2"
		bs         = []byte("this is  test massage ")
	)
	Convey("Decrypt Verify err should return nil", t, func() {
		bs, _ = svr.Decrypt(bs, conf.Conf.Click.AesKey, conf.Conf.Click.AesIv)
		svr.Verify(bs, conf.Conf.Click.AesSalt, time.Now())
	})
	Convey("Report err should return nil", t, func() {
		err = svr.Report(c, playedTime, cid, tp, "", realtime, aid, mid, "", "", dt, "1516695880")
		So(err, ShouldBeNil)
	})
	Convey("CheckDid err should return nil", t, func() {
		svr.CheckDid("1516695880")
	})
	Convey("GenDid err should return nil", t, func() {
		svr.GenDid("127.0.0.1", time.Now())
	})
	Convey("Play err should return nil", t, func() {
		svr.Play(c, "web", "128546345", "12345", "", "14771787", "1", "", "", "", "127.0.0.1", "2", "1", "1212", "", "1", "3", "2", "4", "", "", "", "", "", "")
	})
}
