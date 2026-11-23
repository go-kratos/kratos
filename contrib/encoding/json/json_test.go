package json

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type testEmbed struct {
	Level1a int `json:"a"`
	Level1b int `json:"b"`
	Level1c int `json:"c"`
}

type testMessage struct {
	Field1 string     `json:"a"`
	Field2 string     `json:"b"`
	Field3 string     `json:"c"`
	Embed  *testEmbed `json:"embed,omitempty"`
}

type mock struct {
	value int
}

const (
	Unknown = iota
	Gopher
	Zebra
)

func (a *mock) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		a.value = Unknown
	case "gopher":
		a.value = Gopher
	case "zebra":
		a.value = Zebra
	}

	return nil
}

func (a *mock) MarshalJSON() ([]byte, error) {
	var s string
	switch a.value {
	default:
		s = "unknown"
	case Gopher:
		s = "gopher"
	case Zebra:
		s = "zebra"
	}

	return json.Marshal(s)
}

func TestJSON_Marshal(t *testing.T) {
	tests := []struct {
		input  any
		expect string
	}{
		{
			input:  &testMessage{},
			expect: `{"a":"","b":"","c":""}`,
		},
		{
			input:  &testMessage{Field1: "a", Field2: "b", Field3: "c"},
			expect: `{"a":"a","b":"b","c":"c"}`,
		},
		{
			input:  &mock{value: Gopher},
			expect: `"gopher"`,
		},
		{
			input:  wrapperspb.Int64(1),
			expect: `"1"`,
		},
	}
	for _, v := range tests {
		data, err := (codec{}).Marshal(v.input)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
		}
		if got, want := string(data), v.expect; strings.ReplaceAll(got, " ", "") != want {
			if strings.Contains(want, "\n") {
				t.Errorf("marshal(%#v):\nHAVE:\n%s\nWANT:\n%s", v.input, got, want)
			} else {
				t.Errorf("marshal(%#v):\nhave %#q\nwant %#q", v.input, got, want)
			}
		}
	}
}

func TestJSON_Unmarshal(t *testing.T) {
	p := testMessage{}
	p2 := &mock{}
	p3 := wrapperspb.Int64(0)
	tests := []struct {
		input  string
		expect any
	}{
		{
			input:  `{"a":"","b":"","c":""}`,
			expect: &testMessage{},
		},
		{
			input:  `{"a":"a","b":"b","c":"c"}`,
			expect: &p,
		},
		{
			input:  `"zebra"`,
			expect: p2,
		},
		{
			input:  `"2"`,
			expect: p3,
		},
	}
	for _, v := range tests {
		want := []byte(v.input)
		err := (codec{}).Unmarshal(want, v.expect)
		if err != nil {
			t.Errorf("unmarshal(%#v): %s", v.input, err)
		}
		got, err := codec{}.Marshal(v.expect)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
		}
		if !reflect.DeepEqual(strings.ReplaceAll(string(got), " ", ""), strings.ReplaceAll(string(want), " ", "")) {
			t.Errorf("marshal(%#v):\nhave %#q\nwant %#q", v.input, got, want)
		}
	}
}

func TestJSON_UnmarshalPointerToProtoMessage(t *testing.T) {
	var msg *wrapperspb.StringValue
	if err := (codec{}).Unmarshal([]byte(`"kratos"`), &msg); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if msg == nil {
		t.Fatal("Unmarshal() message = nil, want allocated message")
	}
	if got, want := msg.Value, "kratos"; got != want {
		t.Fatalf("Unmarshal() value = %s, want %s", got, want)
	}
}

func TestJSON_EmptyData(t *testing.T) {
	var msg wrapperspb.StringValue
	if err := (codec{}).Unmarshal(nil, &msg); err != nil {
		t.Fatalf("Unmarshal(nil) error = %v", err)
	}
}

func TestJSON_Name(t *testing.T) {
	if got, want := (codec{}).Name(), Name; got != want {
		t.Fatalf("Name() = %s, want %s", got, want)
	}
}
