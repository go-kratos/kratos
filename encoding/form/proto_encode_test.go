package form

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/go-kratos/kratos/v2/internal/testdata/complex"
)

func TestEncodeValues(t *testing.T) {
	in := &complex.Complex{
		Id:          2233,
		NoOne:       "2233",
		Simple:      &complex.Simple{Component: "5566"},
		Simples:     []string{"3344", "5566"},
		B:           true,
		Sex:         complex.Sex_woman,
		Age:         18,
		A:           19,
		Count:       3,
		Price:       11.23,
		D:           22.22,
		Byte:        []byte("123"),
		Map:         map[string]string{"kratos": "https://go-kratos.dev/", "kratos_start": "https://go-kratos.dev/en/docs/getting-started/start/"},
		MapInt64Key: map[int64]string{1: "kratos", 2: "go-zero"},

		Timestamp: timestamppb.New(time.Date(1970, 1, 1, 0, 0, 20, 2, time.Local)),
		Duration:  &durationpb.Duration{Seconds: 120, Nanos: 22},
		Field:     &fieldmaskpb.FieldMask{Paths: []string{"1", "2"}},
		Double:    &wrapperspb.DoubleValue{Value: 12.33},
		Float:     &wrapperspb.FloatValue{Value: 12.34},
		Int64:     &wrapperspb.Int64Value{Value: 64},
		Int32:     &wrapperspb.Int32Value{Value: 32},
		Uint64:    &wrapperspb.UInt64Value{Value: 64},
		Uint32:    &wrapperspb.UInt32Value{Value: 32},
		Bool:      &wrapperspb.BoolValue{Value: false},
		String_:   &wrapperspb.StringValue{Value: "go-kratos"},
		Bytes:     &wrapperspb.BytesValue{Value: []byte("123")},
	}
	query, err := EncodeValues(in)
	if err != nil {
		t.Fatal(err)
	}
	want := "a=19&age=18&b=true&bool=false&byte=MTIz&bytes=MTIz&count=3&d=22.22&double=12.33&duration=2m0.000000022s&field=1%2C2&float=12.34&id=2233&int32=32&int64=64&map%5Bkratos%5D=https%3A%2F%2Fgo-kratos.dev%2F&map%5Bkratos_start%5D=https%3A%2F%2Fgo-kratos.dev%2Fen%2Fdocs%2Fgetting-started%2Fstart%2F&map_int64_key%5B1%5D=kratos&map_int64_key%5B2%5D=go-zero&numberOne=2233&price=11.23&sex=woman&simples=3344&simples=5566&string=go-kratos&timestamp=1970-01-01T00%3A00%3A20.000000002Z&uint32=32&uint64=64&very_simple.component=5566" // nolint:lll
	if got := query.Encode(); want != got {
		t.Errorf("\nwant: %s, \ngot: %s", want, got)
	}
}

func TestJsonCamelCase(t *testing.T) {
	tests := []struct {
		camelCase string
		snakeCase string
	}{
		{
			"userId", "user_id",
		},
		{
			"user", "user",
		},
		{
			"userIdAndUsername", "user_id_and_username",
		},
		{
			"", "",
		},
	}
	for _, test := range tests {
		t.Run(test.snakeCase, func(t *testing.T) {
			camel := jsonCamelCase(test.snakeCase)
			if camel != test.camelCase {
				t.Errorf("want: %s, got: %s", test.camelCase, camel)
			}
		})
	}
}

func TestIsASCIILower(t *testing.T) {
	tests := []struct {
		b     byte
		lower bool
	}{
		{
			'A', false,
		},
		{
			'a', true,
		},
		{
			',', false,
		},
		{
			'1', false,
		},
		{
			' ', false,
		},
	}
	for _, test := range tests {
		t.Run(string(test.b), func(t *testing.T) {
			lower := isASCIILower(test.b)
			if test.lower != lower {
				t.Errorf("'%s' is not ascii lower", string(test.b))
			}
		})
	}
}
