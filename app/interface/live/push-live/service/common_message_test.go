package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/live/push-live/dao"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/cache/redis"
	"math/rand"
	"strconv"
	"testing"
)

func makeTestCommonPushTask(title, body, linkValue, group string, business, expireTime int) (task *model.ApPushTask) {
	m := &model.LiveCommonMessage{}
	m.MsgContent = model.LiveCommonMessageContent{
		Business:   business,
		Group:      group,
		Mids:       "",
		AlertTitle: title,
		AlertBody:  body,
		LinkValue:  linkValue,
		ExpireTime: expireTime,
	}
	task = s.InitCommonTask(m)
	return
}

func TestService_InitCommonTask(t *testing.T) {
	initd()
	Convey("should return init struct", t, func() {
		title := "room_title"
		body := "测试"
		group := "group"
		linkValue := strconv.Itoa(rand.Intn(9999))
		expireTime := rand.Intn(10000) + 1
		business := rand.Intn(9999)
		task := makeTestCommonPushTask(title, body, linkValue, group, business, expireTime)

		So(task.AlertTitle, ShouldResemble, title)
		So(task.AlertBody, ShouldResemble, body)
		So(task.ExpireTime, ShouldResemble, expireTime)
		So(task.LinkValue, ShouldResemble, linkValue)
		So(task.MidSource, ShouldEqual, business)
		So(task.Group, ShouldEqual, group)
	})
}

func TestService_setPushInterval(t *testing.T) {
	initd()
	Convey("test setPushInterval", t, func() {
		var (
			resTotal int
			total    int
			business int
			task     *model.ApPushTask
			mids     []int64
			err      error
		)
		Convey("test business will not exec logic", func() {
			business = rand.Intn(100)
			task = &model.ApPushTask{}
			total = 10
			mids = makeMids(total)
			resTotal, err = s.setPushInterval(business, rand.Int31(), mids, task)
			So(err, ShouldBeNil)
			So(resTotal, ShouldEqual, 0)
		})

		Convey("test business will exec logic", func() {
			var conn redis.Conn
			business = 111
			task = &model.ApPushTask{
				LinkValue: "test",
			}
			total = 10
			mids = makeMids(total)
			resTotal, err = s.setPushInterval(business, 300, mids, task)
			So(err, ShouldBeNil)
			So(resTotal, ShouldEqual, total)
			// clean
			conn, err = redis.Dial(s.c.Redis.PushInterval.Proto, s.c.Redis.PushInterval.Addr, s.dao.RedisOption()...)
			So(err, ShouldBeNil)
			for _, mid := range mids {
				key := dao.GetIntervalKey(mid)
				conn.Do("DEL", key)
			}
			conn.Close()
		})
	})
}

func makeMids(total int) []int64 {
	mids := make([]int64, 0, total)
	for i := 0; i < total; i++ {
		mids = append(mids, rand.Int63())
	}
	return mids
}
