package tab

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-card/model/card/operate"
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

func TestDao_Menus(t *testing.T) {
	type args struct {
		c context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantMenus []*operate.Menu
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMenus, err := d.Menus(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.Menus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMenus, tt.wantMenus) {
				t.Errorf("Dao.Menus() = %v, want %v", gotMenus, tt.wantMenus)
			}
		})
	}
}

func TestDao_Actives(t *testing.T) {
	type args struct {
		c context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantAcs []*operate.Active
		wantErr error
	}{}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotAcs, err := d.Actives(tt.args.c)
			So(err, ShouldEqual, tt.wantErr)
			So(gotAcs, ShouldResemble, tt.wantAcs)
		})
	}
}
