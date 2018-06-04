package db

import (
	"fmt"
	"strings"
)

// arrayToString converts array to string
func arrayToString(arr []int64) string {
	return strings.Trim(strings.Replace(fmt.Sprint(arr), " ", ",", -1), "[]")
}

// versionArrayToString converts version array to string
func versionArrayToString(arr [][2]int64) string {
	return strings.Trim(strings.Replace(fmt.Sprint(arr), " ", ",", -1), "[]")
}
