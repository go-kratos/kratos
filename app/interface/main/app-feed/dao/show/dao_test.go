package show

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-card/model/card/show"
	"go-common/app/interface/main/app-feed/conf"

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

func TestDao_Card(t *testing.T) {
	type args struct {
		c   context.Context
		ids []int64
	}
	tests := []struct {
		name    string
		d       *Dao
		args    args
		wantRs  map[int64]*show.Shopping
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotRs, err := tt.d.Card(tt.args.c, tt.args.ids)
			So(gotRs, ShouldEqual, tt.wantRs)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}
