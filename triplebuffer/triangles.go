package triplebuffer

import "math"

type TriangleBuffer struct {
	tripleIntBuffer
}

func NewTriangleBuffer() TriangleBuffer {
	return TriangleBuffer{tripleIntBuffer{
		[]int{},
		make(map[int][]int),
	}}
}

// func (tb *TriangleBuffer) area(i int) [3]int {
//   return [0,0,0]
// }

// func (tb *TriangleBuffer) EnsureUniqueness(indices ...int) {
// 	return
// }

func (tb *TriangleBuffer) ApplyMap(mapp map[int]int) {
	for n, i := range tb.Buffer {
		if _, ok := mapp[n]; ok {
			tb.Buffer[i] = mapp[n]
		}
	}
}

// Assumes an up to date index exists
func (tb *TriangleBuffer) NeighboursOf(i int) (result []int) {
	verts := tb.Get(tb.Index[i]...)
	set := make(map[int]bool)
	for _, vi := range verts {
		set[vi] = true
	}
	for vi, _ := range set {
		if vi != i {
			result = append(result, vi)
		}
	}
	return
}

// calculate a normal vector for triangle with vertices a, b and c using the
// "right-hand rule"
func Normal(a, b, c []float64) (x, y, z float64) {
	// Derive the vectors of two sides of the triangle
	v1 := [3]float64{b[0] - a[0], b[1] - a[1], b[2] - a[2]}
	v2 := [3]float64{c[0] - a[0], c[1] - a[1], c[2] - a[2]}

	// Calculate the cross product
	cpx := v1[1]*v2[2] - v1[2]*v2[1]
	cpy := v1[2]*v2[0] - v1[0]*v2[2]
	cpz := v1[0]*v2[1] - v1[1]*v2[0]

	// Normalize the cross product to arrive at the normal vector
	l := math.Sqrt(cpx*cpx + cpy*cpy + cpz*cpz)
	x = cpx / l
	y = cpy / l
	z = cpz / l
	return
}
