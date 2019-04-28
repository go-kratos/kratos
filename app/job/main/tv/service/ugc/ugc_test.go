package ugc

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"encoding/json"
	"go-common/app/job/main/tv/conf"
	"go-common/app/job/main/tv/model/ugc"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(srv)
	}
}

func TestService_ArcHandle(t *testing.T) {
	Convey("TestService_ArcHandle", t, WithService(func(s *Service) {
		msg := []byte(`{"action":"update","table":"archive","new":{"id":0,"aid":10110186,"mid":27515615,"typeid":76,"videos":1,"title":"xxx","cover":"http://i0.hdslb.com/bfs/archive/b5beb958f94f5deb6d3ba8775c4e81d2cf0f4bc1.jpg","content":"万年不变小电视","duration":10,"attribute":2113536,"copyright":1,"access":0,"pubtime":"2018-06-12 19:40:23","ctime":"2018-06-12 19:40:28","mtime":"",
"state":0,"mission_id":0,"order_id":0,"redirect_url":"","forward":0,"dynamic":""},"old":{"id":0,"aid":10110186,"mid":27515615,"typeid":76,"videos":1,"title":"XXXX","cover":"http://i0.hdslb.com/bfs/archive/b5beb958f94f5deb6d3ba8775c4e81d2cf0f4bc1.jpg","content":"万年不变小电视","duration":10,"attribute":2113536,"copyright":1,"access":0,"pubtime":"2018-06-12 19:40:23","ctime":"2018-06-12 19:40:28","mtime":"",
"state":0,"mission_id":0,"order_id":0,"redirect_url":"","forward":0,"dynamic":""}}`)
		var ms = &ugc.ArcMsg{}
		err := json.Unmarshal(msg, ms)
		s.ArcHandle(ms)
		So(err, ShouldBeNil)
	}))
}
