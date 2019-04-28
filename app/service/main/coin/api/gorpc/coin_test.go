package coin

import (
	"context"
	"testing"
	"time"

	coin "go-common/app/service/main/coin/model"
)

const (
	mid    = 23675773
	aid    = 1
	realIP = "127.0.0.1"
)

func TestCoin(t *testing.T) {
	s := New(nil)
	time.Sleep(1 * time.Second)

	// coin
	testAddCoins(t, s)
	testArchiveUserCoins(t, s)
}

func testAddCoins(t *testing.T, s *Service) {
	arg := coin.ArgAddCoin{Mid: mid, Aid: aid, Multiply: 1, RealIP: realIP}
	if err := s.AddCoins(context.TODO(), &arg); err != nil {
		t.Logf("call.AddCoins error(%v)", err)
	}
}

func testArchiveUserCoins(t *testing.T, s *Service) {
	arg := coin.ArgCoinInfo{Mid: mid, Aid: aid, RealIP: realIP}
	if res, err := s.ArchiveUserCoins(context.TODO(), &arg); err != nil && res != nil {
		t.Logf("call.ArchiveUserCoins error(%v)", err)
	}
}
