package tag

import (
	"context"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-interface/conf"
	tagmdl "go-common/app/interface/main/app-interface/model/tag"
	tagrpc "go-common/app/interface/main/tag/rpc/client"

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
			So(gotD, ShouldEqual, tt.wantD)
		})
	}
}

func TestDao_ArcTags(t *testing.T) {
	type fields struct {
		tagRPC *tagrpc.Service
	}
	type args struct {
		c   context.Context
		aid int64
		mid int64
		ip  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantTs  []*tagmdl.Tag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				tagRPC: tt.fields.tagRPC,
			}
			gotTs, err := d.ArcTags(tt.args.c, tt.args.aid, tt.args.mid, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.ArcTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTs, tt.wantTs) {
				t.Errorf("Dao.ArcTags() = %v, want %v", gotTs, tt.wantTs)
			}
		})
	}
}
