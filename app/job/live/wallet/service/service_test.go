package Service

import (
	"context"
	"go-common/app/job/live/wallet/conf"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	conf.ConfPath = "../cmd/live-wallet-test.toml"
	once.Do(startService)
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	Convey("Ping", t, func() {
		once.Do(startService)
		time.Sleep(time.Second)
		s.Ping(context.TODO())
	})
}

func TestClose(t *testing.T) {
	Convey("Close", t, func() {
		once.Do(startService)
		time.Sleep(time.Second)
		s.Close()
	})
}

func TestWait(t *testing.T) {
	Convey("Wait", t, func() {
		once.Do(startService)
		time.Sleep(time.Second)
		s.Wait()
	})
}
