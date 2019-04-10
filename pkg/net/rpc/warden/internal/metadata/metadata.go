package metadata

const (
	CPUUsage = "cpu_usage"
)

// MD is context metadata for balancer and resolver
type MD struct {
	Weight uint64
	Color  string
}
