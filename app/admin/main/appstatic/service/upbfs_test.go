package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/appstatic/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	c   = context.TODO()
	now = time.Now().Unix()
)

func TestService_AddFile(t *testing.T) {
	Convey("AddFile should return without err", t, WithService(func(svf *Service) {
		// simulate resource & file data
		res := &model.Resource{}
		res.ID = 1
		resFile := &model.ResourceFile{
			ID:         0,
			Name:       "tName",
			Type:       "image/png",
			Md5:        "d41d8cd98f00b204e9800998ecf8427e",
			Size:       20,
			URL:        "http://uat-i0.hdslb.com/bfs/app-static/timg.png",
			ResourceID: 2,
			Ctime:      xtime.Time(now),
			Mtime:      xtime.Time(now),
		}
		// DB insert
		err := svf.AddFile(c, resFile, 3)
		// testing
		So(err, ShouldBeNil)
	}))
}

func TestService_Upload(t *testing.T) {
	Convey("Upload should return without err", t, WithService(func(svf *Service) {
		// simulate file content
		str2 := "Testing File Content"
		data2 := []byte(str2)
		// upload to bfs
		location, err := svf.Upload(c, "", "image/png", now, data2)
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
		fInfo, err := svf.ParseFile(data2)
		// testing
		So(err, ShouldBeNil)
		So(fInfo.Size, ShouldBeGreaterThan, 0)
		So(fInfo.Type, ShouldNotBeBlank)
		So(fInfo.Md5, ShouldNotBeBlank)
	}))
}

func TestService_TypeCheck(t *testing.T) {
	Convey("File Type Check", t, WithService(func(svf *Service) {
		So(svf.TypeCheck("application/zip"), ShouldBeTrue)
		So(svf.TypeCheck("application/xml"), ShouldBeTrue)
		So(svf.TypeCheck("application/json"), ShouldBeFalse)
	}))
}
