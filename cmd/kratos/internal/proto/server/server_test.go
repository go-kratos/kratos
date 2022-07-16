package server

import "testing"

func Test_ucFirst(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ucFirst on lowercase words",
			args: args{str: "helloworld"},
			want: "Helloworld",
		},
		{
			name: "ucFirst on uppercase words",
			args: args{str: "HELLOWORLD"},
			want: "HELLOWORLD",
		},
		{
			name: "ucFirst on lowercase words with spaces",
			args: args{str: "hello world"},
			want: "Hello world",
		},
		{
			name: "ucFirst on uppercase words with spaces",
			args: args{str: "HELLO WORLD"},
			want: "HELLO WORLD",
		},
		{
			name: "ucFirst on Lower Camel Case words",
			args: args{str: "helloWorld"},
			want: "HelloWorld",
		},
		{
			name: "ucFirst on Upper Camel Case words",
			args: args{str: "HelloWorld"},
			want: "HelloWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ucFirst(tt.args.str); got != tt.want {
				t.Errorf("ucFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}
