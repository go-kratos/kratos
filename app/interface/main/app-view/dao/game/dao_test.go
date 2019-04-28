package game

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/game"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-view-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
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
		t.Run(tt.name, func(t *testing.T) {
			if gotD := New(tt.args.c); !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("New() = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestDao_Info(t *testing.T) {
	type args struct {
		c      context.Context
		gameID int64
		plat   int8
	}
	tests := []struct {
		name     string
		args     args
		wantInfo *game.Info
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotInfo, err := d.Info(tt.args.c, tt.args.gameID, tt.args.plat)
			So(err, ShouldEqual, tt.wantErr)
			So(gotInfo, ShouldResemble, tt.wantInfo)
		})
	}
}
