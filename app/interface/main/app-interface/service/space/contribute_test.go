package space

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/space"
)

func TestService_Contribute(t *testing.T) {
	type args struct {
		c     context.Context
		plat  int8
		build int
		vmid  int64
		pn    int
		ps    int
		now   time.Time
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.Contributes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.Contribute(tt.args.c, tt.args.plat, tt.args.build, tt.args.vmid, tt.args.pn, tt.args.ps, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Contribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.Contribute() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_Contribution(t *testing.T) {
	type args struct {
		c      context.Context
		plat   int8
		build  int
		vmid   int64
		cursor *model.Cursor
		now    time.Time
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.Contributes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.Contribution(tt.args.c, tt.args.plat, tt.args.build, tt.args.vmid, tt.args.cursor, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Contribution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.Contribution() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_firstContribute(t *testing.T) {
	type args struct {
		c    context.Context
		vmid int64
		size int
		now  time.Time
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.Contributes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.firstContribute(tt.args.c, tt.args.vmid, tt.args.size, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.firstContribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.firstContribute() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_dealContribute(t *testing.T) {
	type args struct {
		c     context.Context
		plat  int8
		build int
		vmid  int64
		attrs *space.Attrs
		items []*space.Item
		now   time.Time
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes *space.Contributes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.dealContribute(tt.args.c, tt.args.plat, tt.args.build, tt.args.vmid, tt.args.attrs, tt.args.items, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.dealContribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.dealContribute() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_Clip(t *testing.T) {
	type args struct {
		c    context.Context
		vmid int64
		pos  int
		size int
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
			if gotRes := tt.s.Clip(tt.args.c, tt.args.vmid, tt.args.pos, tt.args.size); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.Clip() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_Album(t *testing.T) {
	type args struct {
		c    context.Context
		vmid int64
		pos  int
		size int
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
			if gotRes := tt.s.Album(tt.args.c, tt.args.vmid, tt.args.pos, tt.args.size); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.Album() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_AddContribute(t *testing.T) {
	type args struct {
		c     context.Context
		vmid  int64
		attrs *space.Attrs
		items []*space.Item
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
			if err := tt.s.AddContribute(tt.args.c, tt.args.vmid, tt.args.attrs, tt.args.items); (err != nil) != tt.wantErr {
				t.Errorf("Service.AddContribute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
