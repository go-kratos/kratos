package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/live/push-live/model"
	"math/rand"
	"strconv"
	"testing"
)

func TestService_InitPushTask(t *testing.T) {
	initd()
	Convey("should return init struct", t, func() {
		targetID = rand.Int63n(100) + 1
		uname := "测试"
		linkValue := strconv.Itoa(rand.Intn(9999))
		roomTitle := "room_title"
		expireTime := rand.Intn(10000) + 1
		task := makeTestInitPushTask(targetID, uname, linkValue, roomTitle, expireTime)

		So(task.TargetID, ShouldResemble, targetID)
		So(task.AlertTitle, ShouldResemble, uname)
		So(task.AlertBody, ShouldResemble, roomTitle)
		So(task.ExpireTime, ShouldResemble, expireTime)
		So(task.LinkValue, ShouldResemble, linkValue)
	})
}

func TestDao_GetSourceByTypes(t *testing.T) {
	initd()
	Convey("Get mid_source by different types", t, func() {
		types := []string{model.StrategySwitch, model.StrategyFans, model.StrategySpecial, model.StrategySwitchSpecial}
		length := len(types)
		currentX := rand.Intn(length)
		currentY := rand.Intn(length)

		var currentTypes []string
		if currentX >= currentY {
			currentTypes = types[currentY:currentX]
		} else {
			currentTypes = types[currentX:currentY]
		}
		midSource := s.getSourceByTypes(currentTypes)

		So(midSource, ShouldBeGreaterThanOrEqualTo, 0)
		So(midSource, ShouldBeLessThanOrEqualTo, 15)
	})
}

func TestService_GetFansBySwitch(t *testing.T) {
	initd()
	Convey("should find some fans id by given target id", t, func() {
		targetID = 27515316
		fans, fansSP, err := s.GetFansBySwitch(context.Background(), targetID)

		So(len(fans), ShouldBeGreaterThan, 0)
		So(len(fansSP), ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	})
}

func TestService_GetFansBySwitchAndSpecial(t *testing.T) {
	initd()
	Convey("should find some fans id by given target id", t, func() {
		targetID = 27515316
		fans, fansSP, err := s.GetFansBySwitchAndSpecial(context.Background(), targetID)

		So(len(fans), ShouldEqual, 0)
		So(len(fansSP), ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	})
}

func TestService_GetMids(t *testing.T) {
	initd()
	Convey("should find some fans id by given target id", t, func() {
		targetID = 27515316
		uname := "测试"
		linkValue := strconv.Itoa(rand.Intn(9999))
		roomTitle := "room_title"
		expireTime := rand.Intn(10000) + 1
		task := makeTestInitPushTask(targetID, uname, linkValue, roomTitle, expireTime)
		types := []string{"Switch", "Special"}
		s.pushTypes = types
		midMap := s.GetMids(context.Background(), task)

		for _, list := range midMap {
			So(len(list), ShouldBeGreaterThan, 0)
		}
	})
}
