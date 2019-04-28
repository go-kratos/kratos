package live

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-feed/conf"
	model "go-common/app/interface/main/app-feed/model/live"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-feed-test.toml")
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

func TestDao_AppMRoom(t *testing.T) {
	type args struct {
		c       context.Context
		roomids []int64
	}
	tests := []struct {
		name    string
		args    args
		wantRs  map[int64]*live.Room
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotRs, err := d.AppMRoom(tt.args.c, tt.args.roomids)
			So(gotRs, ShouldEqual, tt.wantRs)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestDao_FeedList(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name      string
		args      args
		wantFs    []*model.Feed
		wantCount int
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotFs, gotCount, err := d.FeedList(tt.args.c, tt.args.mid, tt.args.pn, tt.args.ps)
			So(gotFs, ShouldEqual, tt.wantFs)
			So(gotCount, ShouldEqual, tt.wantCount)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}
