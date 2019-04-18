package memcache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "go-common/library/cache/memcache/test"
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
