package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/interface/main/videoup/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Ping(t *testing.T) {
	var (
		c   = context.Background()
		err error
	)
	Convey("Ping", t, WithService(func(s *Service) {
		err = s.Ping(c)
		So(err, ShouldBeNil)
	}))
}

func Test_Close(t *testing.T) {
	Convey("Close", t, WithService(func(s *Service) {
		s.Close()
	}))
}
func Test_WebAdd(t *testing.T) {
	var (
		c    = context.Background()
		err  error
		MID  = int64(27515256)
		body = `{"copyright":1,"cover":"","title":"test","tid":130,"tag":"音乐选集","no_reprint":1,"upos":0,"lang":"zh-CN","mission_id":0,"porder":{},"desc":"123","dynamic":"123","videos":[{"desc":"","filename":"g180126072jadys8fuaz74u18hkxwvnf","title":""}]}
		`
		aid int64
	)
	var ap = &archive.ArcParam{}
	if err = json.Unmarshal([]byte(body), ap); err != nil {
		return
	}
	Convey("webAdd", t, WithService(func(s *Service) {
		aid, err = s.WebAdd(c, MID, ap, true)
		So(err, ShouldBeNil)
		So(aid, ShouldNotBeNil)
	}))
}

func Test_WebEdit(t *testing.T) {
	var (
		c    = context.Background()
		err  error
		MID  = int64(27515256)
		body = `{"copyright":1,"cover":"","title":"test","tid":130,"tag":"音乐选集","no_reprint":1,"upos":0,"lang":"zh-CN","mission_id":0,"porder":{},"desc":"123","dynamic":"123","videos":[{"desc":"","filename":"g180126072jadys8fuaz74u18hkxwvnf","title":""}]}
		`
	)
	var ap = &archive.ArcParam{}
	if err = json.Unmarshal([]byte(body), ap); err != nil {
		return
	}
	Convey("webEdit", t, WithService(func(s *Service) {
		err = s.WebEdit(c, ap, MID)
		So(err, ShouldBeNil)
	}))
}
