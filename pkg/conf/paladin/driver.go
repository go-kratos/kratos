package paladin

// Driver defined paladin remote client impl
// each remote config center driver must do
// 1. implements `New` method
// 2. call `Register` to register itself
type Driver interface {
	New() (Client, error)
}
