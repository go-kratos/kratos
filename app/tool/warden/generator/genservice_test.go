package generator

import "testing"

func TestUnderscore(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "remoteIP",
			args: args{s: "remoteIP"},
			want: "remote_ip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := underscore(tt.args.s); got != tt.want {
				t.Errorf("underscore() = %v, want %v", got, tt.want)
			}
		})
	}
}
