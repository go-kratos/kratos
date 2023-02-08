package apollo

import (
	"testing"
)

func Test_genKey(t *testing.T) {
	type args struct {
		ns  string
		sub string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "blank namespace",
			args: args{
				ns:  "",
				sub: "x.y",
			},
			want: "x.y",
		},
		{
			name: "properties namespace",
			args: args{
				ns:  "application",
				sub: "x.y",
			},
			want: "application.x.y",
		},
		{
			name: "namespace with format",
			args: args{
				ns:  "app.yaml",
				sub: "x.y",
			},
			want: "app.x.y",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genKey(tt.args.ns, tt.args.sub); got != tt.want {
				t.Errorf("genKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_format(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		want      string
	}{
		{
			name:      "properties namespace",
			namespace: "application",
			want:      "json",
		},
		{
			name:      "properties namespace #1",
			namespace: "app.setting",
			want:      "json",
		},
		{
			name:      "namespace with format[yaml]",
			namespace: "app.yaml",
			want:      "yaml",
		},
		{
			name:      "namespace with format[yml]",
			namespace: "app.yml",
			want:      "yml",
		},
		{
			name:      "namespace with format[json]",
			namespace: "app.json",
			want:      "json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.namespace); got != tt.want {
				t.Errorf("format() = %v, want %v", got, tt.want)
			}
		})
	}
}
