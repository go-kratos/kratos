package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyCallOptions(t *testing.T) {
	assert.NoError(t, EmptyCallOption{}.before(&callInfo{}))
	EmptyCallOption{}.after(&callInfo{}, &csAttempt{})
}

func TestContentType(t *testing.T) {
	assert.Equal(t, "aaa", ContentType("aaa").(ContentTypeCallOption).ContentType)
}

func TestContentTypeCallOption_before(t *testing.T) {
	c := &callInfo{}
	err := ContentType("aaa").before(c)
	assert.NoError(t, err)
	assert.Equal(t, "aaa", c.contentType)
}

func TestDefaultCallInfo(t *testing.T) {
	path := "hi"
	rv := defaultCallInfo(path)
	assert.Equal(t, path, rv.pathTemplate)
	assert.Equal(t, path, rv.operation)
	assert.Equal(t, "application/json", rv.contentType)
}

func TestOperation(t *testing.T) {
	assert.Equal(t, "aaa", Operation("aaa").(OperationCallOption).Operation)
}

func TestOperationCallOption_before(t *testing.T) {
	c := &callInfo{}
	err := Operation("aaa").before(c)
	assert.NoError(t, err)
	assert.Equal(t, "aaa", c.operation)
}

func TestPathTemplate(t *testing.T) {
	assert.Equal(t, "aaa", PathTemplate("aaa").(PathTemplateCallOption).Pattern)
}

func TestPathTemplateCallOption_before(t *testing.T) {
	c := &callInfo{}
	err := PathTemplate("aaa").before(c)
	assert.NoError(t, err)
	assert.Equal(t, "aaa", c.pathTemplate)
}

func TestHeader(t *testing.T) {
	h := http.Header{"A": []string{"123"}}
	assert.Equal(t, "123", Header(&h).(HeaderCallOption).header.Get("A"))
}

func TestHeaderCallOption_after(t *testing.T) {
	h := http.Header{"A": []string{"123"}}
	c := &callInfo{}
	cs := &csAttempt{res: &http.Response{Header: h}}
	o := Header(&h)
	o.after(c, cs)
	assert.Equal(t, &h, o.(HeaderCallOption).header)
}
