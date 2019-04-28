package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/service/main/share/conf"
	"go-common/app/service/main/share/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/share-service-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_Add(t *testing.T) {
	Convey("Add", t, func() {
		p := &model.ShareParams{
			OID: 1684013,
			MID: 100,
			TP:  1,
			IP:  "",
		}
		share, err := s.Add(context.Background(), p)
		Println(share)
		Println(err)
	})
}

func Test_Stat(t *testing.T) {
	Convey("Stat", t, func() {
		share, err := s.Stat(context.Background(), 100, 1)
		Println(share)
		Println(err)
	})
}

func Test_Stats(t *testing.T) {
	Convey("Stats", t, func() {
		s.Stats(context.Background(), []int64{1, 2, 3, 34}, 1)
	})
}
