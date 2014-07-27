package cuboid

import (
	"math"
	"testing"
)

type testParams struct {
	cuboids      []Cuboid
	points       [][3]float64
	resultCuboid Cuboid
	resultPoint  [3]float64
	resultFloat  float64
	resultBool   bool
}

// Tests for Cuboid.volume

var volumeTests = []testParams{
	{
		cuboids:     []Cuboid{*New(1, 1, 1, 2, 2, 2)},
		resultFloat: 1,
	},
	{
		cuboids:     []Cuboid{*New(3, -1, 4.5, 2.1, 2, 2)},
		resultFloat: 6.75,
	},
}

func TestVolume(t *testing.T) {
	for _, params := range volumeTests {
		r := params.cuboids[0].Volume()
		if !floatEqual(r, params.resultFloat) {
			t.Error(
				"For Cuboid", params.cuboids[0],
				"expected volume to be", params.resultFloat,
				"got", r,
			)
		}
	}
}

// Tests for Cuboid.center

var centerTests = []testParams{
	{
		cuboids:     []Cuboid{*New(1, 1, 1, 2, 2, 2)},
		resultPoint: [3]float64{1.5, 1.5, 1.5},
	},
	{
		cuboids:     []Cuboid{*New(3, 1, 4.5, 2.1, 2, 2)},
		resultPoint: [3]float64{2.55, 1.5, 3.25},
	},
}

func TestCenter(t *testing.T) {
	for _, params := range centerTests {
		x, y, z := params.cuboids[0].Center()
		if !vectorEqual([]float64{x, y, z}, params.resultPoint[:]) {
			t.Error(
				"For Cuboid", params.cuboids[0],
				"expected center to be", params.resultPoint,
				"got", []float64{x, y, z},
			)
		}
	}
}

// Tests for Cuboid.contains

var containsTests = []testParams{
	{
		cuboids:    []Cuboid{*New(1, 1, 1, 2, 2, 2)},
		points:     [][3]float64{[3]float64{1.5, 1.5, 1.5}},
		resultBool: true,
	},
	{
		cuboids:    []Cuboid{*New(3, 1, 4.5, 2.1, 2, 2)},
		points:     [][3]float64{[3]float64{3, 1, 4.5}},
		resultBool: true,
	},
	{
		cuboids:    []Cuboid{*New(0, 0, -1, -2, -2, -2)},
		points:     [][3]float64{[3]float64{0, 0, 0}},
		resultBool: false,
	},
}

func TestContains(t *testing.T) {
	for _, params := range containsTests {
		r := params.cuboids[0].Contains(
			params.points[0][0],
			params.points[0][1],
			params.points[0][2],
		)
		if r != params.resultBool {
			if params.resultBool {
				t.Error(
					"Expected Cuboid", params.cuboids[0],
					"to contain point", params.points[0],
				)
			} else {
				t.Error(
					"Expected Cuboid", params.cuboids[0],
					"NOT to contain point", params.points[0],
				)
			}
		}
	}
}

// Tests for Cuboid.Intersects

// Tests for Cuboid.Union

// helpers

func floatEqual(a, b float64) bool {
	FLOAT_EQUALITY_THRESHOLD := 0.0000001
	return math.Abs(a-b) < FLOAT_EQUALITY_THRESHOLD
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
