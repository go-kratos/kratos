package reply

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ChangeSubjectMid(t *testing.T) {
	Convey("ChangeSubjectMid", t, func() {
		d.ChangeSubjectMid(0, 1684013)
	})
}
