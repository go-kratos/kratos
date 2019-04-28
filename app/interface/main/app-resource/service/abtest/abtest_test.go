package abtest

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/service/main/resource/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-resource-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestExperiment(t *testing.T) {
	Convey("get Experiment data", t, WithService(func(s *Service) {
		res := s.Experiment(context.TODO(), model.PlatAndroid, 100)
		result, _ := json.Marshal(res)
		fmt.Printf("test Experiment (%v) \n", string(result))
		So(res, ShouldNotBeEmpty)
	}))
}
