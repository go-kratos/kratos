package feed

import (
	"context"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_convergeCard(t *testing.T) {
	type args struct {
		c     context.Context
		limit int
		ids   []int64
	}
	tests := []struct {
		name        string
		args        args
		wantCardm   map[int64]*operate.Converge
		wantAids    []int64
		wantRoomIDs []int64
		wantMetaIDs []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotCardm, gotAids, gotRoomIDs, gotMetaIDs := s.convergeCard(tt.args.c, tt.args.limit, tt.args.ids...)
			So(gotCardm, ShouldResemble, tt.wantCardm)
			So(gotAids, ShouldResemble, tt.wantAids)
			So(gotRoomIDs, ShouldResemble, tt.wantRoomIDs)
			So(gotMetaIDs, ShouldResemble, tt.wantMetaIDs)
		})
	}
}

func TestService_downloadCard(t *testing.T) {
	type args struct {
		c   context.Context
		ids []int64
	}
	tests := []struct {
		name      string
		args      args
		wantCardm map[int64]*operate.Download
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			if gotCardm := s.downloadCard(tt.args.c, tt.args.ids...); !reflect.DeepEqual(gotCardm, tt.wantCardm) {
				t.Errorf("Service.downloadCard() = %v, want %v", gotCardm, tt.wantCardm)
			}
		})
	}
}

func TestService_subscribeCard(t *testing.T) {
	type args struct {
		c   context.Context
		ids []int64
	}
	tests := []struct {
		name      string
		args      args
		wantCardm map[int64]*operate.Follow
		wantUpIDs []int64
		wantTids  []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotCardm, gotUpIDs, gotTids := s.subscribeCard(tt.args.c, tt.args.ids...)
			if !reflect.DeepEqual(gotCardm, tt.wantCardm) {
				t.Errorf("Service.subscribeCard() gotCardm = %v, want %v", gotCardm, tt.wantCardm)
			}
			if !reflect.DeepEqual(gotUpIDs, tt.wantUpIDs) {
				t.Errorf("Service.subscribeCard() gotUpIDs = %v, want %v", gotUpIDs, tt.wantUpIDs)
			}
			if !reflect.DeepEqual(gotTids, tt.wantTids) {
				t.Errorf("Service.subscribeCard() gotTids = %v, want %v", gotTids, tt.wantTids)
			}
		})
	}
}

func TestService_channelRcmdCard(t *testing.T) {
	type args struct {
		c   context.Context
		ids []int64
	}
	tests := []struct {
		name      string
		args      args
		wantCardm map[int64]*operate.Follow
		wantUpIDs []int64
		wantTids  []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotCardm, gotUpIDs, gotTids := s.channelRcmdCard(tt.args.c, tt.args.ids...)
			if !reflect.DeepEqual(gotCardm, tt.wantCardm) {
				t.Errorf("Service.channelRcmdCard() gotCardm = %v, want %v", gotCardm, tt.wantCardm)
			}
			if !reflect.DeepEqual(gotUpIDs, tt.wantUpIDs) {
				t.Errorf("Service.channelRcmdCard() gotUpIDs = %v, want %v", gotUpIDs, tt.wantUpIDs)
			}
			if !reflect.DeepEqual(gotTids, tt.wantTids) {
				t.Errorf("Service.channelRcmdCard() gotTids = %v, want %v", gotTids, tt.wantTids)
			}
		})
	}
}

func TestService_liveUpRcmdCard(t *testing.T) {
	type args struct {
		c   context.Context
		ids []int64
	}
	tests := []struct {
		name      string
		args      args
		wantCardm map[int64][]*live.Card
		wantUpIDs []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotCardm, gotUpIDs := s.liveUpRcmdCard(tt.args.c, tt.args.ids...)
			if !reflect.DeepEqual(gotCardm, tt.wantCardm) {
				t.Errorf("Service.liveUpRcmdCard() gotCardm = %v, want %v", gotCardm, tt.wantCardm)
			}
			if !reflect.DeepEqual(gotUpIDs, tt.wantUpIDs) {
				t.Errorf("Service.liveUpRcmdCard() gotUpIDs = %v, want %v", gotUpIDs, tt.wantUpIDs)
			}
		})
	}
}
