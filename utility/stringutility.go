package utility

import (
	"strconv"
	"strings"
)

func StringToIntArray(s string, sep string) ([]int, error) {
	if len(s) == 0 {
		return make([]int, 0), nil
	}
	strarry := strings.Split(s, sep)
	dist := make([]int, len(strarry))

	for i, item := range strarry {
		result, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		dist[i] = result
	}
	return dist, nil
}
