package trace

import (
	"testing"
)

func TestSpanContext(t *testing.T) {
	pctx := &spanContext{
		parentID: genID(),
		spanID:   genID(),
		traceID:  genID(),
		flags:    flagSampled,
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
	if pctx2.parentID != pctx.parentID || pctx2.spanID != pctx.spanID || pctx2.traceID != pctx.traceID || pctx2.flags != pctx.flags {
		t.Errorf("wrong spancontext get %+v -> %+v", pctx, pctx2)
	}
}
