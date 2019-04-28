package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Asobinlogconsumeproc(t *testing.T) {
	once.Do(startService)
	Convey("asobinlogcommitproc running", t, func() {
		So(s.asobinlogcommitproc, ShouldBeNil)
	})
}
