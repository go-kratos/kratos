package gorm

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	dsn = "root:secret@tcp(127.0.0.1:3306)/mysql?timeout=20s&readTimeout=20s"
)

func TestOpen(t *testing.T) {
	type args struct {
		dialect string
		options *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				dialect: "mysql",
				options: &Config{
					DSN:         dsn,
					Debug:       false,
					Idle:        10,
					Active:      100,
					IdleTimeout: time.Second * 300,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.args.dialect, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
			rows, err := got.Table("user").Select("User").Rows()
			fmt.Printf("rows = %+v\n", rows)
			assert.Nil(t, err)
			assert.NotNil(t, rows)
		})
	}
}
