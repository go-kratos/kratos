package status

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestToGRPCCode(t *testing.T) {
	tests := []struct {
		name string
		code int
		want codes.Code
	}{
		{"http.StatusOK", http.StatusOK, codes.OK},
		{"http.StatusBadRequest", http.StatusBadRequest, codes.InvalidArgument},
		{"http.StatusUnauthorized", http.StatusUnauthorized, codes.Unauthenticated},
		{"http.StatusForbidden", http.StatusForbidden, codes.PermissionDenied},
		{"http.StatusNotFound", http.StatusNotFound, codes.NotFound},
		{"http.StatusConflict", http.StatusConflict, codes.Aborted},
		{"http.StatusTooManyRequests", http.StatusTooManyRequests, codes.ResourceExhausted},
		{"http.StatusInternalServerError", http.StatusInternalServerError, codes.Internal},
		{"http.StatusNotImplemented", http.StatusNotImplemented, codes.Unimplemented},
		{"http.StatusServiceUnavailable", http.StatusServiceUnavailable, codes.Unavailable},
		{"http.StatusGatewayTimeout", http.StatusGatewayTimeout, codes.DeadlineExceeded},
		{"StatusClientClosed", ClientClosed, codes.Canceled},
		{"else", 100000, codes.Unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToGRPCCode(tt.code); got != tt.want {
				t.Errorf("GRPCCodeFromStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromGRPCCode(t *testing.T) {
	tests := []struct {
		name string
		code codes.Code
		want int
	}{
		{"codes.OK", codes.OK, http.StatusOK},
		{"codes.Canceled", codes.Canceled, ClientClosed},
		{"codes.Unknown", codes.Unknown, http.StatusInternalServerError},
		{"codes.InvalidArgument", codes.InvalidArgument, http.StatusBadRequest},
		{"codes.DeadlineExceeded", codes.DeadlineExceeded, http.StatusGatewayTimeout},
		{"codes.NotFound", codes.NotFound, http.StatusNotFound},
		{"codes.AlreadyExists", codes.AlreadyExists, http.StatusConflict},
		{"codes.PermissionDenied", codes.PermissionDenied, http.StatusForbidden},
		{"codes.Unauthenticated", codes.Unauthenticated, http.StatusUnauthorized},
		{"codes.ResourceExhausted", codes.ResourceExhausted, http.StatusTooManyRequests},
		{"codes.FailedPrecondition", codes.FailedPrecondition, http.StatusBadRequest},
		{"codes.Aborted", codes.Aborted, http.StatusConflict},
		{"codes.OutOfRange", codes.OutOfRange, http.StatusBadRequest},
		{"codes.Unimplemented", codes.Unimplemented, http.StatusNotImplemented},
		{"codes.Internal", codes.Internal, http.StatusInternalServerError},
		{"codes.Unavailable", codes.Unavailable, http.StatusServiceUnavailable},
		{"codes.DataLoss", codes.DataLoss, http.StatusInternalServerError},
		{"else", codes.Code(10000), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromGRPCCode(tt.code); got != tt.want {
				t.Errorf("StatusFromGRPCCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
