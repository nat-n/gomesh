package mesh

import (
	"strconv"
)

// Parses a triple of strings into floats
func parse3Floats(strings []string) (result [3]float64, err error) {
	result[0], err = strconv.ParseFloat(strings[0], 64)
	if err != nil {
		return
	}
	result[1], err = strconv.ParseFloat(strings[1], 64)
	if err != nil {
		return
	}
	result[2], err = strconv.ParseFloat(strings[2], 64)
	if err != nil {
		return
	}
	return
}

// Parses a triple of strings into ints
func parse3Ints(strings []string) (result [3]int, err error) {
	int64s := [3]int64{}
	int64s[0], err = strconv.ParseInt(strings[0], 10, 64)
	if err != nil {
		return
	}
	result[0] = int(int64s[0])
	int64s[1], err = strconv.ParseInt(strings[1], 10, 64)
	if err != nil {
		return
	}
	result[1] = int(int64s[1])
	int64s[2], err = strconv.ParseInt(strings[2], 10, 64)
	if err != nil {
		return
	}
	result[2] = int(int64s[2])
	return
}
