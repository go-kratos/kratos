package dynamic

import (
	"context"
	"encoding/json"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

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

func TestDao_dynamicSrv(t *testing.T) {
	type args struct {
		c       context.Context
		reqPath string
		params  string
	}
	tests := []struct {
		name    string
		args    args
		wantRes json.RawMessage
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := d.dynamicSrv(tt.args.c, tt.args.reqPath, tt.args.params)
			So(gotRes, ShouldEqual, tt.wantRes)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestDao_DynamicHistory(t *testing.T) {
	type args struct {
		c      context.Context
		params string
	}
	tests := []struct {
		name    string
		d       *Dao
		args    args
		wantRes json.RawMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.d.DynamicHistory(tt.args.c, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.DynamicHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Dao.DynamicHistory() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestDao_DynamicCount(t *testing.T) {
	type args struct {
		c      context.Context
		params string
	}
	tests := []struct {
		name    string
		d       *Dao
		args    args
		wantRes json.RawMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.d.DynamicCount(tt.args.c, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.DynamicCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Dao.DynamicCount() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestDao_DynamicNew(t *testing.T) {
	type args struct {
		c      context.Context
		params string
	}
	tests := []struct {
		name    string
		d       *Dao
		args    args
		wantRes json.RawMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.d.DynamicNew(tt.args.c, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.DynamicNew() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Dao.DynamicNew() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
