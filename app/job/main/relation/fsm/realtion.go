package fsm

// RelationState is relation state type
type RelationState string

// RelationEvent is relation event type
type RelationEvent string

// consts for relation
var (
	// states
	StateNoRelation = RelationState("no_relation")
	StateWhisper    = RelationState("whisper")
	StateFollowing  = RelationState("following")
	StateBlacked    = RelationState("blacked")
	StateFriend     = RelationState("friend") // StateFriend is the most special state
	// events
	EventAddFollowing = RelationEvent("add_following")
	EventDelFollowing = RelationEvent("del_following")
	EventAddWhisper   = RelationEvent("add_whisper")
	EventDelWhisper   = RelationEvent("del_whisper")
	EventAddBlack     = RelationEvent("add_black")
	EventDelBlack     = RelationEvent("del_black")
	EventDelFollower  = RelationEvent("del_follower")
	EventBeFriend     = RelationEvent("be_friend") // EventBeFriend is the most special event
)

// RelationEventHandler is used to handle all state change for relation
type RelationEventHandler interface {
	AddFollowing(*Event)
	DelFollowing(*Event)
	AddWhisper(*Event)
	DelWhisper(*Event)
	AddBlack(*Event)
	DelBlack(*Event)
	DelFollower(*Event)
}

// DefaultHandler is the default RelationEventHandler
var DefaultHandler = &defaultHandlerImpl{}

type defaultHandlerImpl struct{}

func (*defaultHandlerImpl) AddFollowing(*Event) {}
func (*defaultHandlerImpl) DelFollowing(*Event) {}
func (*defaultHandlerImpl) AddWhisper(*Event)   {}
func (*defaultHandlerImpl) DelWhisper(*Event)   {}
func (*defaultHandlerImpl) AddBlack(*Event)     {}
func (*defaultHandlerImpl) DelBlack(*Event)     {}
func (*defaultHandlerImpl) DelFollower(*Event)  {}

// RelationStateMachine is used to describe all state change for relation
type RelationStateMachine struct {
	*FSM
}

// NewRelationStateMachine will create a RelationStateMachine
func NewRelationStateMachine(initial RelationState, handler RelationEventHandler) *RelationStateMachine {
	rs := &RelationStateMachine{
		FSM: NewFSM(
			string(StateNoRelation),
			Events{
				{
					Name: string(EventAddFollowing),
					Src: []string{
						string(StateNoRelation),
						string(StateWhisper),
						string(StateBlacked),
					},
					Dst: string(StateFollowing),
				},
				{
					Name: string(EventDelFollowing),
					Src: []string{
						string(StateFollowing),
						string(StateFriend),
					},
					Dst: string(StateNoRelation),
				},
				{
					Name: string(EventAddWhisper),
					Src: []string{
						string(StateNoRelation),
						string(StateFollowing),
						string(StateBlacked),
						string(StateFriend),
					},
					Dst: string(StateWhisper),
				},
				{
					Name: string(EventDelWhisper),
					Src: []string{
						string(StateWhisper),
					},
					Dst: string(StateNoRelation),
				},
				{
					Name: string(EventAddBlack),
					Src: []string{
						string(StateNoRelation),
						string(StateFollowing),
						string(StateFriend),
						string(StateWhisper),
					},
					Dst: string(StateBlacked),
				},
				{
					Name: string(EventDelBlack),
					Src: []string{
						string(StateBlacked),
					},
					Dst: string(StateNoRelation),
				},
				{
					Name: string(EventDelBlack),
					Src: []string{
						string(StateBlacked),
					},
					Dst: string(StateNoRelation),
				},
			},
			Callbacks{
				string(EventAddFollowing): handler.AddFollowing,
				string(EventDelFollowing): handler.DelFollowing,
				string(EventAddWhisper):   handler.AddWhisper,
				string(EventDelWhisper):   handler.DelWhisper,
				string(EventAddBlack):     handler.AddBlack,
				string(EventDelBlack):     handler.DelBlack,
				string(EventDelFollower):  handler.DelFollower,
			},
		),
	}
	return rs
}

// Event is used to execute any events
func (r *RelationStateMachine) Event(event RelationEvent, args ...interface{}) error {
	return r.FSM.Event(string(event), args...)
}

// SetState is used to set state
func (r *RelationStateMachine) SetState(state RelationState) {
	r.FSM.SetState(string(state))
}
