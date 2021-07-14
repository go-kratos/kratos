package form

import (
	"github.com/go-playground/form/v4"
	"testing"
)

type testStruct struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func Test_codec2_Unmarshal(t *testing.T) {

	s := &testStruct{}
	decoder := form.NewDecoder()
	decoder.SetTagName("json")
	encoder := form.NewEncoder()
	encoder.SetTagName("json")

	type fields struct {
		encoder *form.Encoder
		decoder *form.Decoder
	}
	type args struct {
		data []byte
		v    interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				encoder: encoder,
				decoder: decoder,
			},
			args: args{
				data: []byte(`
------WebKitFormBoundary8yfyAJD270aASP3u
Content-Disposition: form-data; name="field1"

asdasdasda
------WebKitFormBoundary8yfyAJD270aASP3u
Content-Disposition: form-data; name="field2"

2
------WebKitFormBoundary8yfyAJD270aASP3u--
`),
				v: s,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := multipartCodec{
				encoder: tt.fields.encoder,
				decoder: tt.fields.decoder,
			}
			if err := c.Unmarshal(tt.args.data, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(tt.args.v)
		})
	}
}

func Test_getBoundary(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name         string
		args         args
		wantBoundary string
	}{
		{
			name: "1",
			args: args{
				body: `
------WebKitFormBoundaryVAeYqYuN75Dk4s06
Content-Disposition: form-data; name="file"; filename="体脂秤的顾客 (1).xlsx"
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet


------WebKitFormBoundaryVAeYqYuN75Dk4s06--
`,
			},
			wantBoundary: "----WebKitFormBoundaryVAeYqYuN75Dk4s06",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBoundary := getBoundary(tt.args.body); gotBoundary != tt.wantBoundary {
				t.Errorf("getBoundary() = %v, want %v", gotBoundary, tt.wantBoundary)
			}
		})
	}
}
