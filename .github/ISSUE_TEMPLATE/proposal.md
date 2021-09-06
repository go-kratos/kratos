---
name: "\U0001F9F1 Proposal Request"
about: Implementation draft of feature.
title: "[Proposal]"
labels: proposal
assignees: ''
---

Please see the FAQ in our main README.md before submitting your issue.

<!--
In order to accurately distinguish that the needs put forward by users are the needs of most users and reasonable needs, solicit community opinions through the process, and the features adopted by the community will be realized as new functions.

In order to make the proposal process as simple as possible, the process includes three stages: feature request - > proposal - > pull-request, where feature, proposal is issue and pull-request is the specific function implementation.

### Feature-request

In order to help the community correctly understand the requirements of the feature, the feature request issue needs to describe the functional requirements and relevant references or documents in detail. And the feature request issue can contain the basic description of the function, which can be used as a reference for the function implementation in the proposal.

### Proposal

Proposal contains the basic implementation methods of functions, such as interface definition, general usage of functions, etc.

### Pull-request

After the function is realized, a merge request will be initiated to associate the proposal issue with the function issue. After the merger is completed, all questions will be closed and the process will end.

### Decision process

When more than five maintainer members agree to implement the feature, a proposal issue will be created for detailed design. The status of the proposal is divided into: under discussion, finalized and abandoned. After reaching the final status, start specific implementation (PR can also be implemented synchronously during the discussion)

### Final decision maker mechanism

If the maintainer team members have major differences on a requirement, the final decision is made by @Terry Mao.
-->

### Proposal description
<!--
example:
Add event interface for accessing message oriented middleware
-->
### Implementation mode
<!--
```go
example:
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
example:
```go
msg := kafka.NewMessage("kratos", []byte("hello world"), map[string]string{
		"user":  "kratos",
		"phone": "123456",
	})
err := sender.Send(context.Background(), msg)
```
-->