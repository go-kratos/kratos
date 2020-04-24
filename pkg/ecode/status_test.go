package ecode

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/timestamp"
	pkgerr "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/go-kratos/kratos/pkg/ecode/types"
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

func TestParse(t *testing.T) {
	var (
		errMsg string
		st     *Status
	)
	t.Run("parse ecode.Status", func(t *testing.T) {
		st = Parse(FromCode(RequestErr))
		assert.Equal(t, RequestErr.Code(), st.Code())

		errMsg = "name is required"
		st = Parse(Error(RequestErr, errMsg))
		assert.Equal(t, RequestErr.Code(), st.Code())
		assert.Equal(t, errMsg, st.Message())
	})

	t.Run("parse ecode.Code", func(t *testing.T) {
		st = Parse(ServerErr)

		assert.Equal(t, ServerErr.Error(), st.Message())
	})

	t.Run("parse wrap error", func(t *testing.T) {
		st = Parse(pkgerr.Wrap(ServerErr, "db is unavailable"))

		assert.Equal(t, ServerErr.Code(), st.Code())
	})

	t.Run("parse general error", func(t *testing.T) {
		errMsg = "something is wrong!"

		st = Parse(errors.New(errMsg))
		assert.Equal(t, errMsg, st.Message())
		assert.Equal(t, ServerErr.Code(), st.Code())
	})

	t.Run("parse raw Canceled", func(t *testing.T) {
		st = Parse(context.Canceled)

		assert.Equal(t, Canceled.Code(), st.Code())
		assert.Equal(t, Canceled.Message(), st.Message())
	})

	t.Run("parse raw DeadlineExceeded", func(t *testing.T) {
		st = Parse(context.DeadlineExceeded)

		assert.Equal(t, Deadline.Code(), st.Code())
		assert.Equal(t, Deadline.Message(), st.Message())
	})

	t.Run("parse google grpc status", func(t *testing.T) {
		errMsg = "record not found"

		st = Parse(status.Error(codes.NotFound, errMsg))

		assert.Equal(t, int(codes.NotFound), st.Code())
		assert.Equal(t, errMsg, st.Message())
	})

	t.Run("parse nil", func(t *testing.T) {
		st = Parse(nil)

		assert.Equal(t, OK.Code(), st.Code())
	})
}
