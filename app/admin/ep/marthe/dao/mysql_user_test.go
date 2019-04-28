package dao

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/ep/marthe/model"
)

var (
	testTime = time.Now().Format("2006_01_02_15_04_05")
	testUser = model.User{
		Name:  testTime,
		EMail: testTime + "@bilibili.com"}
)

func Test_User(t *testing.T) {
	Convey("test CreateUser", t, func() {
		err := d.CreateUser(&testUser)
		So(err, ShouldBeNil)
	})

	Convey("find user by user name", t, func() {
		userInDb, err := d.FindUserByUserName(testUser.Name)
		So(userInDb.EMail, ShouldEqual, testUser.EMail)
		So(err, ShouldBeNil)
	})

	Convey("find user by id", t, func() {
		userID := testUser.ID
		userInDb, err := d.FindUserByID(userID)
		So(userInDb.EMail, ShouldEqual, testUser.EMail)
		So(err, ShouldBeNil)
	})

	Convey("delete user", t, func() {
		err := d.DelUser(&testUser)
		So(err, ShouldBeNil)
	})

}
