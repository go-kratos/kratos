package favorite

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-interface/model/favorite"

	"github.com/smartystreets/goconvey/convey"
)

var s *Service

func TestFolder(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Folder", t, func(ctx convey.C) {
		folder, err := s.Folder(c, "", "", "", "", "", 0, 1, 27515397, 27515397)
		b, _ := json.Marshal(folder)
		fmt.Printf("%s", b)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestService_FolderVideo(t *testing.T) {
	type args struct {
		c         context.Context
		accessKey string
		actionKey string
		device    string
		mobiApp   string
		platform  string
		keyword   string
		order     string
		build     int
		tid       int
		pn        int
		ps        int
		mid       int64
		fid       int64
		vmid      int64
	}
	tests := []struct {
		name       string
		s          *Service
		args       args
		wantFolder *favorite.FavideoList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFolder := tt.s.FolderVideo(tt.args.c, tt.args.accessKey, tt.args.actionKey, tt.args.device, tt.args.mobiApp, tt.args.platform, tt.args.keyword, tt.args.order, tt.args.build, tt.args.tid, tt.args.pn, tt.args.ps, tt.args.mid, tt.args.fid, tt.args.vmid); !reflect.DeepEqual(gotFolder, tt.wantFolder) {
				t.Errorf("Service.FolderVideo() = %v, want %v", gotFolder, tt.wantFolder)
			}
		})
	}
}

func TestService_Topic(t *testing.T) {
	type args struct {
		c         context.Context
		accessKey string
		actionKey string
		device    string
		mobiApp   string
		platform  string
		build     int
		ps        int
		pn        int
		mid       int64
	}
	tests := []struct {
		name      string
		s         *Service
		args      args
		wantTopic *favorite.TopicList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTopic := tt.s.Topic(tt.args.c, tt.args.accessKey, tt.args.actionKey, tt.args.device, tt.args.mobiApp, tt.args.platform, tt.args.build, tt.args.ps, tt.args.pn, tt.args.mid); !reflect.DeepEqual(gotTopic, tt.wantTopic) {
				t.Errorf("Service.Topic() = %v, want %v", gotTopic, tt.wantTopic)
			}
		})
	}
}

func TestService_Article(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name        string
		s           *Service
		args        args
		wantArticle *favorite.ArticleList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotArticle := tt.s.Article(tt.args.c, tt.args.mid, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotArticle, tt.wantArticle) {
				t.Errorf("Service.Article() = %v, want %v", gotArticle, tt.wantArticle)
			}
		})
	}
}

func TestService_Clips(t *testing.T) {
	type args struct {
		c         context.Context
		mid       int64
		accessKey string
		actionKey string
		device    string
		mobiApp   string
		platform  string
		build     int
		pn        int
		ps        int
	}
	tests := []struct {
		name      string
		s         *Service
		args      args
		wantClips *favorite.ClipsList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotClips := tt.s.Clips(tt.args.c, tt.args.mid, tt.args.accessKey, tt.args.actionKey, tt.args.device, tt.args.mobiApp, tt.args.platform, tt.args.build, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotClips, tt.wantClips) {
				t.Errorf("Service.Clips() = %v, want %v", gotClips, tt.wantClips)
			}
		})
	}
}

func TestService_Albums(t *testing.T) {
	type args struct {
		c         context.Context
		mid       int64
		accessKey string
		actionKey string
		device    string
		mobiApp   string
		platform  string
		build     int
		pn        int
		ps        int
	}
	tests := []struct {
		name       string
		s          *Service
		args       args
		wantAlbums *favorite.AlbumsList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAlbums := tt.s.Albums(tt.args.c, tt.args.mid, tt.args.accessKey, tt.args.actionKey, tt.args.device, tt.args.mobiApp, tt.args.platform, tt.args.build, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotAlbums, tt.wantAlbums) {
				t.Errorf("Service.Albums() = %v, want %v", gotAlbums, tt.wantAlbums)
			}
		})
	}
}

func TestService_Specil(t *testing.T) {
	type args struct {
		c         context.Context
		accessKey string
		actionKey string
		device    string
		mobiApp   string
		platform  string
		build     int
		pn        int
		ps        int
	}
	tests := []struct {
		name       string
		s          *Service
		args       args
		wantSpecil *favorite.SpList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSpecil := tt.s.Specil(tt.args.c, tt.args.accessKey, tt.args.actionKey, tt.args.device, tt.args.mobiApp, tt.args.platform, tt.args.build, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotSpecil, tt.wantSpecil) {
				t.Errorf("Service.Specil() = %v, want %v", gotSpecil, tt.wantSpecil)
			}
		})
	}
}

func TestService_Audio(t *testing.T) {
	type args struct {
		c         context.Context
		accessKey string
		mid       int64
		pn        int
		ps        int
	}
	tests := []struct {
		name      string
		s         *Service
		args      args
		wantAudio *favorite.AudioList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAudio := tt.s.Audio(tt.args.c, tt.args.accessKey, tt.args.mid, tt.args.pn, tt.args.ps); !reflect.DeepEqual(gotAudio, tt.wantAudio) {
				t.Errorf("Service.Audio() = %v, want %v", gotAudio, tt.wantAudio)
			}
		})
	}
}

func TestTab(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Tab", t, func(ctx convey.C) {
		tab, err := s.Tab(c, "", "", "", "", "", "", 0, 27515397)
		b, _ := json.Marshal(tab)
		fmt.Printf("%s", b)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
