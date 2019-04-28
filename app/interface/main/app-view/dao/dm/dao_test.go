package dm

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-view/conf"
	dm "go-common/app/interface/main/dm2/model"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../../cmd/app-view-test.toml")
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
		t.Run(tt.name, func(t *testing.T) {
			if gotD := New(tt.args.c); !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("New() = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestDao_SubjectInfos(t *testing.T) {
	type args struct {
		c    context.Context
		typ  int32
		plat int8
		oids []int64
	}
	tests := []struct {
		name    string
		args    args
		wantRes map[int64]*dm.SubjectInfo
		wantErr error
	}{
		{
			"normal",
			args{
				context.TODO(),
				1,
				1,
				[]int64{12412},
			},
			map[int64]*dm.SubjectInfo{},
			nil,
		},
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotRes, err := d.SubjectInfos(tt.args.c, tt.args.typ, tt.args.plat, tt.args.oids...)
			So(err, ShouldEqual, tt.wantErr)
			So(gotRes, ShouldResemble, tt.wantRes)
		})
	}
}
