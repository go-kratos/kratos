package memcache

import (
	"testing"

	pb "github.com/bilibili/kratos/pkg/cache/memcache/test"

	"github.com/stretchr/testify/assert"
)

func TestItemUtil(t *testing.T) {
	item1 := RawItem("test", []byte("hh"), 0, 0)
	assert.Equal(t, "test", item1.Key)
	assert.Equal(t, []byte("hh"), item1.Value)
	assert.Equal(t, FlagRAW, FlagRAW&item1.Flags)

	item1 = JSONItem("test", &Item{}, 0, 0)
	assert.Equal(t, "test", item1.Key)
	assert.NotNil(t, item1.Object)
	assert.Equal(t, FlagJSON, FlagJSON&item1.Flags)

	item1 = ProtobufItem("test", &pb.TestItem{}, 0, 0)
	assert.Equal(t, "test", item1.Key)
	assert.NotNil(t, item1.Object)
	assert.Equal(t, FlagProtobuf, FlagProtobuf&item1.Flags)
}

func TestLegalKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test empty key",
			want: false,
		},
		{
			name: "test too large key",
			args: args{func() string {
				var data []byte
				for i := 0; i < 255; i++ {
					data = append(data, 'k')
				}
				return string(data)
			}()},
			want: false,
		},
		{
			name: "test invalid char",
			args: args{"hello world"},
			want: false,
		},
		{
			name: "test invalid char",
			args: args{string([]byte{0x7f})},
			want: false,
		},
		{
			name: "test normal key",
			args: args{"hello"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := legalKey(tt.args.key); got != tt.want {
				t.Errorf("legalKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
