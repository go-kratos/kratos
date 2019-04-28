package region

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"go-common/app/interface/main/app-feed/conf"
	"go-common/app/interface/main/app-feed/model/tag"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-feed-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func TestService_HotTags(t *testing.T) {
	type args struct {
		c    context.Context
		mid  int64
		rid  int16
		ver  string
		plat int8
		now  time.Time
	}
	tests := []struct {
		name        string
		args        args
		wantHs      []*tag.Hot
		wantVersion string
		wantErr     error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHs, gotVersion, err := s.HotTags(tt.args.c, tt.args.mid, tt.args.rid, tt.args.ver, tt.args.plat, tt.args.now)
			So(gotHs, ShouldEqual, tt.wantHs)
			So(gotVersion, ShouldEqual, tt.wantVersion)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestService_SubTags(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name  string
		s     *Service
		args  args
		wantT *tag.SubTag
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotT := tt.s.SubTags(tt.args.c, tt.args.mid, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("Service.SubTags() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func TestService_AddTag(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		tid int64
		now time.Time
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddTag(tt.args.c, tt.args.mid, tt.args.tid, tt.args.now); (err != nil) != tt.wantErr {
				t.Errorf("Service.AddTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_CancelTag(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		tid int64
		now time.Time
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.CancelTag(tt.args.c, tt.args.mid, tt.args.tid, tt.args.now); (err != nil) != tt.wantErr {
				t.Errorf("Service.CancelTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
