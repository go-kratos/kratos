package tests

import (
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"go-common/app/tool/gengo/cmd/deepcopy-gen/examples/api/v1"
	"go-common/app/tool/gengo/cmd/deepcopy-gen/examples/model"
	xtime "go-common/library/time"
)

func TestCopyFromAllBuitin(t *testing.T) {
	base := &model.MemberBase{
		Mid:      1,
		Name:     "name",
		Sex:      1,
		Face:     "aaa.jpg",
		Sign:     "sign",
		Rank:     10000,
		Birthday: xtime.Time(time.Now().Unix()),
	}
	reply := new(v1.BaseInfoReply)
	reply.DeepCopyFromMemberBase(base)

	assert.Equal(t, base.Mid, reply.Mid)
	assert.Equal(t, base.Name, reply.Name)
	assert.Equal(t, base.Sex, reply.Sex)
	assert.Equal(t, base.Face, reply.Face)
	assert.Equal(t, base.Sign, reply.Sign)
	assert.Equal(t, base.Rank, reply.Rank)
	assert.Equal(t, base.Birthday, reply.Birthday)

	copied := reply.DeepCopyAsMemberBase()
	assert.Equal(t, base, copied)
}

func TestCopyFromMap(t *testing.T) {
	names := &model.Names{
		Names: map[int64]string{
			1: "1",
			2: "2",
		},
	}
	reply := new(v1.NamesReply)
	reply.DeepCopyFromNames(names)

	assert.Equal(t, names.Names[1], reply.Names[1])
	assert.Equal(t, names.Names[2], reply.Names[2])
	assert.Equal(t, unsafe.Pointer(&names.Names), unsafe.Pointer(&names.Names))
	assert.Equal(t, unsafe.Pointer(&reply.Names), unsafe.Pointer(&reply.Names))
	assert.NotEqual(t, unsafe.Pointer(&names.Names), unsafe.Pointer(&reply.Names))

	names.Names[3] = "3"
	assert.Contains(t, names.Names, int64(3))
	assert.NotContains(t, reply.Names, int64(3))
}

func TestCopyFromSlice(t *testing.T) {
	mids := &model.Mids{
		Mids: []int64{1, 2},
	}
	reply := new(v1.MidsReply)
	reply.DeepCopyFromMids(mids)

	assert.Equal(t, mids.Mids[0], reply.Mids[0])
	assert.Equal(t, mids.Mids[1], reply.Mids[1])
	assert.Equal(t, unsafe.Pointer(&mids.Mids), unsafe.Pointer(&mids.Mids))
	assert.Equal(t, unsafe.Pointer(&reply.Mids), unsafe.Pointer(&reply.Mids))
	assert.NotEqual(t, unsafe.Pointer(&mids.Mids), unsafe.Pointer(&reply.Mids))

	mids.Mids = append(mids.Mids, 3)
	assert.Equal(t, mids.Mids[2], int64(3))
	assert.NotContains(t, reply.Mids, int64(3))
}
