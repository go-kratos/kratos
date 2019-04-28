package feed

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-common/app/interface/main/app-card/model/card/ai"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_IndexInfoc(t *testing.T) {
	type args struct {
		c           context.Context
		mid         int64
		plat        int8
		build       int
		buvid       string
		disid       string
		api         string
		userFeature json.RawMessage
		style       int
		code        int
		items       []*ai.Item
		isRcmd      bool
		pull        bool
		newUser     bool
		now         time.Time
		zoneID      int64
		adResponse  string
		deviceID    string
		network     string
		flush       int
		autoPlay    string
		deviceType  int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			s.IndexInfoc(tt.args.c, tt.args.mid, tt.args.plat, tt.args.build, tt.args.buvid, tt.args.disid, tt.args.api, tt.args.userFeature, tt.args.style, tt.args.code, tt.args.items, tt.args.isRcmd, tt.args.pull, tt.args.newUser, tt.args.now, tt.args.adResponse, tt.args.deviceID, tt.args.network, tt.args.flush, tt.args.autoPlay, tt.args.deviceType)
		})
	}
}

func TestService_infoc(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.infoc(tt.args.i)
		})
	}
}

func TestService_infocproc(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.infocproc()
		})
	}
}

func Test_gotoMapID(t *testing.T) {
	type args struct {
		gt string
	}
	tests := []struct {
		name   string
		args   args
		wantId string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotId := gotoMapID(tt.args.gt); gotId != tt.wantId {
				t.Errorf("gotoMapID() = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}
