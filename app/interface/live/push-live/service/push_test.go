package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/live/push-live/model"
	"math/rand"
	"strconv"
	"testing"
)

func makeTestInitPushTask(targetID int64, uname, linkValue,
	roomTitle string, expireTime int) (task *model.ApPushTask) {
	m := &model.StartLiveMessage{
		TargetID:   targetID,
		Uname:      uname,
		LinkValue:  linkValue,
		RoomTitle:  roomTitle,
		ExpireTime: expireTime,
	}
	task = s.InitPushTask(m)
	return
}

func TestService_Push(t *testing.T) {
	initd()
	Convey("test push func", t, func() {
		// test empty mids
		targetID := rand.Int63n(100) + 1
		uname := "测试"
		linkValue := strconv.Itoa(rand.Intn(9999))
		roomTitle := "room_title"
		expireTime := rand.Intn(10000) + 1
		task := makeTestInitPushTask(targetID, uname, linkValue, roomTitle, expireTime)

		midMap := make(map[int][]int64)
		midMap[model.RelationAttention] = []int64{}

		total := s.Push(task, midMap)
		So(total, ShouldEqual, 0)
	})
}

func TestService_GetPushGroup(t *testing.T) {
	initd()
	Convey("test get group by different push type", t, func() {
		var (
			group     string
			testGroup = "test_group"
		)
		group = s.GetPushGroup(model.RelationAttention, "")
		So(group, ShouldEqual, model.AttentionGroup)

		group = s.GetPushGroup(model.RelationSpecial, "")
		So(group, ShouldEqual, model.SpecialGroup)

		group = s.GetPushGroup(rand.Intn(9999), testGroup)
		So(group, ShouldEqual, testGroup)
	})
}
