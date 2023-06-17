//go:build linux
// +build linux

package project

import (
	"testing"
)

func Test_processProjectParams(t *testing.T) {
	type args struct {
		projectName      string
		fallbackPlaceDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"absLinux", args{projectName: "/home/kratos/awesome/go/demo", fallbackPlaceDir: ""}, "/home/kratos/awesome/go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := processProjectParams(tt.args.projectName, tt.args.fallbackPlaceDir); got != tt.want {
				t.Errorf("processProjectParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
