package view

import (
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ViewInfoc(t *testing.T) {
	Convey("ViewInfoc", t, func() {
		s.ViewInfoc(0, 0, "test", "0", "", "", "", "", "", "", time.Now(), errors.New("test"), 1, "", "")
	})
}

func Test_RelateInfoc(t *testing.T) {
	Convey("RelateInfoc", t, func() {
		s.RelateInfoc(0, 0, 0, "", "", "", "", "", "", "", "", "", nil, time.Now(), 0)
	})
}
