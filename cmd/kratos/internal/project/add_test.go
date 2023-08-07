package project

import "testing"

func Test_subtractPath(t *testing.T) {
	type args struct {
		basePath     string
		subtractPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-1",
			args: args{
				basePath:     "/home/kratos/code/test/a/c",
				subtractPath: "/home/kratos/code/test",
			},
			want: "a/c",
		},
		{
			name: "test-2",
			args: args{
				basePath:     "code/test/go-kratos/kratos/app/user-service",
				subtractPath: "code/test/go-kratos/kratos",
			},
			want: "app/user-service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := subtractPath(tt.args.basePath, tt.args.subtractPath); got != tt.want {
				t.Errorf("subtractPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
