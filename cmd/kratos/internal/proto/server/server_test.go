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
			name: "serviceName on lowercase words",
			args: args{str: "helloworld"},
			want: "Helloworld",
		},
		{
			name: "serviceName on uppercase words",
			args: args{str: "HELLOWORLD"},
			want: "HELLOWORLD",
		},
		{
			name: "serviceName on lowercase words with spaces",
			args: args{str: "hello world"},
			want: "HelloWorld",
		},
		{
			name: "serviceName on uppercase words with spaces",
			args: args{str: "HELLO WORLD"},
			want: "HELLOWORLD",
		},
		{
			name: "serviceName on Lower Camel Case words",
			args: args{str: "helloWorld"},
			want: "HelloWorld",
		},
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

func Test_parametersName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "parametersName on not nested",
			args: args{
				name: "MessageResponse",
			},
			want: "MessageResponse",
		},
		{
			name: "parametersName on One layer of nesting",
			args: args{
				name: "Message.Response",
			},
			want: "Message_Response",
		},
		{
			name: "parametersName on Two layer of nesting",
			args: args{
				name: "Message.Message2.Response",
			},
			want: "Message_Message2_Response",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parametersName(tt.args.name); got != tt.want {
				t.Errorf("parametersName() = %v, want %v", got, tt.want)
			}
		})
	}
}
