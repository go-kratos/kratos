---
name: "\U0001F9F1 Feature Request"
about: Description of new features.
title: "[Feature]"
labels: feature
assignees: ''
---

Please see the FAQ in our main README.md before submitting your issue.

<!--
In order to accurately distinguish whether the needs put forward by users are the needs or reasonable needs of most users, solicit opinions from the community through the proposal process, and the proposals adopted by the community will be realized as new functions. 
In order to make the proposal process as simple as possible, the process includes three stages: proposal feature and PR, in which proposal feature is issue and PR is the specific function implementation. 
In order to facilitate the community to correctly understand the requirements of the proposal, the proposal issue needs to describe the functional requirements in detail and relevant references or literature.
The proposal can include the approximate implementation mode of the function, such as interface definition, which can be used as a reference for the function implementation in the feature issue.
When most community users agree with the proposal, A feature issue will be created to associate the proposal issue. The feature issue needs to describe in detail the implementation method and function demonstration of the function as a reference for the final function implementation.
After the function is implemented, a merge request will be initiated to associate the proposal issue and feature issue. After the merge is completed, all issues will be closed.
-->

### Feature description
<!--
Add event interface for accessing message oriented middleware
-->
### Implementation mode
<!--
```go
type Message interface {
    Key() string
    Value() []byte
    Header() map[string]string
    Ack() error
    Nack() error
}

type Handler func(context.Context, Message) error

type Event interface {
    Send(ctx context.Context, key string, value []byte]) error
    Receive(ctx context.Context, handler Handler) error
    Close() error
}
````
-->
### Usage demonstration
<!-- 
```go
msg := kafka.NewMessage("kratos", []byte("hello world"), map[string]string{
		"user":  "kratos",
		"phone": "123456",
	})
err := sender.Send(context.Background(), msg)
```
-->