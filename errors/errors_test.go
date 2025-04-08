package errors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	text "text/template"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestError struct{ message string }

func (e *TestError) Error() string { return e.message }

func TestErrors(t *testing.T) {
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

	gs2 := status.New(codes.InvalidArgument, "bad request")
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
			e:    New(404, "test", ""),
			err:  New(http.StatusNotFound, "test", ""),
			want: true,
		},
		{
			name: "false",
			e:    New(0, "test", ""),
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

func TestCause(t *testing.T) {
	testError := &TestError{message: "test"}
	err := BadRequest("foo", "bar").WithCause(testError)
	if !errors.Is(err, testError) {
		t.Fatalf("want %v but got %v", testError, err)
	}
	if te := new(TestError); errors.As(err, &te) {
		if te.message != testError.message {
			t.Fatalf("want %s but got %s", testError.message, te.message)
		}
	}
}

func TestOther(t *testing.T) {
	err := Errorf(10001, "test code 10001", "message")
	// Code
	if !reflect.DeepEqual(Code(nil), 200) {
		t.Errorf("Code(nil) = %v, want %v", Code(nil), 200)
	}
	if !reflect.DeepEqual(Code(errors.New("test")), UnknownCode) {
		t.Errorf(`Code(errors.New("test")) = %v, want %v`, Code(nil), 200)
	}
	if !reflect.DeepEqual(Code(err), 10001) {
		t.Errorf(`Code(err) = %v, want %v`, Code(err), 10001)
	}
	// Reason
	if !reflect.DeepEqual(Reason(nil), UnknownReason) {
		t.Errorf(`Reason(nil) = %v, want %v`, Reason(nil), UnknownReason)
	}
	if !reflect.DeepEqual(Reason(errors.New("test")), UnknownReason) {
		t.Errorf(`Reason(errors.New("test")) = %v, want %v`, Reason(nil), UnknownReason)
	}
	if !reflect.DeepEqual(Reason(err), "test code 10001") {
		t.Errorf(`Reason(err) = %v, want %v`, Reason(err), "test code 10001")
	}
	// Clone
	err400 := Newf(http.StatusBadRequest, "BAD_REQUEST", "param invalid")
	err400.Metadata = map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	if cerr := Clone(err400); cerr == nil || cerr.Error() != err400.Error() {
		t.Errorf("Clone(err) = %v, want %v", Clone(err400), err400)
	}
	if cerr := Clone(nil); cerr != nil {
		t.Errorf("Clone(nil) = %v, want %v", Clone(err400), err400)
	}
}

// MockI18n is a mock implementation of the I18n interface
type MockI18n struct {
	localizeFunc func(ctx context.Context, reason string, data any) string
}

// Localize mocks the Localize method of the I18n interface
func (m *MockI18n) Localize(ctx context.Context, reason string, data any) string {
	return m.localizeFunc(ctx, reason, data)
}

func textFormatter(format string, data any) (s string, err error) {
	t, err := text.New("text/template").Parse(format)
	if err != nil {
		return "", nil
	}
	tmpl := text.Must(t, err)
	resultIoWriter := new(strings.Builder)

	if err = tmpl.Execute(resultIoWriter, data); err != nil {
		return "", err
	}
	return resultIoWriter.String(), nil
}

// TestNewWithContext_GlobalI18nIsNil tests the NewWithContext function when globalI18n is nil.
// It verifies that the returned Error object has an empty Message field.
func TestNewWithContext_GlobalI18nIsNil(t *testing.T) {
	// Arrange
	ctx := context.Background()
	code := 404
	reason := "not_found"
	data := map[string]string{"resource": "user"}

	// Act
	err := NewWithContext(ctx, code, reason, data)

	// Assert
	if !errors.Is(err, New(code, reason, "")) {
		t.Errorf("Expected error to be %v, but got %v", New(code, reason, ""), err)
	}
	if err.Message != "" {
		t.Errorf("Expected error Message to be empty, but got %v", err.Message)
	}
}

type langContextKey struct{}

func withLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langContextKey{}, lang)
}

func getLang(ctx context.Context) string {
	lang, _ := ctx.Value(langContextKey{}).(string)
	return lang
}

// TestNewWithContext_GlobalI18nIsNotNil tests the NewWithContext function when globalI18n is not nil.
// It verifies that the returned Error object has a localized Message field.
func TestNewWithContext_GlobalI18nIsNotNil(t *testing.T) {
	// Arrange
	ctx := context.Background()
	code := 404
	reason := "not_found"
	data := map[string]string{"resource": "Kratos"}
	lang := map[string]string{
		"en_CH": "{{ .resource }} is the best",
		"zh_CN": "{{ .resource }}是最好的",
	}

	// Mock the globalI18n
	mockI18n := &MockI18n{
		localizeFunc: func(ctx context.Context, reason string, data any) string {
			message := lang[getLang(ctx)]
			s, err := textFormatter(message, data)
			if err != nil {
				return ""
			}
			return s
		},
	}
	RegisterI18nManager(mockI18n)

	// Act
	err := NewWithContext(withLang(ctx, "en_CH"), code, reason, data)
	// Assert
	if !errors.Is(err, New(code, reason, "")) {
		t.Errorf("Expected error to be %v, but got %v", New(code, reason, ""), err)
	}
	if err.Message != "kratos is the best" {
		t.Errorf("Expected error message to be %v, but got %v", "kratos is the best", err.Message)
	}

	err = NewWithContext(withLang(ctx, "zh_CN"), code, reason, data)
	if !errors.Is(err, New(code, reason, "")) {
		t.Errorf("Expected error to be %v, but got %v", New(code, reason, ""), err)
	}
	if err.Message != "kratos是最好的" {
		t.Errorf("Expected error message to be %v, but got %v", "kratos是最好的", err.Message)
	}
}
