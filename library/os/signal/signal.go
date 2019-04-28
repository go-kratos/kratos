package signal

import (
	"os"
	"os/signal"
)

// Notify redirect os/signal.Notify
func Notify(c chan<- os.Signal, sig ...os.Signal) {
	signal.Notify(c, sig...)
}
