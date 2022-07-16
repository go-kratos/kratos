package server

import "testing"

func Test_serviceName(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "serviceName on Lower Camel Case words",
			args: args{str: "helloWorld"},
			want: "HelloWorld",
		},
		{
			name: "serviceName on Upper Camel Case words",
			args: args{str: "HelloWorld"},
			want: "HelloWorld",
		},
		{
			name: "serviceName on Upper Camel Case words",
			args: args{str: "hello_world"},
			want: "HelloWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serviceName(tt.args.str); got != tt.want {
				t.Errorf("serviceName() = %v, want %v", got, tt.want)
			}
		})
	}
}
