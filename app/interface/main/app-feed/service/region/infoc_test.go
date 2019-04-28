package region

import (
	"testing"
	"time"

	"go-common/app/interface/main/app-feed/model/tag"
)

func TestService_TagsInfoc(t *testing.T) {
	type args struct {
		mid   int64
		plat  int8
		build int
		buvid string
		disid string
		ip    string
		api   string
		tags  []*tag.Tag
		now   time.Time
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
			tt.s.TagsInfoc(tt.args.mid, tt.args.plat, tt.args.build, tt.args.buvid, tt.args.disid, tt.args.ip, tt.args.api, tt.args.tags, tt.args.now)
		})
	}
}

func TestService_ChangeTagsInfoc(t *testing.T) {
	type args struct {
		mid   int64
		plat  int8
		build int
		buvid string
		disid string
		ip    string
		api   string
		tags  []*tag.Tag
		now   time.Time
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
			tt.s.ChangeTagsInfoc(tt.args.mid, tt.args.plat, tt.args.build, tt.args.buvid, tt.args.disid, tt.args.ip, tt.args.api, tt.args.tags, tt.args.now)
		})
	}
}

func TestService_AddTagInfoc(t *testing.T) {
	type args struct {
		mid   int64
		plat  int8
		build int
		buvid string
		disid string
		ip    string
		api   string
		rid   int
		tid   int64
		now   time.Time
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
			tt.s.AddTagInfoc(tt.args.mid, tt.args.plat, tt.args.build, tt.args.buvid, tt.args.disid, tt.args.ip, tt.args.api, tt.args.rid, tt.args.tid, tt.args.now)
		})
	}
}

func TestService_CancelTagInfoc(t *testing.T) {
	type args struct {
		mid   int64
		plat  int8
		build int
		buvid string
		disid string
		ip    string
		api   string
		rid   int
		tid   int64
		now   time.Time
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
			tt.s.CancelTagInfoc(tt.args.mid, tt.args.plat, tt.args.build, tt.args.buvid, tt.args.disid, tt.args.ip, tt.args.api, tt.args.rid, tt.args.tid, tt.args.now)
		})
	}
}

func TestService_infocproc(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.infocproc()
		})
	}
}
