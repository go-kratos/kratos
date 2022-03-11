package json

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	testData "github.com/SeeMusic/kratos/v2/internal/testdata/encoding"
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
		{
			input:  &testData.TestModel{Id: 1, Name: "go-kratos", Hobby: []string{"1", "2"}},
			expect: `{"id":"1","name":"go-kratos","hobby":["1","2"],"attrs":{}}`,
		},
		{
			input:  &mock{value: Gopher},
			expect: `"gopher"`,
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
	p2 := testData.TestModel{}
	p3 := &testData.TestModel{}
	p4 := &mock{}
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
		{
			input:  `{"id":"1","name":"go-kratos","hobby":["1","2"],"attrs":{}}`,
			expect: &p2,
		},
		{
			input:  `{"id":1,"name":"go-kratos","hobby":["1","2"]}`,
			expect: &p3,
		},
		{
			input:  `"zebra"`,
			expect: p4,
		},
	}
	for _, v := range tests {
		want := []byte(v.input)
		err := (codec{}).Unmarshal(want, v.expect)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
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
