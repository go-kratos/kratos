package view

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PushFail(t *testing.T) {
	Convey("PushFail", t, func() {
		d.PushFail(context.TODO(), nil)
	})
}

func Test_PopFail(t *testing.T) {
	Convey("PopFail", t, func() {
		d.PopFail(context.TODO())
	})
}

func Test_PingRedis(t *testing.T) {
	Convey("PingRedis", t, func() {
		d.PingRedis(context.TODO())
	})
}
