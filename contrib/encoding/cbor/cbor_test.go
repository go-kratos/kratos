package cbor

import (
	"reflect"
	"testing"
)

func TestCodec_Name(t *testing.T) {
	c := new(codec)
	if !reflect.DeepEqual(c.Name(), "cbor") {
		t.Errorf("Name() should be cbor, but got %s", c.Name())
	}
}

func TestCodec_Marshal(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "default",
			args: args{v: struct{ String string }{String: "string"}},
			want: []byte{161, 102, 83, 116, 114, 105, 110, 103, 102, 115, 116, 114, 105, 110, 103},
		},
		{
			name: "tag: keyasint",
			args: args{v: struct {
				String string `cbor:"1,keyasint"`
			}{String: "string"}},
			want: []byte{161, 1, 102, 115, 116, 114, 105, 110, 103},
		},
		{
			name: "tag: toarray",
			args: args{v: struct {
				Ints []int `cbor:",toarray"`
			}{Ints: []int{1, 2, 3}}},
			want: []byte{161, 100, 73, 110, 116, 115, 131, 1, 2, 3},
		},
		{
			name: "tag: omitempty",
			args: args{v: struct {
				Empty []int `cbor:",toarray,omitempty"`
			}{Empty: []int{}}},
			want: []byte{160},
		},
		{
			name:    "invalid: unsupported type",
			args:    args{v: func() {}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(codec)
			got, err := c.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodec_Unmarshal(t *testing.T) {
	type args struct {
		data []byte
		v    interface{}
	}
	type want struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				data: []byte{161, 102, 83, 116, 114, 105, 110, 103, 102, 115, 116, 114, 105, 110, 103},
				v:    &struct{ String string }{},
			},
			want: want{
				v: &struct {
					String string
				}{String: "string"},
			},
		},
		{
			name: "tag: keyasint",
			args: args{
				data: []byte{161, 1, 102, 115, 116, 114, 105, 110, 103},
				v: &struct {
					String string `cbor:"1,keyasint"`
				}{},
			},
			want: want{
				v: &struct {
					String string `cbor:"1,keyasint"`
				}{String: "string"},
			},
		},
		{
			name: "tag: toarray",
			args: args{
				data: []byte{161, 100, 73, 110, 116, 115, 131, 1, 2, 3},
				v: &struct {
					Ints []int `cbor:",toarray"`
				}{},
			},
			want: want{
				v: &struct {
					Ints []int `cbor:",toarray"`
				}{Ints: []int{1, 2, 3}},
			},
		},
		{
			name: "tag: omitempty",
			args: args{
				data: []byte{160},
				v: &struct {
					Empty []int `cbor:",toarray,omitempty"`
				}{},
			},
			want: want{
				v: &struct {
					Empty []int `cbor:",toarray,omitempty"`
				}{Empty: nil},
			},
		},
		{
			name: "invalid: malformed data",
			args: args{
				data: []byte{160},
				v:    nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(codec)
			err := c.Unmarshal(tt.args.data, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.want.v, tt.args.v) {
				t.Errorf("Unmarshal() got = %v, want %v", tt.args.v, tt.want.v)
			}
		})
	}
}
