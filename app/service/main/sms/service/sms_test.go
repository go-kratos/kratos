package service

import (
	"context"
	"flag"
	"path/filepath"
	"strings"
	"testing"

	pb "go-common/app/service/main/sms/api"
	"go-common/app/service/main/sms/conf"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/sms-service-test.toml")
	flag.Set("conf", dir)
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	s = New(conf.Conf)
}

func Test_Send(t *testing.T) {
	Convey("send", t, func() {
		req := &pb.SendReq{
			Mobile: "17621660828",
			Tcode:  "acc_01111",
		}
		_, err := s.Send(context.Background(), req)
		So(err, ShouldBeNil)
	})
}

func Test_Param(t *testing.T) {
	Convey("solve param", t, func() {
		var (
			template = "test#[code].very#[code] code"
			strs     []string
			param    = make(map[string]string)
			ss       []string
		)
		strs = strings.SplitAfter(template, "#[")
		for i, v := range strs {
			if i == 0 {
				continue
			}
			t.Logf("%s", v)
			t.Logf("%d", strings.Index(v, "]"))
			k := v[0:strings.Index(v, "]")]
			param[k] = ""
		}
		for k := range param {
			ss = append(ss, k)
		}
		t.Logf("%v", ss)
		t.Logf(strings.Replace(template, "#[codes]", "888", -1))
	})
}
