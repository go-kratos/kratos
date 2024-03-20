package lumberjack

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	lumberjack2 "gopkg.in/natefinch/lumberjack.v2"
)

func MustParse(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

func TestNewLoggerWithURL(t *testing.T) {
	type args struct {
		config *Config
		u      *url.URL
	}
	tests := []struct {
		name string
		args args
		want *Logger
	}{
		{
			name: "default",
			want: newDefaultConfig(),
		},
		{
			name: "custom",
			args: args{
				config: &Config{
					Maxsize:    1023,
					Maxbackups: 10,
					Maxage:     10,
					Compress:   true,
					Localtime:  true,
					Filename:   "test.log",
				},
			},
			want: &Logger{
				Logger: &lumberjack2.Logger{
					MaxSize:    1023,
					MaxBackups: 10,
					MaxAge:     10,
					Compress:   true,
					LocalTime:  true,
					Filename:   "test.log",
				},
			},
		},
		{
			name: "custom-url-ignore-1",
			args: args{
				config: &Config{
					Filename: "test.log",
				},
				u: MustParse("lumberjack"),
			},
			want: &Logger{
				Logger: &lumberjack2.Logger{
					MaxSize:    1024,
					MaxBackups: 3,
					MaxAge:     7,
					Filename:   "test.log",
				},
			},
		},
		{
			name: "custom-url-ignore-2",
			args: args{
				config: &Config{
					Filename: "test.log",
				},
				u: MustParse("lumberjack:"),
			},
			want: &Logger{
				Logger: &lumberjack2.Logger{
					MaxSize:    1024,
					MaxBackups: 3,
					MaxAge:     7,
					Filename:   "test.log",
				},
			},
		},
		{
			name: "custom-url-ignore-3",
			args: args{
				config: &Config{
					Filename: "test.log",
				},
				u: MustParse("lumberjack:/logs/info.log"),
			},
			want: &Logger{
				Logger: &lumberjack2.Logger{
					MaxSize:    1024,
					MaxBackups: 3,
					MaxAge:     7,
					Filename:   "test.log",
				},
			},
		},
		{
			name: "custom-url-0",
			args: args{
				u: MustParse("lumberjack:logs/info.log"),
			},
			want: &Logger{
				Logger: &lumberjack2.Logger{
					MaxSize:    1024,
					MaxBackups: 3,
					MaxAge:     7,
					Filename:   "logs/info.log",
				},
			},
		},
		{
			name: "custom-url-1",
			args: args{
				u: MustParse("lumberjack:logs/info.log"),
			},
			want: &Logger{
				Logger: &lumberjack2.Logger{
					MaxSize:    1024,
					MaxBackups: 3,
					MaxAge:     7,
					Filename:   "logs/info.log",
				},
			},
		},
		{
			name: "custom-url-2",
			args: args{
				u: MustParse("lumberjack:/logs/info.log"),
			},
			want: &Logger{
				Logger: &lumberjack2.Logger{
					MaxSize:    1024,
					MaxBackups: 3,
					MaxAge:     7,
					Filename:   "/logs/info.log",
				},
			},
		},
		{
			name: "custom-url-3",
			args: args{
				u: MustParse("lumberjack:./logs/info.log"),
			},
			want: &Logger{
				Logger: &lumberjack2.Logger{
					MaxSize:    1024,
					MaxBackups: 3,
					MaxAge:     7,
					Filename:   "./logs/info.log",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLoggerWithURL(tt.args.config, tt.args.u)
			assert.Equal(t, got.Logger.Compress, tt.want.Logger.Compress)
			assert.Equal(t, got.Logger.MaxAge, tt.want.Logger.MaxAge)
			assert.Equal(t, got.Logger.MaxBackups, tt.want.Logger.MaxBackups)
			assert.Equal(t, got.Logger.MaxSize, tt.want.Logger.MaxSize)
			assert.Equal(t, got.Logger.LocalTime, tt.want.Logger.LocalTime)
			assert.Equal(t, got.Logger.Filename, tt.want.Logger.Filename)
		})
	}
}
