package plain

import (
	"errors"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
)

/* Name is the name registered for the plain codec.
	Header should set:
	Content-Type:text/plain
	Accept:text/plain
*/
const Name = "plain"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with plain.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	reply,ok:=v.(*wrappers.StringValue)

	if !ok {
		return nil, errors.New("please use wrappers.StringValue")
	}

	return []byte(reply.GetValue()),nil
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case *wrappers.StringValue:
		str:=wrapperspb.String(string(data))

		val := reflect.ValueOf(m)
		val.Elem().Set(reflect.ValueOf(*str))

		return nil
	default:
		return errors.New("please use wrappers.StringValue")
	}
}

func (codec) Name() string {
	return Name
}
