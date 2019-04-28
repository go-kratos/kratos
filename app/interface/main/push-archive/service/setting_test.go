package service

import (
	"context"
	ht "net/http"
	"sync"
	"testing"

	"go-common/app/interface/main/push-archive/conf"
	"go-common/app/interface/main/push-archive/model"
	pb "go-common/app/service/main/push/api/grpc/v1"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/smartystreets/goconvey/convey"
)

func Test_getsetting(t *testing.T) {
	initd()

	mid := int64(11111111)
	s.SetSetting(context.TODO(), mid, nil)
	convey.Convey("获取默认的用户开关设置", t, func() {
		st, err := s.Setting(context.TODO(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(st.Type, convey.ShouldEqual, model.PushTypeSpecial)
	})
}

func Test_setsetting(t *testing.T) {
	initd()

	mid := int64(11111111)
	convey.Convey("存储用户开关设置", t, func() {
		err := s.SetSetting(context.TODO(), mid, &model.Setting{Type: model.PushTypeAttention})
		convey.So(err, convey.ShouldBeNil)
	})
}

func Test_ps(t *testing.T) {
	initd()

	url := "http://127.0.0.1:7031/x/push-archive/setting/get?access_key=848657e31639317257ab274741c6ec7d"
	client := bm.NewClient(conf.Conf.HTTPClient)
	type _response struct {
		Code int           `json:"code"`
		Data model.Setting `json:"data"`
	}

	wg := sync.WaitGroup{}
	all := 10
	busyIn := 0
	errIn := 0
	normalIn := 0
	mx := sync.Mutex{}
	for i := 0; i < all; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			defer mx.Unlock()

			req, err := ht.NewRequest("GET", url, nil)
			if err != nil {
				mx.Lock()
				errIn++
				return
			}

			r := &_response{}
			err = client.Do(context.TODO(), req, r)
			mx.Lock()
			if err != nil {
				errIn++
				return
			}

			if r.Code == ecode.ServiceUnavailable.Code() {
				busyIn++
				return
			}
			normalIn++
		}(i)
	}

	wg.Wait()
	convey.Convey("并行访问推送开关，有频率限制", t, func() {
		convey.So(errIn, convey.ShouldEqual, 0)
		convey.So(busyIn, convey.ShouldBeGreaterThan, 0)
		convey.So(errIn+busyIn+normalIn, convey.ShouldEqual, all)
	})
	t.Logf("the errin(%d), busyin(%d), normalin(%d)", errIn, busyIn, normalIn)
}

func Test_pushrpc(t *testing.T) {
	initd()

	convey.Convey("推送平台rpc设置推送开关", t, func() {
		set := &pb.SetSettingRequest{
			Mid:   int64(111111111),
			Type:  1,
			Value: 3,
		}
		res, err := s.pushRPC.Setting(context.TODO(), &pb.SettingRequest{Mid: set.Mid})
		convey.So(err, convey.ShouldBeNil)
		t.Logf("setting(%+v)", res)

		_, err = s.pushRPC.SetSetting(context.TODO(), set)
		convey.So(err, convey.ShouldBeNil)

		res, err = s.pushRPC.Setting(context.TODO(), &pb.SettingRequest{Mid: set.Mid})
		convey.So(err, convey.ShouldBeNil)
		t.Logf("setting(%+v)", res)
	})

}
