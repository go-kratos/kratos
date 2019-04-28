package ecode

import (
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go-common/library/ecode/internal/types"
)

// Error new status with code and message
func Error(code Code, message string) *Status {
	return &Status{s: &types.Status{Code: int32(code.Code()), Message: message}}
}

// Errorf new status with code and message
func Errorf(code Code, format string, args ...interface{}) *Status {
	return Error(code, fmt.Sprintf(format, args...))
}

var _ Codes = &Status{}

// Status statusError is an alias of a status proto
// implement ecode.Codes
type Status struct {
	s *types.Status
}

// Error implement error
func (s *Status) Error() string {
	return s.Message()
}

// Code return error code
func (s *Status) Code() int {
	return int(s.s.Code)
}

// Message return error message for developer
func (s *Status) Message() string {
	if s.s.Message == "" {
		return strconv.Itoa(int(s.s.Code))
	}
	return s.s.Message
}

// Details return error details
func (s *Status) Details() []interface{} {
	if s == nil || s.s == nil {
		return nil
	}
	details := make([]interface{}, 0, len(s.s.Details))
	for _, any := range s.s.Details {
		detail := &ptypes.DynamicAny{}
		if err := ptypes.UnmarshalAny(any, detail); err != nil {
			details = append(details, err)
			continue
		}
		details = append(details, detail.Message)
	}
	return details
}

// WithDetails WithDetails
func (s *Status) WithDetails(pbs ...proto.Message) (*Status, error) {
	for _, pb := range pbs {
		anyMsg, err := ptypes.MarshalAny(pb)
		if err != nil {
			return s, err
		}
		s.s.Details = append(s.s.Details, anyMsg)
	}
	return s, nil
}

// Equal for compatible.
// Deprecated: please use ecode.EqualError.
func (s *Status) Equal(err error) bool {
	return EqualError(s, err)
}

// Proto return origin protobuf message
func (s *Status) Proto() *types.Status {
	return s.s
}

// FromCode create status from ecode
func FromCode(code Code) *Status {
	return &Status{s: &types.Status{Code: int32(code)}}
}

// FromProto new status from grpc detail
func FromProto(pbMsg proto.Message) Codes {
	if msg, ok := pbMsg.(*types.Status); ok {
		if msg.Message == "" {
			// NOTE: if message is empty convert to pure Code, will get message from config center.
			return Code(msg.Code)
		}
		return &Status{s: msg}
	}
	return Errorf(ServerErr, "invalid proto message get %v", pbMsg)
}
