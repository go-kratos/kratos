package shop

import (
	"context"
	"testing"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/shop"
	httpx "go-common/library/net/http/blademaster"

	. "github.com/smartystreets/goconvey/convey"
)

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

func TestDao_Info(t *testing.T) {
	type fields struct {
		client *httpx.Client
		info   string
	}
	type args struct {
		c       context.Context
		mid     int64
		mobiApp string
		device  string
		build   int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantInfo *shop.Info
		wantErr  error
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			d := &Dao{
				client: tt.fields.client,
				info:   tt.fields.info,
			}
			gotInfo, err := d.Info(tt.args.c, tt.args.mid, tt.args.mobiApp, tt.args.device, tt.args.build)
			So(gotInfo, ShouldResemble, tt.wantInfo)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}
