package archive

// Databus db
type Databus struct {
	ID        int64
	Group     string
	Topic     string
	Partition int8
	Offset    int64
}
