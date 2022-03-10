package errors

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
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
		t.Errorf("should be matches: %v", err)
	}
	if !IsBadRequest(err) {
		t.Errorf("should be matches: %v", err)
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
	if FromError(nil) != nil {
		t.Errorf("FromError(nil) should be nil")
	}
	e := FromError(errors.New("test"))
	if !reflect.DeepEqual(e.Code, int32(UnknownCode)) {
		t.Errorf("no expect value: %v, but got: %v", e.Code, int32(UnknownCode))
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name string
		e    *Error
		err  error
		want bool
	}{
		{
			name: "true",
			e:    &Error{Code: 404, Reason: "test"},
			err:  New(http.StatusNotFound, "test", ""),
			want: true,
		},
		{
			name: "false",
			e:    &Error{Reason: "test"},
			err:  errors.New("test"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.e.Is(tt.err); ok != tt.want {
				t.Errorf("Error.Error() = %v, want %v", ok, tt.want)
			}
		})
	}
}

func TestOther(t *testing.T) {
	if !reflect.DeepEqual(Code(nil), 200) {
		t.Errorf("Code(nil) = %v, want %v", Code(nil), 200)
	}
	if !reflect.DeepEqual(Code(errors.New("test")), UnknownCode) {
		t.Errorf(`Code(errors.New("test")) = %v, want %v`, Code(nil), 200)
	}
	if !reflect.DeepEqual(Reason(errors.New("test")), UnknownReason) {
		t.Errorf(`Reason(errors.New("test")) = %v, want %v`, Reason(nil), UnknownReason)
	}
	err := Errorf(10001, "test code 10001", "message")
	if !reflect.DeepEqual(Code(err), 10001) {
		t.Errorf(`Code(err) = %v, want %v`, Code(err), 10001)
	}
	if !reflect.DeepEqual(Reason(err), "test code 10001") {
		t.Errorf(`Reason(err) = %v, want %v`, Reason(err), "test code 10001")
	}
}
