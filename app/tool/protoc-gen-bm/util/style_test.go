package util

import "testing"

func TestCamelCase(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{name: "_test"},
			want: "_test",
		},
		{
			name: "test2",
			args: args{name: "hello_world"},
			want: "HelloWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelCase(tt.args.name); got != tt.want {
				t.Errorf("CamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
