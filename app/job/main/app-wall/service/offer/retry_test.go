package offer

import (
	"context"
	"testing"
	"time"
)

func TestService_retryproc(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.retryproc()
		})
	}
}

func Test_retry(t *testing.T) {
	type args struct {
		callback func() error
		retry    int
		sleep    time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := retry(tt.args.callback, tt.args.retry, tt.args.sleep); (err != nil) != tt.wantErr {
				t.Errorf("retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_syncRetry(t *testing.T) {
	type args struct {
		c         context.Context
		action    string
		os        string
		imei      string
		androidid string
		mac       string
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
			if err := tt.s.syncRetry(tt.args.c, tt.args.action, tt.args.os, tt.args.imei, tt.args.androidid, tt.args.mac); (err != nil) != tt.wantErr {
				t.Errorf("Service.syncRetry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
