package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestSubject d.Subject() unit test.
func TestSubject(t *testing.T) {
	Convey("test dao.Subject", t, func() {
		s, err := testDao.Subject(context.TODO(), 1, 1221)
		So(err, ShouldBeNil)
		So(s, ShouldNotBeNil)
	})
}
