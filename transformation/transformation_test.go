package transformation

import (
	"math"
	"testing"
)

type testParams struct {
	transfromations []Transformation
	vector          [3]float64
	cell            []int
	col             int
	row             int
	product         Transformation
	sum             Transformation
	translation     []float64
	scale           float64
	scaleDims       float64
	rotation        []float64
	resultVec       []float64
	resultTrans     Transformation
	resultFloat     float64
}

// Transformation.cell

var cellTests = []testParams{
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		cell:        []int{1, 0},
		resultFloat: 1,
	},
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		cell:        []int{2, 2},
		resultFloat: 10,
	},
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		cell:        []int{0, 3},
		resultFloat: 12,
	},
}

func TestCell(t *testing.T) {
	for _, params := range cellTests {
		v := params.transfromations[0].Cell(params.cell[0], params.cell[1])
		if v != params.resultFloat {
			t.Error(
				"For cell", params.cell,
				"expected", params.resultFloat,
				"got", v,
			)
		}
	}
}

// Transformation.col

var colTests = []testParams{
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		col:       0,
		resultVec: []float64{0, 4, 8, 12},
	},
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		col:       1,
		resultVec: []float64{1, 5, 9, 13},
	},
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		col:       3,
		resultVec: []float64{3, 7, 11, 15},
	},
}

func TestCol(t *testing.T) {
	for _, params := range colTests {
		r := params.transfromations[0].Col(params.col)
		a := [4]float64{}
		copy(a[:], params.resultVec)
		if r != a {
			t.Error(
				"For col", params.col,
				"expected", a,
				"got", r,
			)
		}
	}
}

// Transformation.row

var rowTests = []testParams{
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		row:       0,
		resultVec: []float64{0, 1, 2, 3},
	},
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		row:       1,
		resultVec: []float64{4, 5, 6, 7},
	},
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
		},
		row:       3,
		resultVec: []float64{12, 13, 14, 15},
	},
}

func TestRow(t *testing.T) {
	for _, params := range rowTests {
		r := params.transfromations[0].Row(params.row)
		a := [4]float64{}
		copy(a[:], params.resultVec)
		if r != a {
			t.Error(
				"For row", params.row,
				"expected", a,
				"got", r,
			)
		}
	}
}

// Transformation.multiply

var multiplyTests = []testParams{
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
			Transformation{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
		},
		resultTrans: Transformation{
			0, 1, 2, 3,
			4, 5, 6, 7,
			8, 9, 10, 11,
			12, 13, 14, 15,
		},
	},
	{
		transfromations: []Transformation{
			Transformation{
				1, 2, 3, 4,
				3, 2, 1, 5,
				2, 1, 3, 6,
				7, 1, 3, 2,
			},
			Transformation{
				3, 4, 5, 6,
				6, 5, 4, 4,
				4, 6, 5, 9,
				4, 6, 5, 9,
			},
			Transformation{
				7, 6, 5, 4,
				1, 2, 3, 4,
				8, 3, 5, 4,
				4, 6, 5, 9,
			},
		},
		resultTrans: Transformation{
			1049, 976, 1008, 1281,
			1117, 1025, 1064, 1344,
			1263, 1181, 1221, 1569,
			1268, 1146, 1199, 1515,
		},
	},
}

func TestMultiply(t *testing.T) {
	for _, params := range multiplyTests {
		v := params.transfromations[0].Multiply(params.transfromations[1:]...)
		if !transformationEqual(v, params.resultTrans) {
			t.Error(
				"For matrix multiplication with\n", params.transfromations[0], "*", params.transfromations[1], "\n",
				"expected", params.resultTrans,
				"got", v,
			)
		}
	}
}

// Transformation.add

var addTests = []testParams{
	{
		transfromations: []Transformation{
			Transformation{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
			Transformation{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
		},
		resultTrans: Transformation{
			1, 1, 2, 3,
			4, 6, 6, 7,
			8, 9, 11, 11,
			12, 13, 14, 16,
		},
	},
	{
		transfromations: []Transformation{
			Transformation{
				1, 2, 3, 4,
				3, 2, 1, 5,
				2, 1, 3, 6,
				7, 1, 3, 2,
			},
			Transformation{
				3, 4, 5, 6,
				6, 5, 4, 4,
				4, 6, 5, 9,
				4, 6, 5, 9,
			},
			Transformation{
				7, 6, 5, 4,
				1, 2, 3, 4,
				8, 3, 5, 4,
				4, 6, 5, 9,
			},
		},
		resultTrans: Transformation{
			11, 12, 13, 14,
			10, 9, 8, 13,
			14, 10, 13, 19,
			15, 13, 13, 20,
		},
	},
}

func TestAdd(t *testing.T) {
	for _, params := range addTests {
		v := params.transfromations[0].Add(params.transfromations[1:]...)
		if !transformationEqual(v, params.resultTrans) {
			t.Error(
				"For matrix multiplication with\n", params.transfromations[0], "*", params.transfromations[1], "\n",
				"expected", params.resultTrans,
				"got", v,
			)
		}
	}
}

// Transformation.apply

var applyTests = []testParams{
	{
		transfromations: []Transformation{
			Transformation{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
		},
		vector:    [3]float64{1, 2, 3},
		resultVec: []float64{1, 2, 3},
	},
	{
		transfromations: []Transformation{
			Transformation{
				-1, 0, 0, 0,
				0, 10, 0, 0,
				0, 0, 3, 0,
				0, 0, 0, 1,
			},
		},
		vector:    [3]float64{1, 2, 3},
		resultVec: []float64{-1, 20, 9},
	},
	{
		transfromations: []Transformation{
			Transformation{
				-1, 0, 0, 4.1,
				0, 10, 0, 5,
				0, 0, 3, -6,
				0, 0, 0, 1,
			},
		},
		vector:    [3]float64{1, 2, 3},
		resultVec: []float64{3.1, 25, 3},
	},
}

func TestApply(t *testing.T) {
	for _, params := range applyTests {
		s := params.vector[:]
		params.transfromations[0].Apply(s)
		if !vectorEqual(s, params.resultVec) {
			t.Error(
				"For application of transform with", params.transfromations[0], ",",
				"expected", params.resultVec, ",",
				"got", params.vector,
			)
		}
	}
}

// transformation.Translation

// transformation.Scale

// transformation.ScaleDimensions

// transformation.Rotation

// helper functions

func transformationEqual(a, b Transformation) bool {
	FLOAT_EQUALITY_THRESHOLD := 0.0000001
	for i, _ := range a {
		if math.Abs(a[i]-b[i]) > FLOAT_EQUALITY_THRESHOLD {
			return false
		}
	}
	return true
}

func vectorEqual(a, b []float64) bool {
	FLOAT_EQUALITY_THRESHOLD := 0.0000001
	for i, _ := range a {
		if math.Abs(a[i]-b[i]) > FLOAT_EQUALITY_THRESHOLD {
			return false
		}
	}
	return true
}
