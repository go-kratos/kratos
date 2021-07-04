package testproto

import (
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/form"
	"google.golang.org/genproto/protobuf/field_mask"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func TestTypes(t *testing.T) {
	var (
		timeT  = time.Date(2016, time.December, 15, 12, 23, 32, 49, time.UTC)
		timePb = timestamppb.New(timeT)

		durationT  = 13 * time.Hour
		durationPb = durationpb.New(durationT)

		fieldmaskPb = &field_mask.FieldMask{Paths: []string{"float_value", "double_value"}}

		input = &Proto3Message{
			FloatValue:         1.5,
			DoubleValue:        2.5,
			Int64Value:         -1,
			Int32Value:         -2,
			Uint64Value:        3,
			Uint32Value:        4,
			BoolValue:          true,
			StringValue:        "str",
			BytesValue:         []byte("abc123!?$*&()'-=@~"),
			RepeatedValue:      []string{"a", "b", "c"},
			RepeatedMessage:    []*wrapperspb.UInt64Value{{Value: 1}, {Value: 2}, {Value: 3}},
			EnumValue:          EnumValue_Y,
			RepeatedEnum:       []EnumValue{EnumValue_Y, EnumValue_Z, EnumValue_X},
			TimestampValue:     timePb,
			DurationValue:      durationPb,
			FieldmaskValue:     fieldmaskPb,
			WrapperFloatValue:  &wrapperspb.FloatValue{Value: 1.5},
			WrapperDoubleValue: &wrapperspb.DoubleValue{Value: 2.5},
			WrapperInt64Value:  &wrapperspb.Int64Value{Value: -1},
			WrapperInt32Value:  &wrapperspb.Int32Value{Value: -2},
			WrapperUInt64Value: &wrapperspb.UInt64Value{Value: 3},
			WrapperUInt32Value: &wrapperspb.UInt32Value{Value: 4},
			WrapperBoolValue:   &wrapperspb.BoolValue{Value: true},
			WrapperStringValue: &wrapperspb.StringValue{Value: "str"},
			WrapperBytesValue:  &wrapperspb.BytesValue{Value: []byte("abc123!?$*&()'-=@~")},
		}
		want = &Proto3Message{}
	)
	codec := encoding.GetCodec(form.Name)
	data, err := codec.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	if err = codec.Unmarshal(data, want); err != nil {
		t.Fatal(err)
	}
}
