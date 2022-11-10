package log

import "testing"

func TestLevel_Key(t *testing.T) {
	if LevelInfo.Key() != LevelKey {
		t.Errorf("want: %s, got: %s", LevelKey, LevelInfo.Key())
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name string
		l    Level
		want string
	}{
		{
			name: "DEBUG",
			l:    LevelDebug,
			want: "DEBUG",
		},
		{
			name: "INFO",
			l:    LevelInfo,
			want: "INFO",
		},
		{
			name: "WARN",
			l:    LevelWarn,
			want: "WARN",
		},
		{
			name: "ERROR",
			l:    LevelError,
			want: "ERROR",
		},
		{
			name: "FATAL",
			l:    LevelFatal,
			want: "FATAL",
		},
		{
			name: "other",
			l:    10,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want Level
	}{
		{
			name: "DEBUG",
			want: LevelDebug,
			s:    "DEBUG",
		},
		{
			name: "INFO",
			want: LevelInfo,
			s:    "INFO",
		},
		{
			name: "WARN",
			want: LevelWarn,
			s:    "WARN",
		},
		{
			name: "ERROR",
			want: LevelError,
			s:    "ERROR",
		},
		{
			name: "FATAL",
			want: LevelFatal,
			s:    "FATAL",
		},
		{
			name: "other",
			want: LevelInfo,
			s:    "other",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLevel(tt.s); got != tt.want {
				t.Errorf("ParseLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
