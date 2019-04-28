package dao

import (
	"context"
	"strconv"
	"testing"
	"time"

	"go-common/app/job/main/passport/model"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_AddLoginLogHBase(t *testing.T) {
	ts := time.Now().Unix()
	mid := int64(88888970)
	ipN := InetAtoN("127.0.0.1")
	for i := 0; i < 6; i++ {
		m := &model.LoginLog{
			Mid:       mid,
			Timestamp: ts + int64(i),
			LoginIP:   int64(ipN),
			Type:      1,
			Server:    strconv.FormatInt(int64(i), 10),
		}
		if err := d.AddLoginLogHBase(context.TODO(), m); err != nil {
			t.Logf("dao.AddLoginLogHBase(%+v) error(%v)", m, err)
			t.FailNow()
		}
	}
}

func TestDao_SendLoginLogMsgs(t *testing.T) {
	convey.Convey("SetToken", t, func(ctx convey.C) {
		dsPubConf := &databus.Config{
			Key:          "0QEO9F8JuuIxZzNDvklH",
			Secret:       "0QEO9F8JuuIxZzNDvklI",
			Group:        "PassportLog-Login-P",
			Topic:        "PassportLog-T",
			Action:       "pub",
			Name:         "databus",
			Proto:        "tcp",
			Addr:         "172.16.33.158:6205",
			Active:       1,
			Idle:         1,
			DialTimeout:  xtime.Duration(time.Second),
			WriteTimeout: xtime.Duration(time.Second),
			ReadTimeout:  xtime.Duration(time.Second),
			IdleTimeout:  xtime.Duration(time.Minute),
		}
		dsPub := databus.New(dsPubConf)
		defer dsPub.Close()
		ts := time.Now().Unix()
		mid := int64(88888970)
		ipN := InetAtoN("127.0.0.1")
		for i := 1; i <= 1; i++ {
			v := &model.LoginLog{
				Mid:       mid,
				Timestamp: ts + int64(i),
				LoginIP:   int64(ipN),
				Type:      1,
				Server:    strconv.FormatInt(int64(i), 10),
			}
			k := dsPubConf.Topic + strconv.FormatInt(mid, 10)
			if err := dsPub.Send(context.TODO(), k, v); err != nil {
				t.Errorf("failed to send login log databus message, dsPub.Send(%v, %v) error(%v)", k, v, err)
				t.FailNow()
			}
		}
	})
}

func TestIntCast(t *testing.T) {
	convey.Convey("when ts diff by delta, the uint32(_int64 - ts) diff should be the same", t, func(ctx convey.C) {
		ts := time.Now().Unix()
		delta := int64(10)
		a := uint32(_int64Max - ts)
		b := uint32(_int64Max - ts - delta)
		ctx.So(a-b, convey.ShouldEqual, delta)
	})
}
