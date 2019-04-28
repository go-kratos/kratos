package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelationStateMachine(t *testing.T) {
	rs := NewRelationStateMachine(StateNoRelation, DefaultHandler)
	assert.NotNil(t, rs)
	assert.Equal(t, RelationState(rs.Current()), StateNoRelation)

	assert.NoError(t, rs.Event(EventAddFollowing))
	assert.Equal(t, StateFollowing, RelationState(rs.Current()))
}
