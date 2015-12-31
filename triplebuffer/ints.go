package triplebuffer

import (
	"strconv"
	"strings"
)

type tripleIntBuffer struct {
	Buffer []int
	Index  map[int][]int
}

func (tb *tripleIntBuffer) Len() int {
	return len(tb.Buffer) / 3
}

func (tb *tripleIntBuffer) Get(indices ...int) []int {
	result := make([]int, 0, len(indices)*3)
	for _, i := range indices {
		result = append(result, tb.Buffer[i*3], tb.Buffer[i*3+1], tb.Buffer[i*3+2])
	}
	return result
}

func (tb *tripleIntBuffer) GetAsMap(indices ...int) map[int][3]int {
	result := make(map[int][3]int)
	for _, i := range indices {
		result[i] = [3]int{tb.Buffer[i*3], tb.Buffer[i*3+1], tb.Buffer[i*3+2]}
	}
	return result
}

func (tb *tripleIntBuffer) Append(values ...int) {
	if len(values)%3 != 0 {
		panic("Cannot Append array of length " +
			strconv.Itoa(len(values)) +
			" to tripleIntBuffer.")
	}
	tb.Buffer = append(tb.Buffer, values...)
}

func (tb *tripleIntBuffer) UpdateOne(i int, x, y, z int) {
	tb.Buffer[i*3] = x
	tb.Buffer[i*3+1] = y
	tb.Buffer[i*3+2] = z
}

func (tb *tripleIntBuffer) Update(updates map[int][]int) {
	for i, values := range updates {
		tb.Buffer[i*3] = values[0]
		tb.Buffer[i*3+1] = values[1]
		tb.Buffer[i*3+2] = values[2]
	}
}

// func (tb *tripleIntBuffer) updateMultiple(updates map[int][]int) {

// }

func (tb *tripleIntBuffer) RemoveOne(i int) {
	if i < 0 || i > len(tb.Buffer)/3 {
		panic("Cannot remove item " + strconv.Itoa(i) +
			" from a tripleIntBuffer of length " +
			strconv.Itoa(len(tb.Buffer)/3) + ".")
	}
	tb.Buffer = append(tb.Buffer[:i*3], tb.Buffer[i*3+3:]...)
	return
}

// This method works by copying the contents of the Buffer once, and filtering
// out any triples referenced in indices.
func (tb *tripleIntBuffer) Remove(indices ...int) {
	if len(indices) == 0 {
		return
	}
	// ensure indices are valid and sorted and unique
	for _, i := range indices {
		if i < 0 || i > len(tb.Buffer)/3 {
			panic("Cannot remove item " + strconv.Itoa(i) +
				" from a tripleIntBuffer of length " +
				strconv.Itoa(len(tb.Buffer)/3) + ".")
		}
	}
	indices = uniqifyInts(indices)

	new_len := len(tb.Buffer) - len(indices)*3
	new_Buffer := make([]int, 0, new_len)
	previous := indices[0]
	// copy over triples before the first index
	if previous > 0 {
		new_Buffer = append(new_Buffer, tb.Buffer[:previous*3]...)
	}

	// copy over triples between indices
	for _, i := range indices {
		if i > previous {
			new_Buffer = append(new_Buffer, tb.Buffer[previous*3:i*3]...)
		}
		previous = i + 1
	}

	// copy over triples after the last index
	if len(new_Buffer) < new_len {
		new_Buffer = append(new_Buffer, tb.Buffer[previous*3:]...)
	}

	tb.Buffer = new_Buffer
}

func (tb *tripleIntBuffer) UpdateIndex() {
	tb.Index = make(map[int][]int)
	len3 := tb.Len()
	var index int
	for i := 0; i < len3; i++ {
		index = i * 3
		tb.Index[tb.Buffer[index]] = append(tb.Index[tb.Buffer[index]], i)
		tb.Index[tb.Buffer[index+1]] = append(tb.Index[tb.Buffer[index+1]], i)
		tb.Index[tb.Buffer[index+2]] = append(tb.Index[tb.Buffer[index+2]], i)
	}
}

// Collect/map over the triples in the Buffer
func (tb *tripleIntBuffer) Collect(f func(int, int, int) (int, int, int)) {
	len3 := tb.Len()
	var index int
	for i := 0; i < len3; i++ {
		index = i * 3
		tb.Buffer[index], tb.Buffer[index+1], tb.Buffer[index+2] = f(
			tb.Buffer[index], tb.Buffer[index+1], tb.Buffer[index+2],
		)
	}
}

// Loop with a function over the triples in the Buffer
func (tb *tripleIntBuffer) Each(f func(int, int, int)) {
	len3 := tb.Len()
	var index int
	for i := 0; i < len3; i++ {
		index = i * 3
		f(tb.Buffer[index], tb.Buffer[index+1], tb.Buffer[index+2])
	}
}

// Loop with a function over the indexed triples in the Buffer
func (tb *tripleIntBuffer) EachOf(f func(int, int, int), indices ...int) {
	for _, i := range indices {
		f(tb.Buffer[i*3], tb.Buffer[i*3+1], tb.Buffer[i*3+2])
	}
}

func (tb *tripleIntBuffer) EachWithIndex(f func(int, int, int, int)) {
	len3 := tb.Len()
	var index int
	for i := 0; i < len3; i++ {
		index = i * 3
		f(i, tb.Buffer[index], tb.Buffer[index+1], tb.Buffer[index+2])
	}
}

// Loop with a function over the indexed triples in the Buffer
func (tb *tripleIntBuffer) EachOfWithIndex(f func(int, int, int, int), indices ...int) {
	for _, i := range indices {
		f(i, tb.Buffer[i*3], tb.Buffer[i*3+1], tb.Buffer[i*3+2])
	}
}

func (tb *tripleIntBuffer) TriplesWith(value int) []int {
	indices, found := tb.Index[value]
	if found {
		indices = uniqifyInts(indices)
		return tb.Get(indices...)
	} else {
		return make([]int, 0)
	}
}

func (tb *tripleIntBuffer) Equal(other tripleIntBuffer) bool {
	if len(tb.Buffer) != len(other.Buffer) {
		return false
	}
	for i, _ := range tb.Buffer {
		if tb.Buffer[i] != other.Buffer[i] {
			return false
		}
	}
	return true
}

func (tb *tripleIntBuffer) ToString() string {
	stringInts := make([]string, len(tb.Buffer), len(tb.Buffer))
	for i := 0; i < len(tb.Buffer); i++ {
		stringInts[i] = strconv.FormatInt(int64(tb.Buffer[i]), 10)
	}
	return strings.Join(stringInts, ",")
}
