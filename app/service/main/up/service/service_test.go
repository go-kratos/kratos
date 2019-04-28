package service

import (
	"context"
	"flag"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/model"
	"go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	s    *Service
	pCfg = &databus.Config{
		Key:          "4c76cbb7a985ac90",
		Secret:       "f36fbb15a85c6e21b0ee22a560ef3a67",
		Group:        "Creative-MainArchive-P",
		Topic:        "Creative-T",
		Action:       "pub",
		Name:         "up-service/pub",
		Proto:        "tcp",
		Addr:         "172.16.33.158:6205",
		Active:       10,
		Idle:         5,
		DialTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		IdleTimeout:  xtime.Duration(time.Minute),
	}
)

func init() {
	dir, _ := filepath.Abs("../cmd/up-service.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_Attr(t *testing.T) {
	mdlUp := &model.Up{
		Attribute: 0,
	}
	Convey("testAttr", t, WithService(func(s *Service) {
		Convey("AttrSet", func() {
			mdlUp.AttrSet(1, model.AttrArchiveNewUp)
		})
		Convey("AttrVal", func() {
			attr := mdlUp.AttrVal(model.AttrArchiveNewUp)
			intAttr := model.InAttr(attr)
			spew.Dump(attr, intAttr)
		})
	}))
}

func Test_service(t *testing.T) {
	var (
		c    = context.Background()
		from = uint8(1)
	)
	Convey("db", t, WithService(func(s *Service) {
		Convey("Ping", func() {
			s.Ping(context.Background())
		})
	}))
	Convey("up", t, WithService(func(s *Service) {
		Convey("AddOrUpdateDB", func() {
			_, err := s.Edit(c, int64(2089809), 1, from)
			So(err, ShouldBeNil)
		})
		Convey("Info", func() {
			isAuthor, err := s.Info(c, int64(2089809), from)
			So(err, ShouldBeNil)
			println(isAuthor)
		})
	}))
	Convey("db", t, WithService(func(s *Service) {
		Convey("Close", func() {
			s.Close()
		})
	}))
}

//func testSub(t *testing.T, d *databus.Databus) {
//	for {
//		m, ok := <-d.Messages()
//		if !ok {
//			return
//		}
//		t.Logf("sub message: %s", string(m.Value))
//		if err := m.Commit(); err != nil {
//			t.Errorf("sub commit error(%v)\n", err)
//		}
//	}
//}

func testPub(t *testing.T, mid int64, from int, isauthor int) {
	pub := databus.New(pCfg)
	m := &model.Msg{
		MID:       mid,
		From:      from,
		IsAuthor:  isauthor, //偶数=0,基数=1
		TimeStamp: time.Now().Unix(),
	}
	if err := pub.Send(context.Background(), strconv.FormatInt(m.MID, 10), m); err != nil {
		t.Errorf("d.Send(test) error(%v)", err)
	} else {
		t.Logf("pub message %v", m)
	}
}

func Test_databus(t *testing.T) {
	mid1 := int64(1)
	testPub(t, mid1, 0, 1)
	testPub(t, mid1, 1, 0)
	testPub(t, mid1, 2, 0)
	//go testSub(t, s.upSub)
	//time.Sleep(time.Second * 3)
	//s.upSub.Close()
}

func TestService_SpecialAdd(t *testing.T) {
	var (
		c = &blademaster.Context{}
	)
	Convey("service", t, WithService(func(s *Service) {
		Convey("SpecialAdd", func() {
			var _, e = s.SpecialAdd(c, "", &model.UpSpecial{}, 12345)
			So(e, ShouldNotBeNil)
		})
	}))
}

func TestService_BmGetStringOrDefault(t *testing.T) {
	var (
		c = &blademaster.Context{}
	)
	Convey("service", t, WithService(func(s *Service) {
		Convey("BmGetStringOrDefault", func() {
			var _, ok = BmGetStringOrDefault(c, "t", "default")
			So(ok, ShouldBeFalse)
		})
	}))
}
