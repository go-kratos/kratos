package json

import (
	"bytes"
	"strings"
	"testing"
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

func TestJSON_Marshal(t *testing.T) {
	tests := []struct {
		input  interface{}
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
	}
	for _, v := range tests {
		data, err := (codec{}).Marshal(v.input)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
		}
		if got, want := string(data), v.expect; got != want {
			if strings.Contains(want, "\n") {
				t.Errorf("marshal(%#v):\nHAVE:\n%s\nWANT:\n%s", v.input, got, want)
			} else {
				t.Errorf("marshal(%#v):\nhave %#q\nwant %#q", v.input, got, want)
			}
		}
	}
}

func TestJSON_Unmarshal(t *testing.T) {
	p := &testMessage{}
	tests := []struct {
		input  string
		expect interface{}
	}{
		{
			input:  `{"a":"","b":"","c":""}`,
			expect: &testMessage{},
		},
		{
			input:  `{"a":"a","b":"b","c":"c"}`,
			expect: &p,
		},
	}
	for _, v := range tests {
		want := []byte(v.input)
		err := (codec{}).Unmarshal(want, &v.expect)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
		}
		got, err := codec{}.Marshal(v.expect)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("marshal(%#v):\nhave %#q\nwant %#q", v.input, got, want)
		}
	}
}
