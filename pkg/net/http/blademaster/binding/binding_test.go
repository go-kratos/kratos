package binding

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" validate:"required"`
}

type FooBarStruct struct {
	FooStruct
	Bar   string   `msgpack:"bar" json:"bar" form:"bar" xml:"bar" validate:"required"`
	Slice []string `form:"slice" validate:"max=10"`
}

type ComplexDefaultStruct struct {
	Int        int     `form:"int" default:"999"`
	String     string  `form:"string" default:"default-string"`
	Bool       bool    `form:"bool" default:"false"`
	Int64Slice []int64 `form:"int64_slice,split" default:"1,2,3,4"`
	Int8Slice  []int8  `form:"int8_slice,split" default:"1,2,3,4"`
}

type Int8SliceStruct struct {
	State []int8 `form:"state,split"`
}

type Int64SliceStruct struct {
	State []int64 `form:"state,split"`
}

type StringSliceStruct struct {
	State []string `form:"state,split"`
}

func TestBindingDefault(t *testing.T) {
	assert.Equal(t, Default("GET", ""), Form)
	assert.Equal(t, Default("GET", MIMEJSON), Form)
	assert.Equal(t, Default("GET", MIMEJSON+"; charset=utf-8"), Form)

	assert.Equal(t, Default("POST", MIMEJSON), JSON)
	assert.Equal(t, Default("PUT", MIMEJSON), JSON)

	assert.Equal(t, Default("POST", MIMEJSON+"; charset=utf-8"), JSON)
	assert.Equal(t, Default("PUT", MIMEJSON+"; charset=utf-8"), JSON)

	assert.Equal(t, Default("POST", MIMEXML), XML)
	assert.Equal(t, Default("PUT", MIMEXML2), XML)

	assert.Equal(t, Default("POST", MIMEPOSTForm), Form)
	assert.Equal(t, Default("PUT", MIMEPOSTForm), Form)

	assert.Equal(t, Default("POST", MIMEPOSTForm+"; charset=utf-8"), Form)
	assert.Equal(t, Default("PUT", MIMEPOSTForm+"; charset=utf-8"), Form)

	assert.Equal(t, Default("POST", MIMEMultipartPOSTForm), Form)
	assert.Equal(t, Default("PUT", MIMEMultipartPOSTForm), Form)

}

func TestStripContentType(t *testing.T) {
	c1 := "application/vnd.mozilla.xul+xml"
	c2 := "application/vnd.mozilla.xul+xml; charset=utf-8"
	assert.Equal(t, stripContentTypeParam(c1), c1)
	assert.Equal(t, stripContentTypeParam(c2), "application/vnd.mozilla.xul+xml")
}

func TestBindInt8Form(t *testing.T) {
	params := "state=1,2,3"
	req, _ := http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q := new(Int8SliceStruct)
	Form.Bind(req, q)
	assert.EqualValues(t, []int8{1, 2, 3}, q.State)

	params = "state=1,2,3,256"
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(Int8SliceStruct)
	assert.Error(t, Form.Bind(req, q))

	params = "state="
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(Int8SliceStruct)
	assert.NoError(t, Form.Bind(req, q))
	assert.Len(t, q.State, 0)

	params = "state=1,,2"
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(Int8SliceStruct)
	assert.NoError(t, Form.Bind(req, q))
	assert.EqualValues(t, []int8{1, 2}, q.State)
}

func TestBindInt64Form(t *testing.T) {
	params := "state=1,2,3"
	req, _ := http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q := new(Int64SliceStruct)
	Form.Bind(req, q)
	assert.EqualValues(t, []int64{1, 2, 3}, q.State)

	params = "state="
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(Int64SliceStruct)
	assert.NoError(t, Form.Bind(req, q))
	assert.Len(t, q.State, 0)
}

func TestBindStringForm(t *testing.T) {
	params := "state=1,2,3"
	req, _ := http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q := new(StringSliceStruct)
	Form.Bind(req, q)
	assert.EqualValues(t, []string{"1", "2", "3"}, q.State)

	params = "state="
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(StringSliceStruct)
	assert.NoError(t, Form.Bind(req, q))
	assert.Len(t, q.State, 0)

	params = "state=p,,p"
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(StringSliceStruct)
	Form.Bind(req, q)
	assert.EqualValues(t, []string{"p", "p"}, q.State)
}

func TestBindingJSON(t *testing.T) {
	testBodyBinding(t,
		JSON, "json",
		"/", "/",
		`{"foo": "bar"}`, `{"bar": "foo"}`)
}

func TestBindingForm(t *testing.T) {
	testFormBinding(t, "POST",
		"/", "/",
		"foo=bar&bar=foo&slice=a&slice=b", "bar2=foo")
}

func TestBindingForm2(t *testing.T) {
	testFormBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"", "")
}

func TestBindingQuery(t *testing.T) {
	testQueryBinding(t, "POST",
		"/?foo=bar&bar=foo", "/",
		"foo=unused", "bar2=foo")
}

func TestBindingQuery2(t *testing.T) {
	testQueryBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"foo=unused", "")
}

func TestBindingXML(t *testing.T) {
	testBodyBinding(t,
		XML, "xml",
		"/", "/",
		"<map><foo>bar</foo></map>", "<map><bar>foo</bar></map>")
}

func createFormPostRequest() *http.Request {
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", bytes.NewBufferString("foo=bar&bar=foo"))
	req.Header.Set("Content-Type", MIMEPOSTForm)
	return req
}

func createFormMultipartRequest() *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	mw.SetBoundary(boundary)
	mw.WriteField("foo", "bar")
	mw.WriteField("bar", "foo")
	req, _ := http.NewRequest("POST", "/?foo=getfoo&bar=getbar", body)
	req.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	return req
}

func TestBindingFormPost(t *testing.T) {
	req := createFormPostRequest()
	var obj FooBarStruct
	FormPost.Bind(req, &obj)

	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
}

func TestBindingFormMultipart(t *testing.T) {
	req := createFormMultipartRequest()
	var obj FooBarStruct
	FormMultipart.Bind(req, &obj)

	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
}

func TestValidationFails(t *testing.T) {
	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func TestValidationDisabled(t *testing.T) {
	backup := Validator
	Validator = nil
	defer func() { Validator = backup }()

	var obj FooStruct
	req := requestWithBody("POST", "/", `{"bar": "foo"}`)
	err := JSON.Bind(req, &obj)
	assert.NoError(t, err)
}

func TestExistsSucceeds(t *testing.T) {
	type HogeStruct struct {
		Hoge *int `json:"hoge" binding:"exists"`
	}

	var obj HogeStruct
	req := requestWithBody("POST", "/", `{"hoge": 0}`)
	err := JSON.Bind(req, &obj)
	assert.NoError(t, err)
}

func TestFormDefaultValue(t *testing.T) {
	params := "int=333&string=hello&bool=true&int64_slice=5,6,7,8&int8_slice=5,6,7,8"
	req, _ := http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q := new(ComplexDefaultStruct)
	assert.NoError(t, Form.Bind(req, q))
	assert.Equal(t, 333, q.Int)
	assert.Equal(t, "hello", q.String)
	assert.Equal(t, true, q.Bool)
	assert.EqualValues(t, []int64{5, 6, 7, 8}, q.Int64Slice)
	assert.EqualValues(t, []int8{5, 6, 7, 8}, q.Int8Slice)

	params = "string=hello&bool=false"
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(ComplexDefaultStruct)
	assert.NoError(t, Form.Bind(req, q))
	assert.Equal(t, 999, q.Int)
	assert.Equal(t, "hello", q.String)
	assert.Equal(t, false, q.Bool)
	assert.EqualValues(t, []int64{1, 2, 3, 4}, q.Int64Slice)
	assert.EqualValues(t, []int8{1, 2, 3, 4}, q.Int8Slice)

	params = "strings=hello"
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(ComplexDefaultStruct)
	assert.NoError(t, Form.Bind(req, q))
	assert.Equal(t, 999, q.Int)
	assert.Equal(t, "default-string", q.String)
	assert.Equal(t, false, q.Bool)
	assert.EqualValues(t, []int64{1, 2, 3, 4}, q.Int64Slice)
	assert.EqualValues(t, []int8{1, 2, 3, 4}, q.Int8Slice)

	params = "int=&string=&bool=true&int64_slice=&int8_slice="
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	q = new(ComplexDefaultStruct)
	assert.NoError(t, Form.Bind(req, q))
	assert.Equal(t, 999, q.Int)
	assert.Equal(t, "default-string", q.String)
	assert.Equal(t, true, q.Bool)
	assert.EqualValues(t, []int64{1, 2, 3, 4}, q.Int64Slice)
	assert.EqualValues(t, []int8{1, 2, 3, 4}, q.Int8Slice)
}

func testFormBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Form
	assert.Equal(t, b.Name(), "form")

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")

	obj = FooBarStruct{}
	req = requestWithBody(method, badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testQueryBinding(t *testing.T, method, path, badPath, body, badBody string) {
	b := Query
	assert.Equal(t, b.Name(), "query")

	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")
	assert.Equal(t, obj.Bar, "foo")
}

func testBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, b.Name(), name)

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, obj.Foo, "bar")

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
func BenchmarkBindingForm(b *testing.B) {
	req := requestWithBody("POST", "/", "foo=bar&bar=foo&slice=a&slice=b&slice=c&slice=w")
	req.Header.Add("Content-Type", MIMEPOSTForm)
	f := Form
	for i := 0; i < b.N; i++ {
		obj := FooBarStruct{}
		f.Bind(req, &obj)
	}
}
