package form

import (
	"encoding/base64"
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/go-kratos/kratos/v2/encoding"
	bdtest "github.com/go-kratos/kratos/v2/internal/testdata/binding"
	"github.com/go-kratos/kratos/v2/internal/testdata/complex"
	ectest "github.com/go-kratos/kratos/v2/internal/testdata/encoding"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type TestModel struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

const contentType = "x-www-form-urlencoded"

func TestFormCodecMarshal(t *testing.T) {
	req := &LoginRequest{
		Username: "kratos",
		Password: "kratos_pwd",
	}
	content, err := encoding.GetCodec(contentType).Marshal(req)
	if err != nil {
		t.Errorf("marshal error: %v", err)
	}
	if !reflect.DeepEqual([]byte("password=kratos_pwd&username=kratos"), content) {
		t.Errorf("expect %v, got %v", []byte("password=kratos_pwd&username=kratos"), content)
	}

	req = &LoginRequest{
		Username: "kratos",
		Password: "",
	}
	content, err = encoding.GetCodec(contentType).Marshal(req)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual([]byte("username=kratos"), content) {
		t.Errorf("expect %v, got %v", []byte("username=kratos"), content)
	}

	m := &TestModel{
		ID:   1,
		Name: "kratos",
	}
	content, err = encoding.GetCodec(contentType).Marshal(m)
	t.Log(string(content))
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual([]byte("id=1&name=kratos"), content) {
		t.Errorf("expect %v, got %v", []byte("id=1&name=kratos"), content)
	}
}

func TestFormCodecUnmarshal(t *testing.T) {
	req := &LoginRequest{
		Username: "kratos",
		Password: "kratos_pwd",
	}
	content, err := encoding.GetCodec(contentType).Marshal(req)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}

	bindReq := new(LoginRequest)
	err = encoding.GetCodec(contentType).Unmarshal(content, bindReq)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual("kratos", bindReq.Username) {
		t.Errorf("expect %v, got %v", "kratos", bindReq.Username)
	}
	if !reflect.DeepEqual("kratos_pwd", bindReq.Password) {
		t.Errorf("expect %v, got %v", "kratos_pwd", bindReq.Password)
	}
}

func TestProtoEncodeDecode(t *testing.T) {
	in := &complex.Complex{
		Id:      2233,
		NoOne:   "2233",
		Simple:  &complex.Simple{Component: "5566"},
		Simples: []string{"3344", "5566"},
		B:       true,
		Sex:     complex.Sex_woman,
		Age:     18,
		A:       19,
		Count:   3,
		Price:   11.23,
		D:       22.22,
		Byte:    []byte("123"),
		Map:     map[string]string{"kratos": "https://go-kratos.dev/"},

		Timestamp: &timestamppb.Timestamp{Seconds: 20, Nanos: 2},
		Duration:  &durationpb.Duration{Seconds: 120, Nanos: 22},
		Field:     &fieldmaskpb.FieldMask{Paths: []string{"1", "2"}},
		Double:    &wrapperspb.DoubleValue{Value: 12.33},
		Float:     &wrapperspb.FloatValue{Value: 12.34},
		Int64:     &wrapperspb.Int64Value{Value: 64},
		Int32:     &wrapperspb.Int32Value{Value: 32},
		Uint64:    &wrapperspb.UInt64Value{Value: 64},
		Uint32:    &wrapperspb.UInt32Value{Value: 32},
		Bool:      &wrapperspb.BoolValue{Value: false},
		String_:   &wrapperspb.StringValue{Value: "go-kratos"},
		Bytes:     &wrapperspb.BytesValue{Value: []byte("123")},
	}
	content, err := encoding.GetCodec(contentType).Marshal(in)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual("a=19&age=18&b=true&bool=false&byte=MTIz&bytes=MTIz&count=3&d=22.22&double=12.33&duration="+
		"2m0.000000022s&field=1%2C2&float=12.34&id=2233&int32=32&int64=64&map%5Bkratos%5D=https%3A%2F%2Fgo-kratos.dev%2F&"+
		"numberOne=2233&price=11.23&sex=woman&simples=3344&simples=5566&string=go-kratos"+
		"&timestamp=1970-01-01T00%3A00%3A20.000000002Z&uint32=32&uint64=64&very_simple.component=5566", string(content)) {
		t.Errorf("rawpath is not equal to %v", string(content))
	}
	in2 := &complex.Complex{}
	err = encoding.GetCodec(contentType).Unmarshal(content, in2)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(int64(2233), in2.Id) {
		t.Errorf("expect %v, got %v", int64(2233), in2.Id)
	}
	if !reflect.DeepEqual("2233", in2.NoOne) {
		t.Errorf("expect %v, got %v", "2233", in2.NoOne)
	}
	if reflect.DeepEqual(in2.Simple, nil) {
		t.Errorf("expect %v, got %v", nil, in2.Simple)
	}
	if !reflect.DeepEqual("5566", in2.Simple.Component) {
		t.Errorf("expect %v, got %v", "5566", in2.Simple.Component)
	}
	if reflect.DeepEqual(in2.Simples, nil) {
		t.Errorf("expect %v, got %v", nil, in2.Simples)
	}
	if !reflect.DeepEqual(len(in2.Simples), 2) {
		t.Errorf("expect %v, got %v", 2, len(in2.Simples))
	}
	if !reflect.DeepEqual("3344", in2.Simples[0]) {
		t.Errorf("expect %v, got %v", "3344", in2.Simples[0])
	}
	if !reflect.DeepEqual("5566", in2.Simples[1]) {
		t.Errorf("expect %v, got %v", "5566", in2.Simples[1])
	}
}

func TestDecodeStructPb(t *testing.T) {
	req := new(ectest.StructPb)
	query := `data={"name":"kratos"}&data_list={"name1": "kratos"}&data_list={"name2": "go-kratos"}`
	if err := encoding.GetCodec(contentType).Unmarshal([]byte(query), req); err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual("kratos", req.Data.GetFields()["name"].GetStringValue()) {
		t.Errorf("except %v, got %v", "kratos", req.Data.GetFields()["name"].GetStringValue())
	}
	if len(req.DataList) != 2 {
		t.Errorf("execpt %v, got %v", 2, len(req.DataList))
		return
	}
	if !reflect.DeepEqual("kratos", req.DataList[0].GetFields()["name1"].GetStringValue()) {
		t.Errorf("except %v, got %v", "kratos", req.Data.GetFields()["name1"].GetStringValue())
	}
	if !reflect.DeepEqual("go-kratos", req.DataList[1].GetFields()["name2"].GetStringValue()) {
		t.Errorf("except %v, got %v", "go-kratos", req.Data.GetFields()["name2"].GetStringValue())
	}
}

func TestDecodeBytesValuePb(t *testing.T) {
	url := "https://example.com/xx/?a=1&b=2&c=3"
	val := base64.URLEncoding.EncodeToString([]byte(url))
	content := "bytes=" + val
	in2 := &complex.Complex{}
	if err := encoding.GetCodec(contentType).Unmarshal([]byte(content), in2); err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(url, string(in2.Bytes.Value)) {
		t.Errorf("except %v, got %v", val, string(in2.Bytes.Value))
	}
}

func TestEncodeFieldMask(t *testing.T) {
	req := &bdtest.HelloRequest{
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"foo", "bar"}},
	}
	if v := EncodeFieldMask(req.ProtoReflect()); v != "updateMask=foo,bar" {
		t.Errorf("got %s", v)
	} else {
		t.Log(v)
	}
}

func TestOptional(t *testing.T) {
	v := int32(100)
	req := &bdtest.HelloRequest{
		Name:     "foo",
		Sub:      &bdtest.Sub{Name: "bar"},
		OptInt32: &v,
	}
	if v, _ := EncodeValues(req); v.Encode() != "name=foo&optInt32=100&sub.naming=bar" {
		t.Errorf("got %s", v.Encode())
	} else {
		t.Log(v)
	}
}
