package trace

import (
	"testing"
)

func TestSpanContext(t *testing.T) {
	pctx := &spanContext{
		ParentID: genID(),
		SpanID:   genID(),
		TraceID:  genID(),
		Flags:    flagSampled,
	}
	if !pctx.isSampled() {
		t.Error("expect sampled")
	}
	value := pctx.String()
	t.Logf("bili-trace-id: %s", value)
	pctx2, err := contextFromString(value)
	if err != nil {
		t.Error(err)
	}
	if pctx2.ParentID != pctx.ParentID || pctx2.SpanID != pctx.SpanID || pctx2.TraceID != pctx.TraceID || pctx2.Flags != pctx.Flags {
		t.Errorf("wrong spancontext get %+v -> %+v", pctx, pctx2)
	}
}
