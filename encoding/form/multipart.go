package form

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"reflect"
	"regexp"

	"github.com/go-kratos/kratos/v2/encoding"

	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

// Name is form codec name
const multipartName = "form-data"

func init() {
	decoder := form.NewDecoder()
	decoder.SetTagName("json")
	encoder := form.NewEncoder()
	encoder.SetTagName("json")
	encoding.RegisterCodec(multipartCodec{encoder: encoder, decoder: decoder})
}

type multipartCodec struct {
	encoder *form.Encoder
	decoder *form.Decoder
}

func (c multipartCodec) Marshal(v interface{}) ([]byte, error) {
	var vs url.Values
	var err error
	if m, ok := v.(proto.Message); ok {
		vs, err = EncodeMap(m)
		if err != nil {
			return nil, err
		}
	} else {
		vs, err = c.encoder.Encode(v)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return []byte(vs.Encode()), nil
}

func (c multipartCodec) Unmarshal(data []byte, v interface{}) error {

	boundary := getBoundary(string(data))

	mr := multipart.NewReader(bytes.NewReader(data), boundary)

	f, err := mr.ReadForm(32 << 20)
	if err != nil {
		return err
	}

	values := f.Value

	for key, headers := range f.File {
		file, err := headers[0].Open()

		if err != nil {
			continue
		}

		bts, err := ioutil.ReadAll(file)

		file.Close()

		if err != nil {
			continue
		}

		values[key] = append(values[key], base64.StdEncoding.EncodeToString(bts))
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if m, ok := v.(proto.Message); ok {
		return MapProto(m, f.Value)
	} else if m, ok := reflect.Indirect(reflect.ValueOf(v)).Interface().(proto.Message); ok {
		return MapProto(m, f.Value)
	}

	return c.decoder.Decode(v, f.Value)
}

func (multipartCodec) Name() string {
	return multipartName
}

// 获取form-data 边界
func getBoundary(body string) (boundary string) {
	reg, err := regexp.Compile("--(?P<boundary>[^\\s]+)(--)?\\s*")
	if err != nil {
		return
	}
	s := reg.FindAllString(body, -1)

	m := make(map[string]int)

	max := 0
	for i := range s {
		k := reg.ReplaceAllString(s[i], "${boundary}")
		m[k]++
		if max < m[k] {
			max = m[k]
			boundary = k
		}
	}
	return
}
