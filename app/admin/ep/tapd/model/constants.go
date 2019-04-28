package model

// pagination.
const (
	DefaultPageSize = 5
	DefaultPageNum  = 1
)

// Hook URL Status element
const (
	HookURLStatusEnable  = 1
	HookURLStatusDisable = 2
)

// HookEventStatusEnable
const (
	HookEventStatusEnable  = 1
	HookEventStatusDisable = 2
)

// Event Event
type Event string

// event element
const (
	StoryCreate      Event = "story::create"
	BugCreate              = "bug::create"
	TaskCreate             = "task::create"
	LaunchformCreate       = "launchform::create"

	StoryUpdate      = "story::update"
	BugUpdate        = "bug::update"
	TaskUpdate       = "task::update"
	LaunchformUpdate = "launchform::update"

	StoryDelete = "story::delete"
	BugDelete   = "bug::delete"
	TaskDelete  = "task::delete"
)
