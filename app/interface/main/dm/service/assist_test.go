package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAssist(t *testing.T) {
	var (
		upID int64
		isUp bool
		err  error
	)
	Convey("test assist not up", t, func() {
		upID, isUp, err = svr.assist(context.Background(), int64(27515406), int64(10097377))
		So(err, ShouldBeNil)
		So(upID, ShouldBeGreaterThan, 1)
		So(isUp, ShouldBeFalse)
	})

	Convey("test assist is up", t, func() {
		upID, isUp, err = svr.assist(context.Background(), int64(27515256), int64(10097377))
		So(err, ShouldBeNil)
		So(upID, ShouldBeGreaterThan, 1)
		So(isUp, ShouldBeTrue)
	})
}

func TestAssistBanned(t *testing.T) {
	Convey("test assist banned", t, func() {
		err := svr.AssistBanned(context.TODO(), 27515256, 9967830, []int64{719926094, 719926092})
		So(err, ShouldBeNil)
	})
}

func TestAssistUptBanned(t *testing.T) {
	Convey("test assist upt banned", t, func() {
		err := svr.AssistUptBanned(context.TODO(), 27515256, "hash1", 0)
		So(err, ShouldBeNil)
		err = svr.AssistUptBanned(context.TODO(), 27515256, "hash1", 1)
		So(err, ShouldBeNil)
	})
}

func TestAssistDelBanned2(t *testing.T) {
	Convey("test assist banned2", t, func() {
		err := svr.AssistDelBanned2(context.TODO(), 27515256, 10097377, []string{"hash1", "hash2"})
		So(err, ShouldBeNil)

		err = svr.AssistUptBanned(context.TODO(), 27515256, "hash1", 1)
		So(err, ShouldBeNil)

		err = svr.AssistUptBanned(context.TODO(), 27515256, "hash2", 1)
		So(err, ShouldBeNil)
	})
}

func TestAssistBannedUsers(t *testing.T) {
	Convey("test assist banned users", t, func() {
		rs, err := svr.AssistBannedUsers(context.TODO(), 27515256, 10097377)
		So(err, ShouldBeNil)
		So(rs, ShouldBeGreaterThan, 1)
	})
}

func TestAssistDelete(t *testing.T) {
	Convey("test assist delete dm", t, func() {
		err := svr.AssistDeleteDM(context.TODO(), 27515256, 10108163, []int64{719925514, 719925516})
		So(err, ShouldBeNil)
	})
}
