package core

// ================= COMMON =================

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ------------- ASCENDING -------------

func heapSortAsc(data []uint64, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownAsc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDownAsc(data, lo, i, first)
	}
}

func insertionSortAsc(data []uint64, a, b int) {
	var j int
	for i := a + 1; i < b; i++ {
		for j = i; j > a && data[j] < data[j-1]; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func siftDownAsc(data []uint64, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data[first+child] < data[first+child+1] {
			child++
		}
		if data[first+root] >= data[first+child] {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func medianOfThreeAsc(data []uint64, m1, m0, m2 int) {
	// bubble sort on 3 elements
	if data[m1] < data[m0] {
		data[m1], data[m0] = data[m0], data[m1]
	}
	if data[m2] < data[m1] {
		data[m2], data[m1] = data[m1], data[m2]
	}
	if data[m1] < data[m0] {
		data[m1], data[m0] = data[m0], data[m1]
	}
}

func swapRangeAsc(data []uint64, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a], data[b] = data[b], data[a]
		a++
		b++
	}
}

func doPivotAsc(data []uint64, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeAsc(data, lo, lo+s, lo+2*s)
		medianOfThreeAsc(data, m, m-s, m+s)
		medianOfThreeAsc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeAsc(data, lo, m, hi-1)

	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data[b] < data[pivot] {
				b++
			} else if data[pivot] >= data[b] {
				data[a], data[b] = data[b], data[a]
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data[pivot] < data[c-1] {
				c--
			} else if data[c-1] >= data[pivot] {
				data[c-1], data[d-1] = data[d-1], data[c-1]
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	}

	n := min(b-a, a-lo)
	swapRangeAsc(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeAsc(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortAsc(data []uint64, a, b, maxDepth int) {
	var mlo, mhi int
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortAsc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi = doPivotAsc(data, a, b)
		if mlo-a < b-mhi {
			quickSortAsc(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSortAsc(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {
		insertionSortAsc(data, a, b)
	}
}

// Asc asc
func Asc(data []uint64) {
	maxDepth := 0
	for i := len(data); i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortAsc(data, 0, len(data), maxDepth)
}

// IsSortedAsc sorted by Asc
func IsSortedAsc(data []uint64) bool {
	for i := len(data) - 1; i > 0; i-- {
		if data[i] < data[i-1] {
			return false
		}
	}
	return true
}

// StableAsc stable Asc
func StableAsc(data []uint64) {
	n := len(data)
	blockSize := 20
	a, b := 0, blockSize
	for b <= n {
		insertionSortAsc(data, a, b)
		a = b
		b += blockSize
	}
	insertionSortAsc(data, a, n)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMergeAsc(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		symMergeAsc(data, a, a+blockSize, n)
		blockSize *= 2
	}
}

func symMergeAsc(data []uint64, a, m, b int) {
	if a >= m || m >= b {
		return
	}
	mid := a + (b-a)/2
	n := mid + m
	var start, c, r, p int
	if m > mid {
		start = n - b
		r, p = mid, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c] >= data[c] {
				start = c + 1
			} else {
				r = c
			}
		}
	} else {
		start = a
		r, p = m, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c] >= data[c] {
				start = c + 1
			} else {
				r = c
			}
		}
	}
	end := n - start
	rotateAsc(data, start, m, end)
	symMergeAsc(data, a, start, mid)
	symMergeAsc(data, mid, end, b)
}

func rotateAsc(data []uint64, a, m, b int) {
	i := m - a
	if i == 0 {
		return
	}
	j := b - m
	if j == 0 {
		return
	}
	if i == j {
		swapRangeAsc(data, a, m, i)
		return
	}
	p := a + i
	for i != j {
		if i > j {
			swapRangeAsc(data, p-i, p, j)
			i -= j
		} else {
			swapRangeAsc(data, p-i, p+j-i, i)
			j -= i
		}
	}
	swapRangeAsc(data, p-i, p, i)
}

// ------------- DESCENDING -------------

func heapSortDesc(data []uint64, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownDesc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDownDesc(data, lo, i, first)
	}
}

func insertionSortDesc(data []uint64, a, b int) {
	var j int
	for i := a + 1; i < b; i++ {
		for j = i; j > a && data[j] > data[j-1]; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func siftDownDesc(data []uint64, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data[first+child] > data[first+child+1] {
			child++
		}
		if data[first+root] <= data[first+child] {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func medianOfThreeDesc(data []uint64, m1, m0, m2 int) {
	// bubble sort on 3 elements
	if data[m1] > data[m0] {
		data[m1], data[m0] = data[m0], data[m1]
	}
	if data[m2] > data[m1] {
		data[m2], data[m1] = data[m1], data[m2]
	}
	if data[m1] > data[m0] {
		data[m1], data[m0] = data[m0], data[m1]
	}
}

func swapRangeDesc(data []uint64, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a], data[b] = data[b], data[a]
		a++
		b++
	}
}

func doPivotDesc(data []uint64, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeDesc(data, lo, lo+s, lo+2*s)
		medianOfThreeDesc(data, m, m-s, m+s)
		medianOfThreeDesc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeDesc(data, lo, m, hi-1)

	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data[b] > data[pivot] {
				b++
			} else if data[pivot] <= data[b] {
				data[a], data[b] = data[b], data[a]
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data[pivot] > data[c-1] {
				c--
			} else if data[c-1] <= data[pivot] {
				data[c-1], data[d-1] = data[d-1], data[c-1]
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	}

	n := min(b-a, a-lo)
	swapRangeDesc(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeDesc(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortDesc(data []uint64, a, b, maxDepth int) {
	var mlo, mhi int
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortDesc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi = doPivotDesc(data, a, b)
		if mlo-a < b-mhi {
			quickSortDesc(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSortDesc(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {
		insertionSortDesc(data, a, b)
	}
}

// Desc desc
func Desc(data []uint64) {
	maxDepth := 0
	for i := len(data); i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortDesc(data, 0, len(data), maxDepth)
}

// IsSortedDesc  sorted by Desc
func IsSortedDesc(data []uint64) bool {
	for i := len(data) - 1; i > 0; i-- {
		if data[i] > data[i-1] {
			return false
		}
	}
	return true
}

// StableDesc stable desc
func StableDesc(data []uint64) {
	n := len(data)
	blockSize := 20
	a, b := 0, blockSize
	for b <= n {
		insertionSortDesc(data, a, b)
		a = b
		b += blockSize
	}
	insertionSortDesc(data, a, n)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMergeDesc(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		symMergeDesc(data, a, a+blockSize, n)
		blockSize *= 2
	}
}

func symMergeDesc(data []uint64, a, m, b int) {
	if a >= m || m >= b {
		return
	}
	mid := a + (b-a)/2
	n := mid + m
	var start, c, r, p int
	if m > mid {
		start = n - b
		r, p = mid, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c] < data[c] {
				start = c + 1
			} else {
				r = c
			}
		}
	} else {
		start = a
		r, p = m, n-1
		for start < r {
			c = start + (r-start)/2
			if data[p-c] < data[c] {
				start = c + 1
			} else {
				r = c
			}
		}
	}
	end := n - start
	rotateDesc(data, start, m, end)
	symMergeDesc(data, a, start, mid)
	symMergeDesc(data, mid, end, b)
}

func rotateDesc(data []uint64, a, m, b int) {
	i := m - a
	if i == 0 {
		return
	}
	j := b - m
	if j == 0 {
		return
	}
	if i == j {
		swapRangeDesc(data, a, m, i)
		return
	}
	p := a + i
	for i != j {
		if i > j {
			swapRangeDesc(data, p-i, p, j)
			i -= j
		} else {
			swapRangeDesc(data, p-i, p+j-i, i)
			j -= i
		}
	}
	swapRangeDesc(data, p-i, p, i)
}
