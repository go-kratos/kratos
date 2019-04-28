package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	xtime "go-common/library/time"

	"go-common/app/admin/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_GetPushData(t *testing.T) {
	Convey("getpushData should be nil", t, func() {
		id := 1
		res, err := s.GetPushData(context.TODO(), int64(id))
		bytes, _ := json.Marshal(res)
		t.Logf("%+v", string(bytes))
		So(err, ShouldBeNil)
	})
}

func TestService_SavePushData(t *testing.T) {
	Convey("save push data should be nil", t, func() {
		arg := new(model.VipPushData)
		arg.ID = 1
		arg.GroupName = "test01"
		arg.Title = "title"
		arg.Content = "content"
		arg.Platform = "[{\"name\":\"Android\",\"condition\":\"=\",\"build\":1},{\"name\":\"iPhone\",\"condition\":\"<=\",\"build\":2},{\"name\":\"iPad\",\"condition\":\"=\",\"build\":3}]"
		arg.LinkType = 10
		arg.ExpiredDayStart = -1
		arg.ExpiredDayEnd = 10
		arg.EffectStartDate = xtime.Time(time.Now().Unix())
		arg.EffectEndDate = xtime.Time(time.Now().AddDate(0, 0, 7).Unix())
		arg.PushStartTime = "18:00:00"
		arg.PushEndTime = "20:00:00"
		err := s.SavePushData(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestService_PushDatas(t *testing.T) {
	Convey("push data should be nil", t, func() {
		arg := new(model.ArgPushData)

		arg.Status = 0
		arg.ProgressStatus = 0
		res, count, err := s.PushDatas(context.TODO(), arg)
		bytes, _ := json.Marshal(res)
		t.Logf("res(%+v) count(%v) ", string(bytes), count)
		So(err, ShouldBeNil)
	})
}
