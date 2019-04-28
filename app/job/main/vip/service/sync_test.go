package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SyncAllUser(t *testing.T) {
	Convey("Test_SyncAllUser", t, func() {
		s.SyncAllUser(context.Background())
	})
}
