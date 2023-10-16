package http

import (
	"net/http"
	"reflect"
	"testing"
)

func TestEmptyCallOptions(t *testing.T) {
	e := EmptyCallOption{}
	if e.before(&callInfo{}) != nil {
		t.Error("EmptyCallOption should be ignored")
	}
	e.after(&callInfo{}, &csAttempt{})
}

func TestContentType(t *testing.T) {
	if !reflect.DeepEqual(ContentType("aaa").(ContentTypeCallOption).ContentType, "aaa") {
		t.Errorf("want: %v,got: %v", "aaa", ContentType("aaa").(ContentTypeCallOption).ContentType)
	}
}

func TestContentTypeCallOption_before(t *testing.T) {
	c := &callInfo{}
	err := ContentType("aaa").before(c)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual("aaa", c.contentType) {
		t.Errorf("want: %v, got: %v", "aaa", c.contentType)
	}
}

func TestDefaultCallInfo(t *testing.T) {
	path := "hi"
	rv := defaultCallInfo(path)
	if !reflect.DeepEqual(path, rv.pathTemplate) {
		t.Errorf("expect %v, got %v", path, rv.pathTemplate)
	}
	if !reflect.DeepEqual(path, rv.operation) {
		t.Errorf("expect %v, got %v", path, rv.operation)
	}
	if !reflect.DeepEqual("application/json", rv.contentType) {
		t.Errorf("expect %v, got %v", "application/json", rv.contentType)
	}
}

func TestOperation(t *testing.T) {
	if !reflect.DeepEqual("aaa", Operation("aaa").(OperationCallOption).Operation) {
		t.Errorf("want: %v,got: %v", "aaa", Operation("aaa").(OperationCallOption).Operation)
	}
}

func TestOperationCallOption_before(t *testing.T) {
	c := &callInfo{}
	err := Operation("aaa").before(c)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual("aaa", c.operation) {
		t.Errorf("want: %v, got: %v", "aaa", c.operation)
	}
}

func TestPathTemplate(t *testing.T) {
	if !reflect.DeepEqual("aaa", PathTemplate("aaa").(PathTemplateCallOption).Pattern) {
		t.Errorf("want: %v,got: %v", "aaa", PathTemplate("aaa").(PathTemplateCallOption).Pattern)
	}
}

func TestPathTemplateCallOption_before(t *testing.T) {
	c := &callInfo{}
	err := PathTemplate("aaa").before(c)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual("aaa", c.pathTemplate) {
		t.Errorf("want: %v, got: %v", "aaa", c.pathTemplate)
	}
}

func TestHeader(t *testing.T) {
	h := http.Header{"A": []string{"123"}}
	if !reflect.DeepEqual(Header(&h).(HeaderCallOption).header.Get("A"), "123") {
		t.Errorf("want: %v,got: %v", "123", Header(&h).(HeaderCallOption).header.Get("A"))
	}
}

func TestHeaderCallOption_before(t *testing.T) {
	h := http.Header{"A": []string{"123"}}
	c := &callInfo{}
	o := Header(&h)
	err := o.before(c)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(&h, c.headerCarrier) {
		t.Errorf("want: %v,got: %v", &h, o.(HeaderCallOption).header)
	}
}

func TestHeaderCallOption_after(t *testing.T) {
	h := http.Header{"A": []string{"123"}}
	c := &callInfo{}
	cs := &csAttempt{res: &http.Response{Header: h}}
	o := Header(&h)
	o.after(c, cs)
	if !reflect.DeepEqual(&h, o.(HeaderCallOption).header) {
		t.Errorf("want: %v,got: %v", &h, o.(HeaderCallOption).header)
	}
}
