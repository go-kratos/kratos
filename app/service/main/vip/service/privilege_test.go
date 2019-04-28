package service

import (
	"testing"

	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

//  go test  -test.v -test.run TestServicePrivilegesBySid
func TestServicePrivilegesBySid(t *testing.T) {
	Convey("TestServicePrivilegesBySid test ", t, func() {
		mvps, err := s.PrivilegesBySid(c, &model.ArgPrivilegeBySid{
			Sid:      7,
			Platform: "ios",
		})
		t.Logf("res(%+v)", mvps)
		for _, v := range mvps.List {
			t.Logf("r(%+v)", v)
		}
		So(err, ShouldBeNil)
	})
}

//  go test  -test.v -test.run TestServicePrivilegesByType
func TestServicePrivilegesByType(t *testing.T) {
	Convey("TestServicePrivilegesByType test ", t, func() {
		mvps, err := s.PrivilegesByType(c, &model.ArgPrivilegeDetail{
			Type:     model.AllPrivilege,
			Platform: "ios",
		})
		for _, v := range mvps {
			t.Logf("res(+%v)", v)
		}
		So(err, ShouldBeNil)
	})
}

func TestServicePrivilegesLangauage(t *testing.T) {
	Convey("TestServicePrivilegesLangauage", t, func() {
		lt := s.getPrivileges("test")
		So(lt, ShouldNotBeNil)
	})
}
