package dao

import (
	"strconv"
	"testing"
	"time"

	"go-common/app/admin/ep/marthe/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tmpIssueNoStr2 = strconv.FormatInt(time.Now().Unix(), 10)

	buglyVersion = &model.BuglyVersion{
		Version:        "Version" + tmpIssueNoStr2,
		BuglyProjectID: 1,
		Action:         model.BuglyVersionActionDisable,
		TaskStatus:     1,
		UpdateBy:       "fengyifeng",
	}
)

func Test_Bugly_Version(t *testing.T) {
	Convey("test insert bugly Version", t, func() {
		err := d.InsertBuglyVersion(buglyVersion)
		So(err, ShouldBeNil)
	})

	Convey("test update bugly Version", t, func() {
		buglyVersion.Version = "update" + tmpIssueNoStr2
		err := d.UpdateBuglyVersion(buglyVersion)
		So(err, ShouldBeNil)
	})

	Convey("test Query Bugly Version By Version", t, func() {
		tmpBuglyVersion, err := d.QueryBuglyVersionByVersion(buglyVersion.Version)
		So(err, ShouldBeNil)
		So(tmpBuglyVersion.Version, ShouldEqual, buglyVersion.Version)
	})

	Convey("test Query Bugly Version By Id", t, func() {
		tmpBuglyVersion, err := d.QueryBuglyVersion(buglyVersion.ID)
		So(err, ShouldBeNil)
		So(tmpBuglyVersion.Version, ShouldEqual, buglyVersion.Version)
	})

}
