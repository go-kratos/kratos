package selector

import (
	"testing"
)

func Test_orMatcher_Match(t *testing.T) {
	type fields struct {
		subMatchers []OperationMatcher
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
		{fields: fields{subMatchers: nil}, want: false},
		{fields: fields{subMatchers: []OperationMatcher{}}, want: false},
		{fields: fields{subMatchers: []OperationMatcher{boolMatcher(false), boolMatcher(false), boolMatcher(false)}}, want: false},
		{fields: fields{subMatchers: []OperationMatcher{boolMatcher(false), boolMatcher(true), boolMatcher(false)}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := orMatcher{
				subMatchers: tt.fields.subMatchers,
			}
			if got := m.Match(tt.args.operation); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_andMatcher_Match(t *testing.T) {
	type fields struct {
		subMatchers []OperationMatcher
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
		{fields: fields{subMatchers: nil}, want: true},
		{fields: fields{subMatchers: []OperationMatcher{}}, want: true},
		{fields: fields{subMatchers: []OperationMatcher{boolMatcher(true), boolMatcher(true), boolMatcher(true)}}, want: true},
		{fields: fields{subMatchers: []OperationMatcher{boolMatcher(true), boolMatcher(false), boolMatcher(true)}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := andMatcher{
				subMatchers: tt.fields.subMatchers,
			}
			if got := m.Match(tt.args.operation); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

type boolMatcher bool

var _ OperationMatcher = boolMatcher(false)

func (b boolMatcher) Match(operation string) bool {
	return bool(b)
}
