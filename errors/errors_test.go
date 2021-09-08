package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/grpc_testing"
)

func TestError(t *testing.T) {
	var base *Error
	err := Newf(http.StatusBadRequest, "reason", "message")
	err2 := Newf(http.StatusBadRequest, "reason", "message")
	err3 := err.WithMetadata(map[string]string{
		"foo": "bar",
	})
	werr := fmt.Errorf("wrap %w", err)

	if errors.Is(err, new(Error)) {
		t.Errorf("should not be equal: %v", err)
	}
	if !errors.Is(werr, err) {
		t.Errorf("should be equal: %v", err)
	}
	if !errors.Is(werr, err2) {
		t.Errorf("should be equal: %v", err)
	}

	if !errors.As(err, &base) {
		t.Errorf("should be matchs: %v", err)
	}
	if !IsBadRequest(err) {
		t.Errorf("should be matchs: %v", err)
	}

	if reason := Reason(err); reason != err3.Reason {
		t.Errorf("got %s want: %s", reason, err)
	}

	if err3.Metadata["foo"] != "bar" {
		t.Error("not expected metadata")
	}

	gs := err.GRPCStatus()
	se := FromError(gs.Err())
	if se.Reason != "reason" {
		t.Errorf("got %+v want %+v", se, err)
	}

	gs2, _ := status.New(codes.InvalidArgument, "bad request").WithDetails(&grpc_testing.Empty{})
	se2 := FromError(gs2.Err())
	// codes.InvalidArgument should convert to http.StatusBadRequest
	if se2.Code != http.StatusBadRequest {
		t.Errorf("convert code err, got %d want %d", UnknownCode, http.StatusBadRequest)
	}
}
