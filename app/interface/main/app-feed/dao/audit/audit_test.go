package audit

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dao *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-feed-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	dao = New(conf.Conf)
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

func TestDao_Audits(t *testing.T) {
	type args struct {
		c context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantRes map[string]map[int]struct{}
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotRes, err := dao.Audits(tt.args.c)
			So(gotRes, ShouldEqual, tt.wantRes)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestDao_PingDB(t *testing.T) {
	type args struct {
		c context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dao.PingDB(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Dao.PingDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDao_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao.Close()
		})
	}
}
