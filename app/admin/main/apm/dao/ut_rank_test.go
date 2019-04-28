package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGitLabFace(t *testing.T) {
	username := "chenjianrong"
	convey.Convey("GitLabFace", t, func() {
		res, err := d.GitLabFace(context.Background(), username)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
