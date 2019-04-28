package service

type Dispatcher struct {
	cache map[uint64]string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		cache: make(map[uint64]string),
	}
}

func (d *Dispatcher) AccessNode(channelID uint64) (string, error) {
	return "127.0.0.1", nil
}
