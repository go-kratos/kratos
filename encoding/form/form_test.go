package form

import (
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/go-kratos/kratos/v2/encoding"
	bdtest "github.com/go-kratos/kratos/v2/internal/testdata/binding"
	"github.com/go-kratos/kratos/v2/internal/testdata/complex"
	ectest "github.com/go-kratos/kratos/v2/internal/testdata/encoding"
)

// This variable can be replaced with -ldflags like below:
// go test "-ldflags=-X github.com/go-kratos/kratos/v2/encoding/form.tagNameTest=form"
var tagNameTest string

func init() {
	if tagNameTest == "" {
		tagNameTest = tagName
	}
}

func TestFormEncoderAndDecoder(t *testing.T) {
	t.Cleanup(func() {
		encoder.SetTagName(tagName)
		decoder.SetTagName(tagName)
	})

	encoder.SetTagName(tagNameTest)
	decoder.SetTagName(tagNameTest)

	type testFormTagName struct {
		Name string `form:"name_form" json:"name_json"`
	}
	v, err := encoder.Encode(&testFormTagName{
		Name: "test tag name",
	})
	if err != nil {
		t.Fatal(err)
	}
	jsonName := v.Get("name_json")
	formName := v.Get("name_form")
	switch tagNameTest {
	case "json":
		if jsonName != "test tag name" {
			t.Errorf("got: %s", jsonName)
		}
		if formName != "" {
			t.Errorf("want: empty, got: %s", formName)
		}
	case "form":
		if formName != "test tag name" {
			t.Errorf("got: %s", formName)
		}
		if jsonName != "" {
			t.Errorf("want: empty, got: %s", jsonName)
		}
	default:
		t.Fatalf("unknown tag name: %s", tagNameTest)
	}

	var tn *testFormTagName
	err = decoder.Decode(&tn, v)
	if err != nil {
		t.Fatal(err)
	}
	if tn == nil {
		t.Fatal("nil tag name")
	}
	if tn.Name != "test tag name" {
		t.Errorf("got %s", tn.Name)
	}
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func TestFormCodecMarshal(t *testing.T) {
	req := &LoginRequest{
		Username: "kratos",
		Password: "kratos_pwd",
	}
	content, err := encoding.GetCodec(Name).Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]byte("password=kratos_pwd&username=kratos"), content) {
		t.Errorf("expect %s, got %s", "password=kratos_pwd&username=kratos", content)
	}

	req = &LoginRequest{
		Username: "kratos",
		Password: "",
	}
	content, err = encoding.GetCodec(Name).Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]byte("username=kratos"), content) {
		t.Errorf("expect %s, got %s", "username=kratos", content)
	}

	m := struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	}{
		ID:   1,
		Name: "kratos",
	}
	content, err = encoding.GetCodec(Name).Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]byte("id=1&name=kratos"), content) {
		t.Errorf("expect %s, got %s", "id=1&name=kratos", content)
	}
}

func TestFormCodecUnmarshal(t *testing.T) {
	req := &LoginRequest{
		Username: "kratos",
		Password: "kratos_pwd",
	}
	content, err := encoding.GetCodec(Name).Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	bindReq := new(LoginRequest)
	err = encoding.GetCodec(Name).Unmarshal(content, bindReq)
	if err != nil {
		t.Fatal(err)
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
		Map:     map[string]string{"kratos": "https://go-kratos.dev/", "kratos_start": "https://go-kratos.dev/en/docs/getting-started/start/"},

		Timestamp: timestamppb.New(time.Date(1970, 1, 1, 0, 0, 20, 2, time.Local)),
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
	content, err := encoding.GetCodec(Name).Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	if "a=19&age=18&b=true&bool=false&byte=MTIz&bytes=MTIz&count=3&d=22.22&double=12.33&duration="+
		"2m0.000000022s&field=1%2C2&float=12.34&id=2233&int32=32&int64=64&"+
		"map%5Bkratos%5D=https%3A%2F%2Fgo-kratos.dev%2F&map%5Bkratos_start%5D=https%3A%2F%2Fgo-kratos.dev%2Fen%2Fdocs%2Fgetting-started%2Fstart%2F&"+
		"numberOne=2233&price=11.23&sex=woman&simples=3344&simples=5566&string=go-kratos"+
		"&timestamp=1970-01-01T00%3A00%3A20.000000002Z&uint32=32&uint64=64&very_simple.component=5566" != string(content) {
		t.Errorf("rawpath is not equal to %s", content)
	}
	in2 := &complex.Complex{}
	err = encoding.GetCodec(Name).Unmarshal(content, in2)
	if err != nil {
		t.Fatal(err)
	}
	if int64(2233) != in2.Id {
		t.Errorf("expect %v, got %v", int64(2233), in2.Id)
	}
	if "2233" != in2.NoOne {
		t.Errorf("expect %v, got %v", "2233", in2.NoOne)
	}
	if in2.Simple == nil {
		t.Errorf("expect %v, got %v", nil, in2.Simple)
	}
	if "5566" != in2.Simple.Component {
		t.Errorf("expect %v, got %v", "5566", in2.Simple.Component)
	}
	if in2.Simples == nil {
		t.Errorf("expect %v, got %v", nil, in2.Simples)
	}
	if len(in2.Simples) != 2 {
		t.Errorf("expect %v, got %v", 2, len(in2.Simples))
	}
	if "3344" != in2.Simples[0] {
		t.Errorf("expect %v, got %v", "3344", in2.Simples[0])
	}
	if "5566" != in2.Simples[1] {
		t.Errorf("expect %v, got %v", "5566", in2.Simples[1])
	}
	if l := len(in2.GetMap()); l != 2 {
		t.Fatalf("in2.Map length want: %d, got: %d", 2, l)
	}
	for key, val := range in.GetMap() {
		if in2Val := in2.GetMap()[key]; in2Val != val {
			t.Errorf("%s want: %q, got: %q", "map["+key+"]", val, in2Val)
		}
	}
}

func TestDecodeStructPb(t *testing.T) {
	req := new(ectest.StructPb)
	query := `data={"name":"kratos"}&data_list={"name1": "kratos"}&data_list={"name2": "go-kratos"}`
	if err := encoding.GetCodec(Name).Unmarshal([]byte(query), req); err != nil {
		t.Fatal(err)
	}
	if "kratos" != req.Data.GetFields()["name"].GetStringValue() {
		t.Errorf("except %v, got %v", "kratos", req.Data.GetFields()["name"].GetStringValue())
	}
	if len(req.DataList) != 2 {
		t.Fatalf("except %v, got %v", 2, len(req.DataList))
	}
	if "kratos" != req.DataList[0].GetFields()["name1"].GetStringValue() {
		t.Errorf("except %v, got %v", "kratos", req.Data.GetFields()["name1"].GetStringValue())
	}
	if "go-kratos" != req.DataList[1].GetFields()["name2"].GetStringValue() {
		t.Errorf("except %v, got %v", "go-kratos", req.Data.GetFields()["name2"].GetStringValue())
	}
}

func TestDecodeBytesValuePb(t *testing.T) {
	url := "https://example.com/xx/?a=1&b=2&c=3"
	val := base64.URLEncoding.EncodeToString([]byte(url))
	content := "bytes=" + val
	in2 := &complex.Complex{}
	if err := encoding.GetCodec(Name).Unmarshal([]byte(content), in2); err != nil {
		t.Fatal(err)
	}
	if url != string(in2.Bytes.Value) {
		t.Errorf("except %s, got %s", val, in2.Bytes.Value)
	}
}

func TestEncodeFieldMask(t *testing.T) {
	req := &bdtest.HelloRequest{
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"foo", "bar"}},
	}
	if v := EncodeFieldMask(req.ProtoReflect()); v != "updateMask=foo,bar" {
		t.Errorf("got %s", v)
	}
}

func TestOptional(t *testing.T) {
	v := int32(100)
	req := &bdtest.HelloRequest{
		Name:     "foo",
		Sub:      &bdtest.Sub{Name: "bar"},
		OptInt32: &v,
	}
	query, _ := EncodeValues(req)
	if query.Encode() != "name=foo&optInt32=100&sub.naming=bar" {
		t.Fatalf("got %s", query.Encode())
	}
}
