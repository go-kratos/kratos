package service

import (
	"testing"

	"go-common/app/job/main/app/model/space"
	xtime "go-common/library/time"
)

func TestService_contributeConsumeproc(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.contributeConsumeproc()
		})
	}
}

func TestService_contributeroc(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.contributeroc()
		})
	}
}

func TestService_contributeCache(t *testing.T) {
	type args struct {
		vmid  int64
		attrs *space.Attrs
		ctime xtime.Time
		ip    string
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
			if err := tt.s.contributeCache(tt.args.vmid, tt.args.attrs, tt.args.ctime, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("Service.contributeCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_contributeUpdate(t *testing.T) {
	type args struct {
		vmid  int64
		attrs *space.Attrs
		items []*space.Item
	}
	tests := []struct {
		name string
		s    *Service
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.contributeUpdate(tt.args.vmid, tt.args.attrs, tt.args.items)
		})
	}
}
