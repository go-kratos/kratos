package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/admin/main/vip/conf"
	"go-common/app/admin/main/vip/model"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	//once sync.Once
	s *Service
	c = context.TODO()
)

func init() {
	flag.Set("conf", "../cmd/vip-admin-test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func Test_ScanUserInfo(t *testing.T) {
	Convey("should return true err == nil", t, func() {

	})
}

func TestService_FrozenCode(t *testing.T) {
	Convey("frozenCode", t, func() {
		codeID := 12
		status := 1
		err := s.FrozenCode(context.TODO(), int64(codeID), int8(status))
		So(err, ShouldBeNil)
	})
}
func TestService_FrozenBatchCode(t *testing.T) {
	Convey("frozen batch code", t, func() {
		batchCodeID := 12
		status := 2
		err := s.FrozenBatchCode(context.TODO(), int64(batchCodeID), int8(status))
		So(err, ShouldBeNil)
	})
}
func TestService_BusinessInfo(t *testing.T) {
	Convey("businessInfo", t, func() {
		id := 12
		_, err := s.BusinessInfo(context.TODO(), id)
		So(err, ShouldBeNil)
	})
}
func TestService_AllVersion(t *testing.T) {
	Convey("version", t, func() {
		_, err := s.AllVersion(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestService_SelCode(t *testing.T) {
	Convey("sel code", t, func() {
		arg := new(model.ArgCode)
		arg.ID = 12
		cursor := 1
		ps := 20
		res, _, _, err := s.SelCode(context.TODO(), arg, "zhaozhihao", int64(cursor), ps)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestService_VipInfo(t *testing.T) {
	Convey("vip info", t, func() {
		mid := 1233
		_, err := s.VipInfo(context.TODO(), int64(mid))
		So(err, ShouldBeNil)
	})
}

func TestService_BatchInfoOfPool(t *testing.T) {
	Convey("batchn info ", t, func() {
		poolID := 12
		_, err := s.BatchInfoOfPool(context.TODO(), poolID)
		So(err, ShouldBeNil)
	})
}

func TestService_BusinessList(t *testing.T) {
	Convey("business list", t, func() {
		res, _, err := s.BusinessList(context.TODO(), 1, 20, -1)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
