package reply

// AscFloors floors sort.
type AscFloors []int64

func (f AscFloors) Len() int {
	return len(f)
}

func (f AscFloors) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f AscFloors) Less(i, j int) bool {
	return f[i] < f[j]
}

// DescFloors floors sort.
type DescFloors []int64

func (f DescFloors) Len() int {
	return len(f)
}

func (f DescFloors) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f DescFloors) Less(i, j int) bool {
	return f[i] > f[j]
}
