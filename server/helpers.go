package server

func appendIfUnique(slice []int64, i int64) []int64 {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func appendVersionIfUnique(slice [][2]int64, i [2]int64) [][2]int64 {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
