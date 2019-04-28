package history

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/tv/conf"

	"context"

	"encoding/json"
	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(srv)
	}
}

func TestService_GetHistory(t *testing.T) {
	Convey("TestService_GetHistory", t, WithService(func(s *Service) {
		cont, err := s.GetHistory(context.Background(), int64(27515401))
		So(err, ShouldBeNil)
		data, _ := json.Marshal(cont)
		fmt.Println(string(data))
	}))
}
