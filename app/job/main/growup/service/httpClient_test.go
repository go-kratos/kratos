package service

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_HTTPClient(t *testing.T) {
	Convey("growup-job HttpClient", t, WithService(func(s *Service) {
		var (
			method  = "GET"
			url     = ""
			params  = make(map[string]string)
			nowTime = time.Now().Add(-24 * time.Hour).Unix()
		)
		res, err := s.HTTPClient(method, url, params, nowTime)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
