package space

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go-common/app/interface/main/app-interface/model/space"
	account "go-common/app/service/main/account/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Space(t *testing.T) {
	type args struct {
		c        context.Context
		mid      int64
		vmid     int64
		plat     int8
		build    int
		pn       int
		ps       int
		platform string
		device   string
		mobiApp  string
		name     string
		now      time.Time
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantSp  *space.Space
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func(t *testing.T) {
			gotSp, err := tt.s.Space(tt.args.c, tt.args.mid, tt.args.vmid, tt.args.plat, tt.args.build, tt.args.pn, tt.args.ps, tt.args.platform, tt.args.device, tt.args.mobiApp, tt.args.name, tt.args.now)
			So(err, ShouldEqual, tt.wantErr)
			So(gotSp, ShouldNotResemble, tt.wantSp)
		})
	}
}

func TestService_UpArcs(t *testing.T) {
	type args struct {
		c   context.Context
		uid int64
		pn  int
		ps  int
		now time.Time
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.ArcList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.UpArcs(tt.args.c, tt.args.uid, tt.args.pn, tt.args.ps, tt.args.now); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.UpArcs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_UpArticles(t *testing.T) {
	type args struct {
		c   context.Context
		uid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.ArticleList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.UpArticles(tt.args.c, tt.args.uid, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.UpArticles() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_favFolders(t *testing.T) {
	type args struct {
		c       context.Context
		mid     int64
		uid     int64
		plat    int8
		build   int
		mobiApp string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.FavList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.favFolders(tt.args.c, tt.args.mid, tt.args.uid, nil, tt.args.plat, tt.args.build, tt.args.mobiApp); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.favFolders() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_Bangumi(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		uid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.BangumiList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.Bangumi(tt.args.c, tt.args.mid, tt.args.uid, nil, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.Bangumi() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_Community(t *testing.T) {
	type args struct {
		c        context.Context
		uid      int64
		pn       int
		ps       int
		ak       string
		platform string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.CommuList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.Community(tt.args.c, tt.args.uid, tt.args.pn, tt.args.ps, tt.args.ak, tt.args.platform); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.Community() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_CoinArcs(t *testing.T) {
	type args struct {
		c   context.Context
		uid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.ArcList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.CoinArcs(tt.args.c, 0, tt.args.uid, nil, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.CoinArcs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_LikeArcs(t *testing.T) {
	type args struct {
		c   context.Context
		uid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.ArcList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.LikeArcs(tt.args.c, 0, tt.args.uid, nil, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.LikeArcs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_upClips(t *testing.T) {
	type args struct {
		c   context.Context
		uid int64
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.ClipList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.upClips(tt.args.c, tt.args.uid); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.upClips() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_upAlbums(t *testing.T) {
	type args struct {
		c   context.Context
		uid int64
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.AlbumList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.upAlbums(tt.args.c, tt.args.uid); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.upAlbums() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_card(t *testing.T) {
	type args struct {
		c    context.Context
		vmid int64
		name string
	}
	tests := []struct {
		name     string
		s        *Service
		args     args
		wantCard *account.Card
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCard, err := tt.s.card(tt.args.c, tt.args.vmid, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.card() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCard, tt.wantCard) {
				t.Errorf("Service.card() = %v, want %v", gotCard, tt.wantCard)
			}
		})
	}
}

func TestService_audios(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.AudioList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.s.audios(tt.args.c, tt.args.mid, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.audios() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_Report(t *testing.T) {
	type args struct {
		c      context.Context
		mid    int64
		reason string
		ak     string
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
			if err := tt.s.Report(tt.args.c, tt.args.mid, tt.args.reason, tt.args.ak); (err != nil) != tt.wantErr {
				t.Errorf("Service.Report() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
