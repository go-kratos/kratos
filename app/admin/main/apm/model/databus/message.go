package databus

// Message Data.
type Message struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Timestamp int64  `json:"timestamp"`
}
