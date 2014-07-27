package triplebuffer

import (
	"math"
	"sort"
)

func uniqifyInts(ints []int) (uniq_ints []int) {
	sort.Ints(ints)
	uniq_ints = make([]int, 0, len(ints))
	found := false
	for _, i := range ints {
		found = false
		for _, j := range uniq_ints {
			if i == j {
				found = true
				break
			}
		}
		if !found {
			uniq_ints = append(uniq_ints, i)
		}
	}
	return
}

func vectorEqual(a, b []float64) bool {
	FLOAT_EQUALITY_THRESHOLD := 0.0000001
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if math.Abs(a[i]-b[i]) > FLOAT_EQUALITY_THRESHOLD {
			return false
		}
	}
	return true
}

func intVectorEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
