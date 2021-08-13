package httputil

import (
	"google.golang.org/grpc/codes"
	"net/http"
	"testing"
)

func TestContentSubtype(t *testing.T) {
	tests := []struct {
		contentType string
		want        string
	}{
		{"text/html; charset=utf-8", "html"},
		{"multipart/form-data; boundary=something", "form-data"},
		{"application/json; charset=utf-8", "json"},
		{"application/json", "json"},
		{"application/xml", "xml"},
		{"text/xml", "xml"},
		{";text/xml", ""},
		{"application", ""},
	}
	for _, test := range tests {
		t.Run(test.contentType, func(t *testing.T) {
			got := ContentSubtype(test.contentType)
			if got != test.want {
				t.Fatalf("want %v got %v", test.want, got)
			}
		})
	}
}

func TestGRPCCodeFromStatus(t *testing.T) {

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
		{"StatusClientClosed", StatusClientClosed, codes.Canceled},
		{"else", 100000, codes.Unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GRPCCodeFromStatus(tt.code); got != tt.want {
				t.Errorf("GRPCCodeFromStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusFromGRPCCode(t *testing.T) {
	tests := []struct {
		name string
		code codes.Code
		want int
	}{
		{"codes.OK", codes.OK, http.StatusOK},
		{"codes.Canceled", codes.Canceled, StatusClientClosed},
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
			if got := StatusFromGRPCCode(tt.code); got != tt.want {
				t.Errorf("StatusFromGRPCCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentType(t *testing.T) {

	tests := []struct {
		name string
		subtype string
		want string
	}{
		{"kratos","kratos","application/kratos"},
		{"json","json","application/json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContentType(tt.subtype); got != tt.want {
				t.Errorf("ContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}