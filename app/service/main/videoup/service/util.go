package service

// InSliceIface checks given interface in interface slice.
func InSliceIface(v interface{}, sl []interface{}) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceUnique cleans repeated values in slice.
func SliceUnique(slice []interface{}) (uniqueslice []interface{}) {
	for _, v := range slice {
		if !InSliceIface(v, uniqueslice) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}

// Slice2String convert slice to string
func Slice2String(slice []interface{}) (uniqueslice []string) {
	for _, v := range slice {
		uniqueslice = append(uniqueslice, v.(string))
	}
	return
}

// Slice2Interface convert slice to interface
func Slice2Interface(slice []string) (uniqueslice []interface{}) {
	for _, v := range slice {
		uniqueslice = append(uniqueslice, v)
	}
	return
}
