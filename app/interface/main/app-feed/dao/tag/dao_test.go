package tag

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-feed/conf"
	tag "go-common/app/interface/main/tag/model"

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

func TestDao_InfoByIDs(t *testing.T) {
	type args struct {
		c    context.Context
		mid  int64
		tids []int64
	}
	tests := []struct {
		name    string
		args    args
		wantTm  map[int64]*tag.Tag
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotTm, err := d.InfoByIDs(tt.args.c, tt.args.mid, tt.args.tids)
			So(gotTm, ShouldEqual, tt.wantTm)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestDao_SubTags(t *testing.T) {
	type args struct {
		c    context.Context
		mid  int64
		vmid int64
		pn   int
		ps   int
	}
	tests := []struct {
		name    string
		args    args
		wantSub *tag.Sub
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSub, err := d.SubTags(tt.args.c, tt.args.mid, tt.args.vmid, tt.args.pn, tt.args.ps)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.SubTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSub, tt.wantSub) {
				t.Errorf("Dao.SubTags() = %v, want %v", gotSub, tt.wantSub)
			}
		})
	}
}
