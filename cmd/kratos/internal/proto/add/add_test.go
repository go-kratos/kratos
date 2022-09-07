package add

import "testing"

func TestUnderscoreToUpperCamelCase(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "hello_world",
			want: "HelloWorld",
		},
		{
			name: "v2_kratos_dev",
			want: "V2KratosDev",
		},
		{
			name: "www_Google_com",
			want: "WwwGoogleCom",
		},
		{
			name: "wwwBaidu_com",
			want: "WwwBaiduCom",
		},
		{
			name: "HelloWorld",
			want: "HelloWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toUpperCamelCase(tt.name); got != tt.want {
				t.Errorf("toUpperCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
