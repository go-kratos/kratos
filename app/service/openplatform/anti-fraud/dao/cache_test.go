package dao

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"testing"
	"time"

	"go-common/app/service/openplatform/anti-fraud/conf"
	"go-common/app/service/openplatform/anti-fraud/model"

	. "github.com/smartystreets/goconvey/convey"
)

const _key = "AntiFraud:BANKID_1527233672941"

var qustion = &model.ArgGetQuestion{UID: "1111", TargetItem: "1111", TargetItemType: 1, Source: 1, Platform: 1, ComponentID: 123}

func init() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(fmt.Errorf("conf.Init() error(%v)", err))
	}
	d = New(conf.Conf)
}

func TestSetex(t *testing.T) {
	Convey("TestSetex", t, func() {

		err := d.Setex(context.TODO(), "BANK_1527061216377_QUESTIONS", 1, time.Hour)

		So(err, ShouldBeNil)
	})
}

func TestSetObj(t *testing.T) {
	Convey("TestSetex", t, func() {

		obj := &model.QuestionBank{}
		reply := []byte(`{"qb_id":1527233672941,"qb_name":"wlt","cd_time":0,"max_retry_time":22,"is_deleted":0}`)

		err := json.Unmarshal(reply, obj)
		if err != nil {
			return
		}

		d.GetObj(context.TODO(), _key, obj)
		data, err := d.GetQusBankInfoCache(context.TODO(), 1527233672941)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	})
}

//
func TestGetObj(t *testing.T) {
	Convey("TestSetex", t, func() {

		data := &model.QuestionBank{}
		err := d.GetObj(context.TODO(), _key, data)
		So(err, ShouldBeNil)
	})
}

func TestGetPic(t *testing.T) {
	Convey("TestGetPic", t, func() {
		oi, _ := d.GetPic(context.TODO(), 4)
		So(oi, ShouldNotBeNil)
	})
}

func TestPushAllPic(t *testing.T) {
	Convey("TestGetPic", t, func() {
		data := &model.ArgGetQuestion{
			UID:            "111111",
			TargetItem:     "1111",
			TargetItemType: 1,
			Source:         1,
			Platform:       1,
			ComponentID:    123,
		}
		_, err := d.GetRandPic(context.TODO(), data)
		So(err, ShouldBeNil)
	})
}

func TestGetCacheQus(t *testing.T) {
	Convey("TestGetCacheQus", t, func() {
		_, err := d.GetCacheQus(context.TODO(), 1527241107344)
		So(err, ShouldBeNil)

	})
}

func TestGetQusKey(t *testing.T) {
	Convey("TestGetQusKey", t, func() {
		oi := d.GetQusKey(_keyAnsweredIds, qustion)
		So(oi, ShouldNotBeNil)
	})
}

func TestQusFetchTime(t *testing.T) {
	Convey("TestQusFetchTime", t, func() {
		oi := d.QusFetchTime(context.TODO(), qustion)
		So(oi, ShouldNotBeNil)
	})
}

func TestSetComponentId(t *testing.T) {
	Convey("TestSetComponentId", t, func() {
		err := d.SetComponentID(context.TODO(), qustion)
		So(err, ShouldBeNil)
	})
}

func TestGetComponentId(t *testing.T) {
	Convey("TestGetComponentId", t, func() {
		oi, err := d.GetComponentID(context.TODO(), qustion)
		So(oi, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestIncrComponentTimes(t *testing.T) {
	Convey("TestIncrComponentTimes", t, func() {
		err := d.IncrComponentTimes(context.TODO(), qustion)
		So(err, ShouldBeNil)
	})
}

func TestGetRandPic(t *testing.T) {
	Convey("TestGetRandPic", t, func() {
		oi, err := d.GetRandPic(context.TODO(), qustion)
		So(err, ShouldBeNil)
		So(oi, ShouldNotBeNil)
	})
}

func TestSetAnsweredID(t *testing.T) {
	Convey("TestSetAnsweredID", t, func() {
		_, err := d.RedisDo(context.TODO(), "SADD", "wlt", 1222)
		So(err, ShouldBeNil)
	})
}
