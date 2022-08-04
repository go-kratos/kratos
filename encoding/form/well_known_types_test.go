package form

import (
	"encoding/base64"
	"testing"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestMarshalTimeStamp(t *testing.T) {
	tests := []struct {
		input  *timestamppb.Timestamp
		expect string
	}{
		{
			input:  timestamppb.New(time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC)),
			expect: "2022-01-02T03:04:05.000000006Z",
		},
		{
			input:  timestamppb.New(time.Date(2022, 13, 1, 13, 61, 61, 100, time.UTC)),
			expect: "2023-01-01T14:02:01.000000100Z",
		},
	}
	for _, v := range tests {
		content, err := marshalTimestamp(v.input.ProtoReflect())
		if err != nil {
			t.Errorf("expect %v,got %v", nil, err)
		}
		if got, want := content, v.expect; got != want {
			t.Errorf("expect %v,got %v", want, got)
		}
	}
}

func TestMarshalDuration(t *testing.T) {
	tests := []struct {
		input  *durationpb.Duration
		expect string
	}{
		{
			input:  durationpb.New(time.Duration(1<<63 - 1)),
			expect: "2562047h47m16.854775807s",
		},
		{
			input:  durationpb.New(time.Duration(-1 << 63)),
			expect: "-2562047h47m16.854775808s",
		},
		{
			input:  durationpb.New(100 * time.Second),
			expect: "1m40s",
		},
		{
			input:  durationpb.New(-100 * time.Second),
			expect: "-1m40s",
		},
	}
	for _, v := range tests {
		content, err := marshalDuration(v.input.ProtoReflect())
		if err != nil {
			t.Errorf("expect %v,got %v", nil, err)
		}
		if got, want := content, v.expect; got != want {
			t.Errorf("expect %v,got %v", want, got)
		}
	}
}

func TestMarshalBytes(t *testing.T) {
	tests := []struct {
		input  protoreflect.Message
		expect string
	}{
		{
			input:  wrapperspb.Bytes([]byte("abc123!?$*&()'-=@~")).ProtoReflect(),
			expect: base64.StdEncoding.EncodeToString([]byte("abc123!?$*&()'-=@~")),
		},
		{
			input:  wrapperspb.Bytes([]byte("kratos")).ProtoReflect(),
			expect: base64.StdEncoding.EncodeToString([]byte("kratos")),
		},
	}
	for _, v := range tests {
		content, err := marshalBytes(v.input)
		if err != nil {
			t.Errorf("expect %v,got %v", nil, err)
		}
		if got, want := content, v.expect; got != want {
			t.Errorf("expect %v,got %v", want, got)
		}
	}
}
