package triplebuffer

import (
	"testing"
)

type floatTestParams struct {
	tripleFloatBuffers  []tripleFloatBuffer
	indices             []int
	floats              []float64
	collectFunction     func(float64, float64, float64) (float64, float64, float64)
	eachFunctionFactory func(tb *tripleFloatBuffer) func(float64, float64, float64)
	resultInts          []int
	resultSlice         []float64
	resultFloatBuffers  []tripleFloatBuffer
	resultMap           map[int][]float64
}

// Tests for tripleFloatBuffer.Len

var testLen = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3},
				make(map[float64][]int),
			},
		},
		resultInts: []int{1},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3, 4, 5, 6},
				make(map[float64][]int),
			},
		},
		resultInts: []int{2},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3, 4, 5, 6, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		resultInts: []int{3},
	},
}

func TestLen(t *testing.T) {
	for _, params := range testLen {
		r := params.tripleFloatBuffers[0].Len()
		if r != params.resultInts[0] {
			t.Error(
				"For Len of", params.tripleFloatBuffers[0],
				"expected", params.resultInts[0],
				"got", r,
			)
		}
	}
}

// Tests for tripleFloatBuffer.Get

var testsGet = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3},
				make(map[float64][]int),
			},
		},
		indices:     []int{0},
		resultSlice: []float64{1, 2, 3},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3.2, 4.5, 52, 6.1, 0.7, 92, 84.4},
				make(map[float64][]int),
			},
		},
		indices:     []int{2, 1},
		resultSlice: []float64{0.7, 92, 84.4, 4.5, 52, 6.1},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.2341, -6.32, 367.333, -83.224, 8.7225, -42435.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		indices:     []int{0, 2},
		resultSlice: []float64{123.2341, -6.32, 367.333, 6, 7, 9},
	},
}

func TestGet(t *testing.T) {
	for _, params := range testsGet {
		r := params.tripleFloatBuffers[0].Get(params.indices...)
		if !vectorEqual(r, params.resultSlice) {
			t.Error(
				"For items from", params.tripleFloatBuffers[0],
				"expected", params.resultSlice,
				"got", r,
			)
		}
	}
}

// Tests for tripleFloatBuffer.Append

var testAppend = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3},
				make(map[float64][]int),
			},
		},
		floats: []float64{1.5, 2, 3.2, 4.5, 52, 6.1, 0.7, 92, 84.4},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3, 1.5, 2, 3.2, 4.5, 52, 6.1, 0.7, 92, 84.4},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3.2, 4.5, 52, 6.1, 0.7, 92, 84.4},
				make(map[float64][]int),
			},
		},
		floats: []float64{2, 1, 33333.333},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3.2, 4.5, 52, 6.1, 0.7, 92, 84.4, 2, 1, 33333.333},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.2341, -6.32, 367.333, -83.224, 8.7225, -42435.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		floats: []float64{-6.32, 367.333, -83.224},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.2341, -6.32, 367.333, -83.224, 8.7225, -42435.2, 6, 7, 9, -6.32, 367.333, -83.224},
				make(map[float64][]int),
			},
		},
	},
}

func TestAppend(t *testing.T) {
	for _, params := range testAppend {
		params.tripleFloatBuffers[0].Append(params.floats...)
		if !vectorEqual(params.tripleFloatBuffers[0].Buffer, params.resultFloatBuffers[0].Buffer) {
			t.Error(
				"For Appending to tripleFloatBuffer",
				"expected", params.resultFloatBuffers[0],
				"got", params.tripleFloatBuffers[0],
			)
		}
	}
}

// Tests for tripleFloatBuffer.UpdateOne

var testUpdateOne = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3},
				make(map[float64][]int),
			},
		},
		indices: []int{0},
		floats:  []float64{1.5, 2, 3333.23424},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3333.23424},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3.2, 4.5, 52, 6.1, 0.7, 92, 84.4},
				make(map[float64][]int),
			},
		},
		indices: []int{1},
		floats:  []float64{2, 1, 33333.333},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3.2, 2, 1, 33333.333, 0.7, 92, 84.4},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.2341, -6.32, 367.333, -83.224, 8.7225, -42435.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		indices: []int{2},
		floats:  []float64{-6.32, 367.333, -83.224},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.2341, -6.32, 367.333, -83.224, 8.7225, -42435.2, -6.32, 367.333, -83.224},
				make(map[float64][]int),
			},
		},
	},
}

func TestUpdateOne(t *testing.T) {
	for _, params := range testUpdateOne {
		params.tripleFloatBuffers[0].UpdateOne(params.indices[0], params.floats[0], params.floats[1], params.floats[2])
		if !vectorEqual(params.tripleFloatBuffers[0].Buffer, params.resultFloatBuffers[0].Buffer) {
			t.Error(
				"For Updating single item in tripleFloatBuffer",
				"expected", params.resultFloatBuffers[0],
				"got", params.tripleFloatBuffers[0],
			)
		}
	}
}

// Tests for tripleFloatBuffer.Update

var testUpdate = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3},
				make(map[float64][]int),
			},
		},
		indices: []int{0},
		floats:  []float64{1.5, 2, 3333.23424},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3333.23424},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3.2, 4.5, 52, 6.1, 0.7, 92, 84.4},
				make(map[float64][]int),
			},
		},
		indices: []int{2, 0},
		floats:  []float64{2, 1, 33333.333, 367.333, -83.224, 8.7225},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{367.333, -83.224, 8.7225, 4.5, 52, 6.1, 2, 1, 33333.333},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.2341, -6.32, 367.333, -83.224, 8.7225, -42435.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		indices: []int{0, 1, 2},
		floats:  []float64{-6.32, 367.333, -83.224, 2, 1, 33333.333, 0.7, 92, 84.4},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{-6.32, 367.333, -83.224, 2, 1, 33333.333, 0.7, 92, 84.4},
				make(map[float64][]int),
			},
		},
	},
}

func TestUpdate(t *testing.T) {
	var updatesMap map[int][]float64
	for _, params := range testUpdate {
		updatesMap = make(map[int][]float64)
		for i, j := range params.indices {
			updatesMap[j] = []float64{
				params.floats[i*3],
				params.floats[i*3+1],
				params.floats[i*3+2],
			}[:]
		}
		params.tripleFloatBuffers[0].Update(updatesMap)
		if !vectorEqual(params.tripleFloatBuffers[0].Buffer, params.resultFloatBuffers[0].Buffer) {
			t.Error(
				"For Updating multiple triples in tripleFloatBuffer",
				"expected", params.resultFloatBuffers[0],
				"got", params.tripleFloatBuffers[0],
			)
		}
	}
}

// Tests for tripleFloatBuffer.RemoveOne

var testRemoveOne = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3},
				make(map[float64][]int),
			},
		},
		indices: []int{0},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3.2, 0.7, 92, 84.4},
				make(map[float64][]int),
			},
		},
		indices: []int{1},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 2, 3.2},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.2341, -6.32, 367.333, -83.224, 8.7225, -42435.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		indices: []int{1},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.2341, -6.32, 367.333, 6, 7, 9},
				make(map[float64][]int),
			},
		},
	},
}

func TestRemoveOne(t *testing.T) {
	for _, params := range testRemoveOne {
		params.tripleFloatBuffers[0].RemoveOne(params.indices[0])
		if !vectorEqual(params.tripleFloatBuffers[0].Buffer, params.resultFloatBuffers[0].Buffer) {
			t.Error(
				"For Removing a single triple from tripleFloatBuffer",
				"expected", params.resultFloatBuffers[0],
				"got", params.tripleFloatBuffers[0],
			)
		}
	}
}

// Tests for tripleFloatBuffer.Remove

var testRemove = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 2, 3, 6, 4, 23234234.2},
				make(map[float64][]int),
			},
		},
		indices: []int{0},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{6, 4, 23234234.2},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 0.2, 3.2, 0.7, 92, 84.4, 643.2, 0.6234, 123.132, 17, 18, 7.6},
				make(map[float64][]int),
			},
		},
		indices: []int{2, 1, 2},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 0.2, 3.2, 17, 18, 7.6},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.23, 41, -6.32, 8, 367.333, -83.224, 8.7225, -424, -35.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		indices: []int{2, 0},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{8, 367.333, -83.224, 6, 7, 9},
				make(map[float64][]int),
			},
		},
	},
}

func TestRemove(t *testing.T) {
	for _, params := range testRemove {
		params.tripleFloatBuffers[0].Remove(params.indices...)
		if !vectorEqual(params.tripleFloatBuffers[0].Buffer, params.resultFloatBuffers[0].Buffer) {
			t.Error(
				"For Removing multiple triples from tripleFloatBuffer",
				"expected", params.resultFloatBuffers[0],
				"got", params.tripleFloatBuffers[0],
			)
		}
	}
}

// Tests for tripleFloatBuffer.UpdateIndex

var testUpdateIndex = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 6.3, 6.3, 2, 3, 6.3, 4, 6.3, 23234234.2},
				make(map[float64][]int),
			},
		},
		floats:     []float64{6.3},
		resultInts: []int{0, 0, 1, 2},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 0.2, 3.2, 0.7, 92, 84.4, 643.2, 0.6234, 123.132, 17, 0.2, 7.6},
				make(map[float64][]int),
			},
		},
		floats:     []float64{0.2},
		resultInts: []int{0, 3},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.23, 41, -6.32, 8, 367.333, -83.224, 8.7225, -424, -35.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		floats:     []float64{-424},
		resultInts: []int{2},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.23, 41, -6.32, 8, 367.333, -83.224, 8.7225, -424, -35.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		floats:     []float64{435},
		resultInts: []int{},
	},
}

func TestUpdateIndex(t *testing.T) {
	for _, params := range testUpdateIndex {
		params.tripleFloatBuffers[0].UpdateIndex()
		if !intVectorEqual(params.tripleFloatBuffers[0].Index[params.floats[0]], params.resultInts) {
			t.Error(
				"For Removing multiple triples from tripleFloatBuffer",
				"expected", params.resultFloatBuffers[0],
				"got", params.tripleFloatBuffers[0],
			)
		}
	}
}

// Tests for tripleFloatBuffer.Collect

var testCollect = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 6.3, 6.3, 2, 3, 6.3, 4, 6.3, 23234234.2},
				make(map[float64][]int),
			},
		},
		collectFunction: func(x, y, z float64) (a, b, c float64) {
			a, b, c = y, z, x
			return
		},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{6.3, 6.3, 1, 3, 6.3, 2, 6.3, 23234234.2, 4},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 0.2, 3.2, 0.7, 92, 84.4, 640.2, 0.6234, 123.132, 17, 0.2, 7.6},
				make(map[float64][]int),
			},
		},
		collectFunction: func(x, y, z float64) (a, b, c float64) {
			a = x + 1
			b = y + 2
			z = 0
			return
		},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{2.5, 2.2, 0, 1.7, 94, 0, 641.2, 2.6234, 0, 18, 2.2, 0},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.23, 41, -6.32, 8, 367.333, -83.224, 8.7225, -424, -35.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		collectFunction: func(x, y, z float64) (a, b, c float64) {
			a = x / y
			b = -y
			c = z * z
			return
		},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{3.005609756, -41, 39.9424, 0.021778604, -367.333, 6926.234176, -0.020571934, 424, 1239.04, 0.857142857, -7, 81},
				make(map[float64][]int),
			},
		},
	},
}

// Tests for tripleFloatBuffer.Each

var testEach = []floatTestParams{
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1, 6.3, 6.3, 2, 3, 6.3, 4, 6.3, 23234234.2},
				make(map[float64][]int),
			},
		},
		eachFunctionFactory: func(tb *tripleFloatBuffer) func(float64, float64, float64) {
			return func(x, y, z float64) {
				tb.Append(y, z, x)
			}
		},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{6.3, 6.3, 1, 3, 6.3, 2, 6.3, 23234234.2, 4},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{1.5, 0.2, 3.2, 0.7, 92, 84.4, 640.2, 0.6234, 123.132, 17, 0.2, 7.6},
				make(map[float64][]int),
			},
		},
		eachFunctionFactory: func(tb *tripleFloatBuffer) func(float64, float64, float64) {
			return func(x, y, z float64) {
				tb.Append(x+1, y+2, 0)
			}
		},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{2.5, 2.2, 0, 1.7, 94, 0, 641.2, 2.6234, 0, 18, 2.2, 0},
				make(map[float64][]int),
			},
		},
	},
	{
		tripleFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{123.23, 41, -6.32, 8, 367.333, -83.224, 8.7225, -424, -35.2, 6, 7, 9},
				make(map[float64][]int),
			},
		},
		eachFunctionFactory: func(tb *tripleFloatBuffer) func(float64, float64, float64) {
			return func(x, y, z float64) {
				tb.Append(x/y, -y, z*z)
			}
		},
		resultFloatBuffers: []tripleFloatBuffer{
			tripleFloatBuffer{
				[]float64{3.005609756, -41, 39.9424, 0.021778604, -367.333, 6926.234176, -0.020571934, 424, 1239.04, 0.857142857, -7, 81},
				make(map[float64][]int),
			},
		},
	},
}

func TestEach(t *testing.T) {
	var dest_buffer tripleFloatBuffer
	for _, params := range testEach {
		dest_buffer = tripleFloatBuffer{
			[]float64{},
			make(map[float64][]int),
		}
		params.tripleFloatBuffers[0].Each(params.eachFunctionFactory(&dest_buffer))
		if !vectorEqual(dest_buffer.Buffer, params.resultFloatBuffers[0].Buffer) {
			t.Error(
				"For iterating over each triple in a tripleFloatBuffer",
				"expected", params.resultFloatBuffers[0],
				"got", dest_buffer,
			)
		}
	}
}
