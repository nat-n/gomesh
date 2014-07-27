package triplebuffer

import (
	"testing"
)

type intsTestParams struct {
	tripleIntBuffers    []tripleIntBuffer
	indices             []int
	ints                []int
	collectFunction     func(int, int, int) (int, int, int)
	eachFunctionFactory func(tb *tripleIntBuffer) func(int, int, int)
	resultInts          []int
	resultIntBuffers    []tripleIntBuffer
	resultMap           map[int][]int
}

// indexing, collect, each

// Tests for tripleIntBuffer.Len

var testIntLen = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3},
				make(map[int][]int),
			},
		},
		resultInts: []int{1},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 4, 5, 6},
				make(map[int][]int),
			},
		},
		resultInts: []int{2},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 4, 5, 6, 6, 7, 9},
				make(map[int][]int),
			},
		},
		resultInts: []int{3},
	},
}

func TestIntLen(t *testing.T) {
	for _, params := range testIntLen {
		r := params.tripleIntBuffers[0].Len()
		if r != params.resultInts[0] {
			t.Error(
				"For Len of", params.tripleIntBuffers[0],
				"expected", params.resultInts[0],
				"got", r,
			)
		}
	}
}

// Tests for tripleIntBuffer.Get

var testsIntGet = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3},
				make(map[int][]int),
			},
		},
		indices:    []int{0},
		resultInts: []int{1, 2, 3},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 4, 52, 6, 0, 92, 84},
				make(map[int][]int),
			},
		},
		indices:    []int{2, 1},
		resultInts: []int{0, 92, 84, 4, 52, 6},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, -6, 367, -83, 8, -42435, 6, 7, 9},
				make(map[int][]int),
			},
		},
		indices:    []int{0, 2},
		resultInts: []int{123, -6, 367, 6, 7, 9},
	},
}

func TestIntGet(t *testing.T) {
	for _, params := range testsIntGet {
		r := params.tripleIntBuffers[0].Get(params.indices...)
		if !intVectorEqual(r, params.resultInts) {
			t.Error(
				"For items from", params.tripleIntBuffers[0],
				"expected", params.resultInts,
				"got", r,
			)
		}
	}
}

// Tests for tripleIntBuffer.Append

var testIntAppend = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3},
				make(map[int][]int),
			},
		},
		ints: []int{1, 2, 3, 4, 52, 6, 0, 92, 84},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 1, 2, 3, 4, 52, 6, 0, 92, 84},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 4, 52, 6, 0, 92, 84},
				make(map[int][]int),
			},
		},
		ints: []int{2, 1, 33333},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 4, 52, 6, 0, 92, 84, 2, 1, 33333},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, -6, 367, -83, 8, -42435, 6, 7, 9},
				make(map[int][]int),
			},
		},
		ints: []int{-6, 367, -83},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, -6, 367, -83, 8, -42435, 6, 7, 9, -6, 367, -83},
				make(map[int][]int),
			},
		},
	},
}

func TestIntAppend(t *testing.T) {
	for _, params := range testIntAppend {
		params.tripleIntBuffers[0].Append(params.ints...)
		if !intVectorEqual(params.tripleIntBuffers[0].Buffer, params.resultIntBuffers[0].Buffer) {
			t.Error(
				"For Appending to tripleIntBuffer",
				"expected", params.resultIntBuffers[0],
				"got", params.tripleIntBuffers[0],
			)
		}
	}
}

// Tests for tripleIntBuffer.UpdateOne

var testIntUpdateOne = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3},
				make(map[int][]int),
			},
		},
		indices: []int{0},
		ints:    []int{1, 2, 3333},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3333},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 4, 52, 6, 0, 92, 84},
				make(map[int][]int),
			},
		},
		indices: []int{1},
		ints:    []int{2, 1, 33333},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 2, 1, 33333, 0, 92, 84},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, -6, 367, -83, 8, -42435, 6, 7, 9},
				make(map[int][]int),
			},
		},
		indices: []int{2},
		ints:    []int{-6, 367, -83},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, -6, 367, -83, 8, -42435, -6, 367, -83},
				make(map[int][]int),
			},
		},
	},
}

func TestIntUpdateOne(t *testing.T) {
	for _, params := range testIntUpdateOne {
		params.tripleIntBuffers[0].UpdateOne(params.indices[0], params.ints[0], params.ints[1], params.ints[2])
		if !intVectorEqual(params.tripleIntBuffers[0].Buffer, params.resultIntBuffers[0].Buffer) {
			t.Error(
				"For Updating single item in tripleIntBuffer",
				"expected", params.resultIntBuffers[0],
				"got", params.tripleIntBuffers[0],
			)
		}
	}
}

// Tests for tripleIntBuffer.Update

var testIntUpdate = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3},
				make(map[int][]int),
			},
		},
		indices: []int{0},
		ints:    []int{1, 2, 3333},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3333},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 4, 52, 6, 0, 92, 84},
				make(map[int][]int),
			},
		},
		indices: []int{2, 0},
		ints:    []int{2, 1, 33333, 367, -83, 8},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{367, -83, 8, 4, 52, 6, 2, 1, 33333},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, -6, 367, -83, 8, -42435, 6, 7, 9},
				make(map[int][]int),
			},
		},
		indices: []int{0, 1, 2},
		ints:    []int{-6, 367, -83, 2, 1, 33333, 0, 92, 84},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{-6, 367, -83, 2, 1, 33333, 0, 92, 84},
				make(map[int][]int),
			},
		},
	},
}

func TestIntUpdate(t *testing.T) {
	var updatesMap map[int][]int
	for _, params := range testIntUpdate {
		updatesMap = make(map[int][]int)
		for i, j := range params.indices {
			updatesMap[j] = []int{
				params.ints[i*3],
				params.ints[i*3+1],
				params.ints[i*3+2],
			}[:]
		}
		params.tripleIntBuffers[0].Update(updatesMap)
		if !intVectorEqual(params.tripleIntBuffers[0].Buffer, params.resultIntBuffers[0].Buffer) {
			t.Error(
				"For Updating multiple triples in tripleIntBuffer",
				"expected", params.resultIntBuffers[0],
				"got", params.tripleIntBuffers[0],
			)
		}
	}
}

// Tests for tripleIntBuffer.RemoveOne

var testIntRemoveOne = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3},
				make(map[int][]int),
			},
		},
		indices: []int{0},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 0, 92, 84},
				make(map[int][]int),
			},
		},
		indices: []int{1},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, -6, 367, -83, 8, -42435, 6, 7, 9},
				make(map[int][]int),
			},
		},
		indices: []int{1},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, -6, 367, 6, 7, 9},
				make(map[int][]int),
			},
		},
	},
}

func TestIntRemoveOne(t *testing.T) {
	for _, params := range testIntRemoveOne {
		params.tripleIntBuffers[0].RemoveOne(params.indices[0])
		if !intVectorEqual(params.tripleIntBuffers[0].Buffer, params.resultIntBuffers[0].Buffer) {
			t.Error(
				"For Removing a single triple from tripleIntBuffer",
				"expected", params.resultIntBuffers[0],
				"got", params.tripleIntBuffers[0],
			)
		}
	}
}

// Tests for tripleIntBuffer.Remove

var testIntRemove = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 2, 3, 6, 4, 23234234},
				make(map[int][]int),
			},
		},
		indices: []int{0},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{6, 4, 23234234},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 0, 3, 0, 92, 84, 643, 6234, 123, 17, 18, 7},
				make(map[int][]int),
			},
		},
		indices: []int{2, 1, 2},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 0, 3, 17, 18, 7},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{123, 41, -632, 8, 367333, -83224, 87225, -424, -352, 6, 7, 9},
				make(map[int][]int),
			},
		},
		indices: []int{2, 0},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{8, 367333, -83224, 6, 7, 9},
				make(map[int][]int),
			},
		},
	},
}

func TestIntRemove(t *testing.T) {
	for _, params := range testIntRemove {
		params.tripleIntBuffers[0].Remove(params.indices...)
		if !intVectorEqual(params.tripleIntBuffers[0].Buffer, params.resultIntBuffers[0].Buffer) {
			t.Error(
				"For Removing multiple triples from tripleIntBuffer",
				"expected", params.resultIntBuffers[0],
				"got", params.tripleIntBuffers[0],
			)
		}
	}
}

// Tests for tripleIntBuffer.UpdateIndex

var testIntUpdateIndex = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 63, 63, 2, 3, 63, 4, 63, 232342342},
				make(map[int][]int),
			},
		},
		ints:       []int{63},
		resultInts: []int{0, 0, 1, 2},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{15, 2, 32, 07, 92, 844, 6432, 06234, 123132, 17, 02, 76},
				make(map[int][]int),
			},
		},
		ints:       []int{2},
		resultInts: []int{0, 3},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{12323, 41, -632, 8, 367333, -83224, 87225, -424, -352, 6, 7, 9},
				make(map[int][]int),
			},
		},
		ints:       []int{-424},
		resultInts: []int{2},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{12323, 41, -632, 8, 367333, -83224, 87225, -424, -352, 6, 7, 9},
				make(map[int][]int),
			},
		},
		ints:       []int{435},
		resultInts: []int{},
	},
}

func TestIntUpdateIndex(t *testing.T) {
	for _, params := range testIntUpdateIndex {
		params.tripleIntBuffers[0].UpdateIndex()
		if !intVectorEqual(params.tripleIntBuffers[0].Index[params.ints[0]], params.resultInts) {
			t.Error(
				"For Removing multiple triples from tripleIntBuffer",
				"expected", params.resultIntBuffers[0],
				"got", params.tripleIntBuffers[0],
			)
		}
	}
}

// Tests for tripleIntBuffer.Collect

var testIntCollect = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 63, 63, 2, 3, 63, 4, 63, 232342342},
				make(map[int][]int),
			},
		},
		collectFunction: func(x, y, z int) (a, b, c int) {
			a, b, c = y, z, x
			return
		},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{63, 63, 1, 3, 63, 2, 63, 232342342, 4},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{15, 2, 32, 7, 92, 844, 6402, 6234, 123132, 17, 2, 76},
				make(map[int][]int),
			},
		},
		collectFunction: func(x, y, z int) (a, b, c int) {
			a = x + 1
			b = y + 2
			z = 0
			return
		},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{25, 22, 0, 17, 94, 0, 6412, 26234, 0, 18, 22, 0},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{12323, 41, -632, 8, 367333, -83224, 87225, -424, -352, 6, 7, 9},
				make(map[int][]int),
			},
		},
		collectFunction: func(x, y, z int) (a, b, c int) {
			a = x / y
			b = -y
			c = z * z
			return
		},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{5609756, -41, 399424, 21778604, -367333, 6234176, -20571934, 424, 123904, 857142857, -7, 81},
				make(map[int][]int),
			},
		},
	},
}

// Tests for tripleIntBuffer.Each

var testIntEach = []intsTestParams{
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{1, 63, 63, 2, 3, 63, 4, 63, 232342342},
				make(map[int][]int),
			},
		},
		eachFunctionFactory: func(tb *tripleIntBuffer) func(int, int, int) {
			return func(x, y, z int) {
				tb.Append(y, z, x)
			}
		},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{63, 63, 1, 3, 63, 2, 63, 232342342, 4},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{15, 2, 32, 7, 92, 844, 6402, 6234, 123132, 17, 2, 76},
				make(map[int][]int),
			},
		},
		eachFunctionFactory: func(tb *tripleIntBuffer) func(int, int, int) {
			return func(x, y, z int) {
				tb.Append(x+1, y+2, 0)
			}
		},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{16, 4, 0, 8, 94, 0, 6403, 6236, 0, 18, 4, 0},
				make(map[int][]int),
			},
		},
	},
	{
		tripleIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{12323, 41, -632, 8, 367333, -832, 87225, -424, -352, 6, 7, 9},
				make(map[int][]int),
			},
		},
		eachFunctionFactory: func(tb *tripleIntBuffer) func(int, int, int) {
			return func(x, y, z int) {
				tb.Append(x/y, -y, z*z)
			}
		},
		resultIntBuffers: []tripleIntBuffer{
			tripleIntBuffer{
				[]int{300, -41, 399424, 0, -367333, 692224, -205, 424, 123904, 0, -7, 81},
				make(map[int][]int),
			},
		},
	},
}

func TestIntEach(t *testing.T) {
	var dest_buffer tripleIntBuffer
	for _, params := range testIntEach {
		dest_buffer = tripleIntBuffer{
			[]int{},
			make(map[int][]int),
		}
		params.tripleIntBuffers[0].Each(params.eachFunctionFactory(&dest_buffer))
		if !intVectorEqual(dest_buffer.Buffer, params.resultIntBuffers[0].Buffer) {
			t.Error(
				"For iterating over each triple in a tripleIntBuffer",
				"expected", params.resultIntBuffers[0],
				"got", dest_buffer,
			)
		}
	}
}
