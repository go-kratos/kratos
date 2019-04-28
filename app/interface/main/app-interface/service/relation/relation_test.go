package relation

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/interface/main/app-interface/conf"
	model "go-common/app/interface/main/app-interface/model/relation"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-interface-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(3 * time.Second)
}

func TestService_Followings(t *testing.T) {
	type args struct {
		c       context.Context
		vmid    int64
		mid     int64
		pn      int
		ps      int
		version uint64
		order   string
	}
	tests := []struct {
		name       string
		args       args
		wantF      []*model.Following
		wantCrc32v uint32
		wantTotal  int
		wantErr    error
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotF, gotCrc32v, gotTotal, err := s.Followings(tt.args.c, tt.args.vmid, tt.args.mid, tt.args.pn, tt.args.ps, tt.args.version, tt.args.order)
			So(gotF, ShouldResemble, tt.wantF)
			So(gotCrc32v, ShouldResemble, tt.wantCrc32v)
			So(gotTotal, ShouldResemble, tt.wantTotal)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestService_Tag(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		tid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantT   []*model.Tag
		wantErr error
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := tt.s.Tag(tt.args.c, tt.args.mid, tt.args.tid, tt.args.pn, tt.args.ps)
			So(gotT, ShouldResemble, tt.wantT)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}
