package trace

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

type mockReport struct {
	sps []*span
}

func (m *mockReport) WriteSpan(sp *span) error {
	m.sps = append(m.sps, sp)
	return nil
}

func (m *mockReport) Close() error {
	return nil
}

func TestDapperPropagation(t *testing.T) {
	t.Run("test HTTP progagation", func(t *testing.T) {
		report := &mockReport{}
		t1 := newTracer("service1", report, &Config{DisableSample: true})
		t2 := newTracer("service2", report, &Config{DisableSample: true})
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
		assert.Equal(t, report.sps[2].context.parentID, uint64(0))
		assert.Equal(t, report.sps[0].context.traceID, report.sps[1].context.traceID)
		assert.Equal(t, report.sps[2].context.traceID, report.sps[1].context.traceID)

		assert.Equal(t, report.sps[1].context.parentID, report.sps[2].context.spanID)
		assert.Equal(t, report.sps[0].context.parentID, report.sps[1].context.spanID)
	})
	t.Run("test gRPC progagation", func(t *testing.T) {
		report := &mockReport{}
		t1 := newTracer("service1", report, &Config{DisableSample: true})
		t2 := newTracer("service2", report, &Config{DisableSample: true})
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
		assert.Equal(t, report.sps[2].context.parentID, uint64(0))
		assert.Equal(t, report.sps[0].context.traceID, report.sps[1].context.traceID)
		assert.Equal(t, report.sps[2].context.traceID, report.sps[1].context.traceID)

		assert.Equal(t, report.sps[1].context.parentID, report.sps[2].context.spanID)
		assert.Equal(t, report.sps[0].context.parentID, report.sps[1].context.spanID)
	})
	t.Run("test normal", func(t *testing.T) {
		report := &mockReport{}
		t1 := newTracer("service1", report, &Config{Probability: 0.000000001})
		sp1 := t1.New("test123")
		sp1.Finish(nil)
	})
	t.Run("test debug progagation", func(t *testing.T) {
		report := &mockReport{}
		t1 := newTracer("service1", report, &Config{})
		t2 := newTracer("service2", report, &Config{})
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
		assert.Equal(t, report.sps[2].context.parentID, uint64(0))
		assert.Equal(t, report.sps[0].context.traceID, report.sps[1].context.traceID)
		assert.Equal(t, report.sps[2].context.traceID, report.sps[1].context.traceID)

		assert.Equal(t, report.sps[1].context.parentID, report.sps[2].context.spanID)
		assert.Equal(t, report.sps[0].context.parentID, report.sps[1].context.spanID)
	})
}

func BenchmarkSample(b *testing.B) {
	err := fmt.Errorf("test error")
	report := &mockReport{}
	t1 := newTracer("service1", report, &Config{})
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
	t1 := newTracer("service1", report, &Config{DisableSample: true})
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
