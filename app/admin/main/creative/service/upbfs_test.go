package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	c   = context.TODO()
	now = time.Now().Unix()
)

func TestService_Upload(t *testing.T) {
	Convey("Upload should return without err", t, WithService(func(svf *Service) {
		// simulate file content
		str2 := "Testing File Content"
		data2 := []byte(str2)
		// upload to bfs
		location, err := svf.Upload(c, "tName", "image/png", now, data2)
		// testing
		So(err, ShouldBeNil)
		So(location, ShouldNotBeNil)
		fmt.Println(location)
	}))
}

func TestService_ParseFile(t *testing.T) {
	Convey("ParseFile can get file info", t, WithService(func(svf *Service) {
		// simulate file content
		str2 := "Testing File Content"
		data2 := []byte(str2)
		// upload to bfs
		fInfo, err := svf.ParseFile(c, data2)
		// testing
		So(err, ShouldBeNil)
		So(fInfo.Size, ShouldBeGreaterThan, 0)
		So(fInfo.Type, ShouldNotBeBlank)
		So(fInfo.Md5, ShouldNotBeBlank)
	}))
}
