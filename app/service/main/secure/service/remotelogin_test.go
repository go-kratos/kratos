package service

import (
	"context"
	"flag"
	"sync"
	"testing"
	"time"

	"go-common/app/service/main/secure/conf"
	model "go-common/app/service/main/secure/model"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	s    *Service
	ctx  = context.Background()
)

func startService() {
	initConf()
	s = New(conf.Conf)
	time.Sleep(time.Second * 2)
}

func initConf() {
	flag.Set("conf", "../cmd/secure-service-test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
}

func (s *Service) initCache(mid int64, uuid string) {
	s.delUnNotify(ctx, mid)
	s.dao.DelCount(ctx, mid)
}
func TestLoginLog(t *testing.T) {
	once.Do(startService)
	commonlog1 := []byte(`{"id":79839245,"loginip":288673833,"mid":3,"server":"","timestamp":1498016443,"type":4}`)
	commonlog2 := []byte(`{"id":79839245,"loginip":288673823,"mid":3,"server":"","timestamp":1498016444,"type":4}`)
	commonlog3 := []byte(`{"id":79839245,"loginip":288673823,"mid":3,"server":"","timestamp":1498016445,"type":4}`)
	commonlog4 := []byte(`{"id":79839245,"loginip":288673823,"mid":3,"server":"","timestamp":1498016446,"type":4}`)
	commonlog5 := []byte(`{"id":79839245,"loginip":288673823,"mid":3,"server":"","timestamp":1498016447,"type":4}`)
	err := s.loginLog(ctx, "insert", commonlog1)
	if err != nil {
		t.Errorf("err %v", err)
	}
	s.commonLoc(ctx, 3)
	s.loginLog(ctx, "insert", commonlog2)
	s.loginLog(ctx, "insert", commonlog3)
	s.loginLog(ctx, "insert", commonlog4)
	s.loginLog(ctx, "insert", commonlog5)
}
func TestCloseNotify(t *testing.T) {
	once.Do(startService)
	s.initCache(1, "1234")
	s.CloseNotify(ctx, 1, "1234")
	b, err := s.dao.UnNotify(ctx, 1, "1234")
	if err != nil {
		t.Error("s.dao.UnNotify err", err)
	}
	if !b {
		t.Errorf("unnotify want true  but get %v", b)
	}
	count, _ := s.dao.Count(ctx, 1, "1234")
	if count != 1 {
		t.Errorf("user close count want 1,but get %d", count)
	}
}
func TestStatus(t *testing.T) {
	once.Do(startService)
	s.initCache(2, "abc")
	commonlog := []byte(`{"id":79839246,"loginip":2886738232,"mid":2,"server":"","timestamp":1498016443,"type":4}`)
	expectlog := []byte(`{"id":79839246,"loginip":3078818617,"mid":2,"server":"","timestamp":1498016443,"type":4}`)
	for i := 0; i < 4; i++ {
		s.loginLog(ctx, "insert", commonlog)
	}
	s.loginLog(ctx, "insert", expectlog)
	res, err := s.Status(ctx, 2, "abc")
	if err != nil || res == nil {
		t.Fatalf("s.Status err %v ", err)
	}
	if !res.Notify {
		t.Errorf("s.Status want notify true  but get %v", res.Notify)
	}
	t.Logf("status notify %v ", res)
	s.CloseNotify(ctx, 2, "abc")
	res, err = s.Status(ctx, 2, "abc")
	if err != nil || res == nil {
		t.Fatalf("s.Status err %v ", err)
	}
	if res.Notify {
		t.Errorf("s.Status want notify true  but get %v", res.Notify)
	}
	if err != nil || res == nil {
		t.Fatalf("s.Status err %v ", err)
	}
	if res.Notify {
		t.Errorf("s.Status want notify false  but get %v", res.Notify)
	}
	for i := 0; i < int(s.c.Expect.CloseCount); i++ {
		s.dao.AddCount(ctx, 2, "1234")
	}
	res, err = s.Status(ctx, 2, "abcd")
	if err != nil || res == nil {
		t.Fatalf("s.Status err %v ", err)
	}
	if !res.Notify {
		t.Errorf("s.Status want notify true  but get %v", res.Notify)
	}
}
func TestGetIP(t *testing.T) {
	once.Do(startService)
	res, err := s.getIPZone(ctx, "139.214.144.59")
	if err != nil {
		t.Errorf("getIPZone err(%v)", err)
	}
	t.Logf("res %v", res)
	res, err = s.getIPZone(ctx, "218.27.198.190")
	if err != nil {
		t.Errorf("getIPZone err(%v)", err)
	}
	t.Logf("res %v", res)
}
func TestAddFeedBack(t *testing.T) {
	once.Do(startService)
	s.AddFeedBack(ctx, 1, 1111, 2, "172.16.33.54")
}

func TestAddExpectlogin(t *testing.T) {
	once.Do(startService)
	l := &model.Log{Mid: 1, Location: "shanghai", LocationID: 111}
	s.addExpectionLog(ctx, l)
}

func TestCommonLoc(t *testing.T) {
	once.Do(startService)
	res, err := s.commonLoc(ctx, 4780461)
	if err != nil {
		t.Errorf("err %v", err)
	}
	t.Log(res)
	for _, r := range res {
		t.Logf("res %v", r)
	}
}
func Test_OftenCheck(t *testing.T) {
	Convey("not often ip", t, func() {
		once.Do(startService)
		res, err := s.OftenCheck(ctx, 4780461, "127.0.0.1")
		So(err, ShouldEqual, nil)
		So(res.Result, ShouldEqual, false)
	})
}
