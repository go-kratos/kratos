package kafka

// Config .
type Config struct {
	Addr       []string
	Topic      string
	Partitions []int32
}
