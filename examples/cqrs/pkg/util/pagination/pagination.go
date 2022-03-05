package pagination

func GetPageOffset(pageNum, pageSize int32) int {
	return int((pageNum - 1) * pageSize)
}
