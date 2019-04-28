package service

import (
	"fmt"
	"sync"
	"testing"

	"go-common/app/service/main/passport-game/conf"
)

var (
	once sync.Once
	s    *Service
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	s = New(conf.Conf)
}

func TestNew(t *testing.T) {
	once.Do(startService)
	if s.c.AccountURI == "" || s.c.PassportURI == "" {
		t.Errorf("conf is not correct, expected account URI and passport URI not empty but not, account URI: %s, passport URI: %s", s.c.AccountURI, s.c.PassportURI)
		t.FailNow()
	} else {
		t.Logf("s.c.AccountURI: %s, s.c.PassportURI: %s", s.c.AccountURI, s.c.PassportURI)
	}
}
