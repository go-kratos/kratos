package apollo

import (
	"reflect"
	"testing"
)

func Test_genKey(t *testing.T) {
	type args struct {
		ns  string
		sub string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				ns:  "",
				sub: "has_no_ns",
			},
			want: "has_no_ns",
		},
		{
			name: "case 2",
			args: args{
				ns:  "ns.ext",
				sub: "sub",
			},
			want: "ns.sub",
		},
		{
			name: "case 3",
			args: args{
				ns:  "",
				sub: "",
			},
			want: "",
		},
		{
			name: "case 4",
			args: args{
				ns:  "ns.ext",
				sub: "sub.sub2.sub3",
			},
			want: "ns.sub.sub2.sub3",
		},
		{
			name: "case 5",
			args: args{
				ns:  "ns.more.ext",
				sub: "sub.sub2.sub3",
			},
			want: "ns.more.sub.sub2.sub3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genKey(tt.args.ns, tt.args.sub); got != tt.want {
				t.Errorf("genKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_format(t *testing.T) {
	type args struct {
		ns string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 0",
			args: args{
				ns: "ns.yaml",
			},
			want: "yaml",
		},
		{
			name: "case 1",
			args: args{
				ns: "ns",
			},
			want: "json",
		},
		{
			name: "case 2",
			args: args{
				ns: "ns.more.json",
			},
			want: "json",
		},
		{
			name: "case 3",
			args: args{
				ns: "",
			},
			want: "json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.args.ns); got != tt.want {
				t.Errorf("format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertProperties(t *testing.T) {
	type args struct {
		key    string
		value  interface{}
		target map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "case 0",
			args: args{
				key:    "application.name",
				value:  "app name",
				target: map[string]interface{}{},
			},
			want: map[string]interface{}{
				"application": map[string]interface{}{
					"name": "app name",
				},
			},
		},
		{
			name: "case 1",
			args: args{
				key:    "application",
				value:  []string{"1", "2", "3"},
				target: map[string]interface{}{},
			},
			want: map[string]interface{}{
				"application": []string{"1", "2", "3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolve(tt.args.key, tt.args.value, tt.args.target)
			if !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("convertProperties() = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}

func Test_convertProperties_duplicate(t *testing.T) {
	target := map[string]interface{}{}
	resolve("application.name", "name", target)
	_, ok := target["application"]
	if !reflect.DeepEqual(ok, true) {
		t.Errorf("ok = %v, want %v", ok, true)
	}
	_, ok = target["application"].(map[string]interface{})["name"]
	if !reflect.DeepEqual(ok, true) {
		t.Errorf("ok = %v, want %v", ok, true)
	}
	if !reflect.DeepEqual(target["application"].(map[string]interface{})["name"], "name") {
		t.Errorf("target[\"application\"][\"name\"] = %v, want %v", target["application"].(map[string]interface{})["name"], "name")
	}

	// cause duplicate, the oldest value will be kept
	resolve("application.name.first", "first name", target)
	_, ok = target["application"]
	if !reflect.DeepEqual(ok, true) {
		t.Errorf("ok = %v, want %v", ok, true)
	}
	_, ok = target["application"].(map[string]interface{})["name"]
	if !reflect.DeepEqual(ok, true) {
		t.Errorf("ok = %v, want %v", ok, true)
	}
	if !reflect.DeepEqual(target["application"].(map[string]interface{})["name"], "name") {
		t.Errorf("target[\"application\"][\"name\"] = %v, want %v", target["application"].(map[string]interface{})["name"], "name")
	}
}
