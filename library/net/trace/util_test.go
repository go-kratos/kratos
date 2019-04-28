package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromContext(t *testing.T) {
	report := &mockReport{}
	t1 := newTracer("service1", report, &Config{DisableSample: true})
	sp1 := t1.New("test123")
	ctx := context.Background()
	ctx = NewContext(ctx, sp1)
	sp2, ok := FromContext(ctx)
	if !ok {
		t.Fatal("nothing from context")
	}
	assert.Equal(t, sp1, sp2)
}
