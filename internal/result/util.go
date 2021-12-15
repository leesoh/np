package result

import (
	"sort"
)

func sortStringMap(m map[string]struct{}) []string {
	var keys []string
	for kk := range m {
		keys = append(keys, kk)
	}
	sort.Strings(keys)
	return keys
}

func sortIntMap(m map[int]struct{}) []int {
	var keys []int
	for kk := range m {
		keys = append(keys, kk)
	}
	sort.Ints(keys)
	return keys
}
