package trace

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

type mockReport struct {
	sps []*Span
}

func (m *mockReport) WriteSpan(sp *Span) error {
	m.sps = append(m.sps, sp)
	return nil
}

func (m *mockReport) Close() error {
	return nil
}

func TestDapperPropagation(t *testing.T) {
	t.Run("test HTTP progagation", func(t *testing.T) {
		report := &mockReport{}
		t1 := NewTracer("service1", report, true)
		t2 := NewTracer("service2", report, true)
		sp1 := t1.New("opt_1")
		sp2 := sp1.Fork("", "opt_client")
		header := make(http.Header)
		t1.Inject(sp2, HTTPFormat, header)
		sp3, err := t2.Extract(HTTPFormat, header)
		if err != nil {
			t.Fatal(err)
		}
		sp3.Finish(nil)
		sp2.Finish(nil)
		sp1.Finish(nil)

		assert.Len(t, report.sps, 3)
		assert.Equal(t, report.sps[2].context.ParentID, uint64(0))
		assert.Equal(t, report.sps[0].context.TraceID, report.sps[1].context.TraceID)
		assert.Equal(t, report.sps[2].context.TraceID, report.sps[1].context.TraceID)

		assert.Equal(t, report.sps[1].context.ParentID, report.sps[2].context.SpanID)
		assert.Equal(t, report.sps[0].context.ParentID, report.sps[1].context.SpanID)
	})
	t.Run("test gRPC progagation", func(t *testing.T) {
		report := &mockReport{}
		t1 := NewTracer("service1", report, true)
		t2 := NewTracer("service2", report, true)
		sp1 := t1.New("opt_1")
		sp2 := sp1.Fork("", "opt_client")
		md := make(metadata.MD)
		t1.Inject(sp2, GRPCFormat, md)
		sp3, err := t2.Extract(GRPCFormat, md)
		if err != nil {
			t.Fatal(err)
		}
		sp3.Finish(nil)
		sp2.Finish(nil)
		sp1.Finish(nil)

		assert.Len(t, report.sps, 3)
		assert.Equal(t, report.sps[2].context.ParentID, uint64(0))
		assert.Equal(t, report.sps[0].context.TraceID, report.sps[1].context.TraceID)
		assert.Equal(t, report.sps[2].context.TraceID, report.sps[1].context.TraceID)

		assert.Equal(t, report.sps[1].context.ParentID, report.sps[2].context.SpanID)
		assert.Equal(t, report.sps[0].context.ParentID, report.sps[1].context.SpanID)
	})
	t.Run("test normal", func(t *testing.T) {
		report := &mockReport{}
		t1 := NewTracer("service1", report, true)
		sp1 := t1.New("test123")
		sp1.Finish(nil)
	})
	t.Run("test debug progagation", func(t *testing.T) {
		report := &mockReport{}
		t1 := NewTracer("service1", report, true)
		t2 := NewTracer("service2", report, true)
		sp1 := t1.New("opt_1", EnableDebug())
		sp2 := sp1.Fork("", "opt_client")
		header := make(http.Header)
		t1.Inject(sp2, HTTPFormat, header)
		sp3, err := t2.Extract(HTTPFormat, header)
		if err != nil {
			t.Fatal(err)
		}
		sp3.Finish(nil)
		sp2.Finish(nil)
		sp1.Finish(nil)

		assert.Len(t, report.sps, 3)
		assert.Equal(t, report.sps[2].context.ParentID, uint64(0))
		assert.Equal(t, report.sps[0].context.TraceID, report.sps[1].context.TraceID)
		assert.Equal(t, report.sps[2].context.TraceID, report.sps[1].context.TraceID)

		assert.Equal(t, report.sps[1].context.ParentID, report.sps[2].context.SpanID)
		assert.Equal(t, report.sps[0].context.ParentID, report.sps[1].context.SpanID)
	})
}

func BenchmarkSample(b *testing.B) {
	err := fmt.Errorf("test error")
	report := &mockReport{}
	t1 := NewTracer("service1", report, true)
	for i := 0; i < b.N; i++ {
		sp1 := t1.New("test_opt1")
		sp1.SetTag(TagString("test", "123"))
		sp2 := sp1.Fork("", "opt2")
		sp3 := sp2.Fork("", "opt3")
		sp3.SetTag(TagString("test", "123"))
		sp3.Finish(nil)
		sp2.Finish(&err)
		sp1.Finish(nil)
	}
}

func BenchmarkDisableSample(b *testing.B) {
	err := fmt.Errorf("test error")
	report := &mockReport{}
	t1 := NewTracer("service1", report, true)
	for i := 0; i < b.N; i++ {
		sp1 := t1.New("test_opt1")
		sp1.SetTag(TagString("test", "123"))
		sp2 := sp1.Fork("", "opt2")
		sp3 := sp2.Fork("", "opt3")
		sp3.SetTag(TagString("test", "123"))
		sp3.Finish(nil)
		sp2.Finish(&err)
		sp1.Finish(nil)
	}
}
