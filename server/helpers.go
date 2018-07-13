package server

import (
	"strconv"
	"strings"
)

func getCurrentHistoricIDs(rawIDs []string) ([]int64, [][2]int64, error) {
	currentIDs := make([]int64, 0)
	historicIDs := make([][2]int64, 0)
	for i := range rawIDs {
		idv := strings.Split(rawIDs[i], "v")
		id, err := strconv.ParseInt(idv[0], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		if len(idv) == 1 {
			currentIDs = appendIfUnique(currentIDs, id)
			continue
		}
		v, err := strconv.ParseInt(idv[1], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		historicIDs = appendVersionIfUnique(historicIDs, [2]int64{id, v})
	}
	return currentIDs, historicIDs, nil
}

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
