package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/service/main/card/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	c = context.TODO()
	s *Service
)

func init() {
	var (
		err error
	)
	flag.Set("conf", "../cmd/test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	c = context.Background()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

// go test  -test.v -test.run TestCard
func TestCard(t *testing.T) {
	Convey("TestCard ", t, func() {
		card := s.Card(c, 1)
		t.Logf("v(%v)", card)
		So(card, ShouldNotBeEmpty)
	})
}

// go test  -test.v -test.run TestCardHots
func TestCardHots(t *testing.T) {
	Convey("TestCardHots ", t, func() {
		card := s.CardHots(c)
		t.Logf("v(%v)", card)
		So(card, ShouldNotBeEmpty)
	})
}

// go test  -test.v -test.run TestCardsByGid
func TestCardsByGid(t *testing.T) {
	Convey("CardsByGid ", t, func() {
		card := s.CardsByGid(c, 1)
		t.Logf("v(%v)", card)
		So(card, ShouldNotBeEmpty)
	})
}
