package offer

import (
	"reflect"
	"testing"

	"go-common/app/job/main/app-wall/model/offer"
)

func TestService_activeConsumer(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.activeConsumer()
		})
	}
}

func TestService_checkMsgIllegal(t *testing.T) {
	type args struct {
		msg []byte
	}
	tests := []struct {
		name       string
		s          *Service
		args       args
		wantActive *offer.ActiveMsg
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActive, err := tt.s.checkMsgIllegal(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.checkMsgIllegal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotActive, tt.wantActive) {
				t.Errorf("Service.checkMsgIllegal() = %v, want %v", gotActive, tt.wantActive)
			}
		})
	}
}

func TestService_activeproc(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.activeproc()
		})
	}
}

func TestService_active(t *testing.T) {
	type args struct {
		msg *offer.ActiveMsg
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
			tt.s.active(tt.args.msg)
		})
	}
}
