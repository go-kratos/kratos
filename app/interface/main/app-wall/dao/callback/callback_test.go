package callback

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-wall/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-wall-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestNew(t *testing.T) {
	type args struct {
		c *conf.Config
	}
	tests := []struct {
		name  string
		args  args
		wantD *Dao
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotD := New(tt.args.c)
			So(gotD, ShouldResemble, tt.wantD)
		})
	}
}

func TestDao_GdtCallback(t *testing.T) {
	type args struct {
		c       context.Context
		appID   string
		appType string
		aderID  string
		idfa    string
		cb      string
		now     time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			err := d.GdtCallback(tt.args.c, tt.args.appID, tt.args.appType, tt.args.aderID, tt.args.idfa, tt.args.cb, tt.args.now)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestDao_ShikeCallback(t *testing.T) {
	type args struct {
		c    context.Context
		idfa string
		cb   string
		now  time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			err := d.ShikeCallback(tt.args.c, tt.args.idfa, tt.args.cb, tt.args.now)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestDao_DontinCallback(t *testing.T) {
	type args struct {
		c       context.Context
		idfa    string
		clickid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			err := d.DontinCallback(tt.args.c, tt.args.idfa, tt.args.clickid)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestDao_ToutiaoCallback(t *testing.T) {
	type args struct {
		c         context.Context
		cb        string
		eventType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			err := d.ToutiaoCallback(tt.args.c, tt.args.cb, tt.args.eventType)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}
