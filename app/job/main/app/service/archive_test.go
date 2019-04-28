package service

import "testing"

func TestService_arcConsumeproc(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.arcConsumeproc()
		})
	}
}

func TestService_archiveUpdate(t *testing.T) {
	type args struct {
		action string
		nwMsg  []byte
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
			tt.s.archiveUpdate(tt.args.action, tt.args.nwMsg)
		})
	}
}

func TestService_upViewCache(t *testing.T) {
	type args struct {
		aid int64
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
			tt.s.upViewCache([]int64{tt.args.aid})
		})
	}
}
