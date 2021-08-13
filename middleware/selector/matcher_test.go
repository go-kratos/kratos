package selector

import (
	"regexp"
	"testing"
)

func Test_regexMatcher_Match(t *testing.T) {
	type fields struct {
		re *regexp.Regexp
	}
	type args struct {
		operation string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{fields: fields{re: regexp.MustCompile(`/hello/[0-9]+`)}, args: args{operation: "/hello/2021"}, want: true},
		{fields: fields{re: regexp.MustCompile(`^/hello/[0-9]+$`)}, args: args{operation: "/hello/2021"}, want: true},
		{fields: fields{re: regexp.MustCompile(`/hello/[0-9]+`)}, args: args{operation: "/hello/2021/"}, want: false},
		{fields: fields{re: regexp.MustCompile(`^/hello/[0-9]+$`)}, args: args{operation: "/hello/2021/"}, want: false},
		{fields: fields{re: regexp.MustCompile(`/hello/[0-9]+`)}, args: args{operation: "/hello/2021/kratos"}, want: false},
		{fields: fields{re: regexp.MustCompile(`^/hello/[0-9]+$`)}, args: args{operation: "/hello/2021/kratos"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := regexMatcher{
				re: tt.fields.re,
			}
			if got := m.Match(tt.args.operation); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
