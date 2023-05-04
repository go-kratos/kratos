package toml

import (
	"reflect"
	"testing"
	"time"
)

func TestCodec_Name(t *testing.T) {
	if (codec{}).Name() != Name {
		t.Fatal("(codec{}).Name() should be consistent with Name")
	}
}

func TestCodec_Unmarshal(t *testing.T) {
	type User struct {
		Name  string `toml:"name"`
		Email string `toml:"email"`
	}
	tests := []struct {
		data  string
		value interface{}
	}{
		{
			data:  "v = \"John Doe\"",
			value: map[string]interface{}{"v": "John Doe"},
		},
		{
			data:  "v = 100",
			value: map[string]interface{}{"v": 100},
		},
		{
			data:  "v = 100_100",
			value: map[string]interface{}{"v": 100100},
		},
		{
			data:  "v = 1.1",
			value: map[string]interface{}{"v": 1.1},
		},
		{
			data:  "v = true",
			value: map[string]interface{}{"v": true},
		},
		{
			data:  "v = [\"apple\", \"banana\", \"cherry\"]",
			value: map[string]interface{}{"": []string{"apple", "banana", "cherry"}},
		},
		{
			data:  "v = 1.618e+01",
			value: map[string]interface{}{"v": 16.18},
		},
		{
			data:  "v = 0xDEADBEEF",
			value: map[string]interface{}{"v": 3735928559},
		},
		{
			data:  "v = 0b1101_0101",
			value: map[string]interface{}{"v": 213},
		},
		{
			data:  "v = 0o755",
			value: map[string]interface{}{"v": 493},
		},
		{
			data:  "v = 2022-04-16T12:13:14Z",
			value: map[string]interface{}{"v": time.Date(2022, time.April, 16, 12, 13, 14, 0, time.UTC)},
		},
		{
			data: "v = { name = \"John\", email = \"john.doe@example.com\"}",
			value: map[string]interface{}{"v": User{
				Name:  "John",
				Email: "john.doe@example.com",
			}},
		},
	}

	for _, tt := range tests {
		v := reflect.ValueOf(tt.value).Type()
		value := reflect.New(v)
		err := (codec{}).Unmarshal([]byte(tt.data), value.Interface())
		if err != nil {
			t.Fatalf("(codec{}).Unmarshal should not return err: %v", err)
		}
	}
}

func TestCodec_Marshal(t *testing.T) {
	value := map[string]string{"v": "hi"}
	got, err := (codec{}).Marshal(value)
	if err != nil {
		t.Fatalf("should not return err")
	}
	//t.Logf("got: %v", string(got))
	if string(got) != "v = \"hi\"\n" {
		t.Fatalf("want v = \"hi\", return %s", string(got))
	}
}
