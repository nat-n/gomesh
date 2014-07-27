package triplebuffer

import (
	"math"
	"strconv"
	"strings"
)

type tripleFloatBuffer struct {
	Buffer []float64
	Index  map[float64][]int
}

func (tb *tripleFloatBuffer) Len() int {
	return len(tb.Buffer) / 3
}

func (tb *tripleFloatBuffer) Get(indices ...int) []float64 {
	result := make([]float64, 0, len(indices)*3)
	for _, i := range indices {
		result = append(result, tb.Buffer[i*3], tb.Buffer[i*3+1], tb.Buffer[i*3+2])
	}
	return result
}

func (tb *tripleFloatBuffer) GetAsMap(indices ...int) map[int][3]float64 {
	result := make(map[int][3]float64)
	for _, i := range indices {
		result[i] = [3]float64{tb.Buffer[i*3], tb.Buffer[i*3+1], tb.Buffer[i*3+2]}
	}
	return result
}

func (tb *tripleFloatBuffer) Append(values ...float64) {
	if len(values)%3 != 0 {
		panic("Cannot append array of length " +
			strconv.Itoa(len(values)) +
			" to tripleFloatBuffer.")
	}
	tb.Buffer = append(tb.Buffer, values...)
}

func (tb *tripleFloatBuffer) UpdateOne(i int, x, y, z float64) {
	tb.Buffer[i*3] = x
	tb.Buffer[i*3+1] = y
	tb.Buffer[i*3+2] = z
}

func (tb *tripleFloatBuffer) Update(updates map[int][]float64) {
	for i, values := range updates {
		tb.Buffer[i*3] = values[0]
		tb.Buffer[i*3+1] = values[1]
		tb.Buffer[i*3+2] = values[2]
	}
}

// func (tb *tripleFloatBuffer) updateMultiple(updates map[int][]float64) {

// }

func (tb *tripleFloatBuffer) RemoveOne(i int) {
	if i < 0 || i > len(tb.Buffer)/3 {
		panic("Cannot remove item " + strconv.Itoa(i) +
			" from a tripleFloatBuffer of length " +
			strconv.Itoa(len(tb.Buffer)/3) + ".")
	}
	tb.Buffer = append(tb.Buffer[:i*3], tb.Buffer[i*3+3:]...)
	return
}

// This method works by copying the contents of the Buffer once, and filtering
// out any triples referenced in indices.
func (tb *tripleFloatBuffer) Remove(indices ...int) {
	// ensure indices are valid and sorted and unique
	for _, i := range indices {
		if i < 0 || i > len(tb.Buffer)/3 {
			panic("Cannot remove item " + strconv.Itoa(i) +
				" from a tripleFloatBuffer of length " +
				strconv.Itoa(len(tb.Buffer)/3) + ".")
		}
	}
	indices = uniqifyInts(indices)

	new_len := len(tb.Buffer) - len(indices)*3
	new_Buffer := make([]float64, 0, new_len)
	previous := indices[0]
	// copy over triples before the first index
	if indices[0] > 0 {
		new_Buffer = append(new_Buffer, tb.Buffer[:previous*3]...)
	}

	// copy over triples between indices
	for _, i := range indices {
		if i > previous {
			new_Buffer = append(new_Buffer, tb.Buffer[(previous)*3:i*3]...)
		}
		previous = i + 1
	}

	// copy over triples after the last index
	if previous <= new_len {
		new_Buffer = append(new_Buffer, tb.Buffer[(previous)*3:]...)
	}

	tb.Buffer = new_Buffer
}

// This should be smarter about building up each array gradually!
// ... although the arrays wont get very big.
func (tb *tripleFloatBuffer) UpdateIndex() {
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
func (tb *tripleFloatBuffer) Collect(f func(float64, float64, float64) (float64, float64, float64)) {
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
func (tb *tripleFloatBuffer) Each(f func(float64, float64, float64)) {
	tripleCount := tb.Len()
	var index int
	for i := 0; i < tripleCount; i++ {
		index = i * 3
		f(tb.Buffer[index], tb.Buffer[index+1], tb.Buffer[index+2])
	}
}

// Loop with a function over the triples in the Buffer
func (tb *tripleFloatBuffer) EachWithIndex(f func(int, float64, float64, float64)) {
	tripleCount := tb.Len()
	var index int
	for i := 0; i < tripleCount; i++ {
		index = i * 3
		f(i, tb.Buffer[index], tb.Buffer[index+1], tb.Buffer[index+2])
	}
}

func (tb *tripleFloatBuffer) TriplesWith(value float64) []float64 {
	indices, found := tb.Index[value]
	if found {
		indices = uniqifyInts(indices)
		return tb.Get(indices...)
	} else {
		return make([]float64, 0)
	}
}

func (tb *tripleFloatBuffer) Equal(other tripleFloatBuffer) bool {
	FLOAT_EQUALITY_THRESHOLD := 0.0000001
	if len(tb.Buffer) != len(other.Buffer) {
		return false
	}
	for i, _ := range tb.Buffer {
		if math.Abs(tb.Buffer[i]-other.Buffer[i]) > FLOAT_EQUALITY_THRESHOLD {
			return false
		}
	}
	return true
}

func (tb *tripleFloatBuffer) IsEmpty() bool {
	for n, _ := range tb.Buffer {
		if n > 0 {
			return false
		}
	}
	return true
}

func (tb *tripleFloatBuffer) Average() (x, y, z float64) {
	tripleCount := tb.Len()
	var index int
	for i := 0; i < tripleCount; i++ {
		index = i * 3
		x += tb.Buffer[index]
		y += tb.Buffer[index+1]
		z += tb.Buffer[index+2]
	}
	x /= float64(tripleCount)
	y /= float64(tripleCount)
	z /= float64(tripleCount)
	return
}

func (tb *tripleFloatBuffer) ToString() string {
	stringFloats := make([]string, len(tb.Buffer), len(tb.Buffer))
	for i := 0; i < len(tb.Buffer); i++ {
		stringFloats[i] = strconv.FormatFloat(tb.Buffer[i], 'f', -1, 64)
	}
	return strings.Join(stringFloats, ",")
}
