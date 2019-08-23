package ecode

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"

	"github.com/bilibili/kratos/pkg/ecode/types"
)

func TestEqual(t *testing.T) {
	var (
		err1 = Error(RequestErr, "test")
		err2 = Errorf(RequestErr, "test")
	)
	assert.Equal(t, err1, err2)
	assert.True(t, Equal(nil, nil))
}

func TestDetail(t *testing.T) {
	m := &timestamp.Timestamp{Seconds: time.Now().Unix()}
	st, _ := Error(RequestErr, "RequestErr").WithDetails(m)

	assert.Equal(t, "RequestErr", st.Message())
	assert.Equal(t, int(RequestErr), st.Code())
	assert.IsType(t, m, st.Details()[0])
}

func TestFromCode(t *testing.T) {
	err := FromCode(RequestErr)

	assert.Equal(t, int(RequestErr), err.Code())
	assert.Equal(t, "-400", err.Message())
}

func TestFromProto(t *testing.T) {
	msg := &types.Status{Code: 2233, Message: "error"}
	err := FromProto(msg)

	assert.Equal(t, 2233, err.Code())
	assert.Equal(t, "error", err.Message())

	m := &timestamp.Timestamp{Seconds: time.Now().Unix()}
	err = FromProto(m)
	assert.Equal(t, -500, err.Code())
	assert.Contains(t, err.Message(), "invalid proto message get")
}

func TestEmpty(t *testing.T) {
	st := &Status{}
	assert.Len(t, st.Details(), 0)

	st = nil
	assert.Len(t, st.Details(), 0)
}
