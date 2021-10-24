package toml

import (
	"reflect"
	"testing"
	"time"
)

type tomlConfig struct {
	Title   string
	Owner   ownerInfo
	Servers map[string]server
	Clients clients
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
	Bio  string
	DOB  time.Time
}

type server struct {
	IP string
	DC string
}

type clients struct {
	Data  [][]interface{}
	Hosts []string
}

func TestCodec_Marshal(t *testing.T) {
	testCases := []struct {
		input interface{}
		want  string
	}{
		{
			input: struct{}{},
			want:  "",
		},
		{
			input: server{
				DC: "dc",
				IP: "ip",
			},
			want: `DC = "dc"
IP = "ip"
`,
		},
	}

	for _, testCase := range testCases {
		b, err := (codec{}).Marshal(testCase.input)
		if err != nil {
			t.Fatalf("(codec{}).Marshal return err:%v\n", err)
		}
		if string(b) != testCase.want {
			t.Fatalf("(codec{}).Marshal return not match want, input:%v\n", testCase.input)
		}
	}

}

func TestCodec_Unmarshal(t *testing.T) {
	testCases := []struct {
		input string
		want  interface{}
	}{
		{
			input: `
title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00

[servers]
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"
  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ]`,
			want: tomlConfig{},
		},
	}

	for _, testCase := range testCases {
		vt := reflect.ValueOf(testCase.want).Type()
		dest := reflect.New(vt).Interface()
		err := (codec{}).Unmarshal([]byte(testCase.input), &dest)
		if err != nil {
			t.Fatalf("(codec{}).Marshal return err:%v\n", err)
		}
	}
}
