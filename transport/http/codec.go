package http

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/gorilla/mux"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const baseContentType = "application"

func contentSubtype(contentType string) string {
	if contentType == baseContentType {
		return ""
	}
	if !strings.HasPrefix(contentType, baseContentType) {
		return ""
	}
	// guaranteed since != baseContentType and has baseContentType prefix
	switch contentType[len(baseContentType)] {
	case '/', ';':
		// this will return true for "application/grpc+" or "application/grpc;"
		// which the previous validContentType function tested to be valid, so we
		// just say that no content-subtype is specified in this case
		return contentType[len(baseContentType)+1:]
	default:
		return ""
	}
}

func codecForReq(req *http.Request) (Codec, error) {
	contentType := req.Header.Get("content-type")
	cc := encoding.GetCodec(contentSubtype(contentType))
	if cc == nil {
		return nil, ErrUnknownCodec(contentType)
	}
	return &codec{codec: cc, req: req}, nil
}

// Codec defines the interface HTTP uses to encode and decode messages.
type Codec interface {
	Param(key string) Value
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(v interface{}) error
}

type codec struct {
	codec encoding.Codec
	req   *http.Request
}

func (c *codec) Param(key string) Value {
	return Value(mux.Vars(c.req)[key])
}

func (c *codec) Marshal(v interface{}) ([]byte, error) {
	return c.codec.Marshal(v)
}

func (c *codec) Unmarshal(v interface{}) error {
	data, err := ioutil.ReadAll(c.req.Body)
	if err != nil {
		return ErrDataLoss(err.Error())
	}
	defer c.req.Body.Close()
	if err = c.codec.Unmarshal(data, v); err != nil {
		return ErrCodecUnmarshal(err.Error())
	}
	return nil
}

// Value is the route variables for the current request.
type Value string

// String converts.
func (v Value) String() string {
	return string(v)
}

// Bool converts.
func (v Value) Bool(val string) (bool, error) {
	return strconv.ParseBool(string(v))
}

// Int32 converts.
func (v Value) Int32() (int32, error) {
	i, err := strconv.ParseInt(v.String(), 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// Int64 converts.
func (v Value) Int64() (int64, error) {
	return strconv.ParseInt(v.String(), 0, 64)
}

// Uint32 converts.
func (v Value) Uint32() (uint32, error) {
	i, err := strconv.ParseUint(v.String(), 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

// Uint64 converts.
func (v Value) Uint64() (uint64, error) {
	return strconv.ParseUint(v.String(), 0, 64)
}

// Float32 converts.
func (v Value) Float32() (float32, error) {
	f, err := strconv.ParseFloat(v.String(), 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

// Float64 converts.
func (v Value) Float64() (float64, error) {
	return strconv.ParseFloat(v.String(), 64)
}

// Timestamp converts.
func (v Value) Timestamp() (*timestamppb.Timestamp, error) {
	var r timestamppb.Timestamp
	unmarshaler := &protojson.UnmarshalOptions{}
	err := unmarshaler.Unmarshal([]byte(v.String()), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Duration converts.
func (v Value) Duration() (*durationpb.Duration, error) {
	var r durationpb.Duration
	unmarshaler := &protojson.UnmarshalOptions{}
	err := unmarshaler.Unmarshal([]byte(v.String()), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
