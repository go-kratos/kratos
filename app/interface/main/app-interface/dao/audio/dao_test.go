package audio

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/audio"
	"go-common/library/net/ip"
	"path/filepath"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../../cmd/app-interface-test.toml")
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

func TestDao_Audios(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		pn  int
		ps  int
	}
	tests := []struct {
		name       string
		args       args
		wantAudios []*audio.Audio
		wantTotal  int
		wantErr    bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAudios, gotTotal, err := d.Audios(tt.args.c, tt.args.mid, tt.args.pn, tt.args.ps)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.Audios() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAudios, tt.wantAudios) {
				t.Errorf("Dao.Audios() gotAudios = %v, want %v", gotAudios, tt.wantAudios)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("Dao.Audios() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func TestDao_AllAudio(t *testing.T) {
	type args struct {
		c    context.Context
		vmid int64
	}
	tests := []struct {
		name    string
		args    args
		wantAus []*audio.Audio
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAus, err := d.AllAudio(tt.args.c, tt.args.vmid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.AllAudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAus, tt.wantAus) {
				t.Errorf("Dao.AllAudio() = %v, want %v", gotAus, tt.wantAus)
			}
		})
	}
}

func TestDao_AudioDetail(t *testing.T) {
	type args struct {
		c   context.Context
		ids []int64
	}
	tests := []struct {
		name    string
		args    args
		wantAum map[int64]*audio.Audio
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAum, err := d.AudioDetail(tt.args.c, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.AudioDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAum, tt.wantAum) {
				t.Errorf("Dao.AudioDetail() = %v, want %v", gotAum, tt.wantAum)
			}
		})
	}
}

func TestDao_FavAudio(t *testing.T) {
	type args struct {
		c         context.Context
		accessKey string
		mid       int64
		pn        int
		ps        int
	}
	tests := []struct {
		name    string
		args    args
		wantAus []*audio.FavAudio
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAus, err := d.FavAudio(tt.args.c, tt.args.accessKey, tt.args.mid, tt.args.pn, tt.args.ps)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.FavAudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAus, tt.wantAus) {
				t.Errorf("Dao.FavAudio() = %v, want %v", gotAus, tt.wantAus)
			}
		})
	}
}

func TestDao_UpperCert(t *testing.T) {
	type args struct {
		c   context.Context
		uid int64
	}
	tests := []struct {
		name     string
		args     args
		wantCert *audio.UpperCert
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCert, err := d.UpperCert(tt.args.c, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.UpperCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCert, tt.wantCert) {
				t.Errorf("Dao.UpperCert() = %v, want %v", gotCert, tt.wantCert)
			}
		})
	}
}

func TestDao_Card(t *testing.T) {
	type args struct {
		c   context.Context
		ip  string
		mid []int64
	}
	tests := []struct {
		name      string
		args      args
		wantCardm map[int64]*audio.Card
		wantErr   error
	}{
		{
			"normal",
			args{
				context.TODO(),
				ip.InternalIP(),
				[]int64{30047},
			},
			map[int64]*audio.Card{30047: &audio.Card{Type: 1, Status: 1}},
			nil,
		},
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotCardm, err := d.Card(tt.args.c, tt.args.mid...)
			So(err, ShouldBeNil)
			So(gotCardm, ShouldResemble, tt.wantCardm)
		})
	}
}

func TestDao_Fav(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
	}
	tests := []struct {
		name    string
		args    args
		wantFav *audio.Fav
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFav, err := d.Fav(tt.args.c, tt.args.mid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.Fav() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFav, tt.wantFav) {
				t.Errorf("Dao.Fav() = %v, want %v", gotFav, tt.wantFav)
			}
		})
	}
}
