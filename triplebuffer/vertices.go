package triplebuffer

import (
	"math"
)

type VertexBuffer struct {
	tripleFloatBuffer
}

func NewVertexBuffer() VertexBuffer {
	return VertexBuffer{tripleFloatBuffer{
		[]float64{},
		make(map[float64][]int),
	}}
}

// TODO:
// Accepts three vertex indices and calculates the normal vector of the triangle
//  defined by them.
// func (vb *VertexBuffer) triangleNormal(a, b, c int) [3]float64 {
// 	return [3]float64{0, 0, 0}
// }

// Accepts the indices of two vertices and calculates the distance between them.
func (vb *VertexBuffer) DistanceBetween(i1, i2 int) float64 {
	return math.Sqrt(math.Pow(vb.Buffer[i1*3]-vb.Buffer[i2*3], 2) +
		math.Pow(vb.Buffer[i1*3+1]-vb.Buffer[i2*3+1], 2) +
		math.Pow(vb.Buffer[i1*3+2]-vb.Buffer[i2*3+2], 2))
}

func (vb *VertexBuffer) DistanceFrom(i1 int, x, y, z float64) float64 {
	return math.Sqrt(math.Pow(vb.Buffer[i1*3]-x, 2) +
		math.Pow(vb.Buffer[i1*3+1]-y, 2) +
		math.Pow(vb.Buffer[i1*3+2]-z, 2))
}

// TODO:
// Finds the shortest distance from the vertex with the given index, to the line
//  of points [x1,y1,z1] and [x2,y2,z2].
// func (vb *VertexBuffer) distanceFromLine(i int,
// 										 x1 float64, y1 float64, z1 float64,
// 										 x2 float64, y2 float64, z2 float64)
// 						(float64) {

//     other = [ x1 - vb.Buffer[i*3],
//               y1 - vb.Buffer[i*3+1],
//               z1 - vb.Buffer[i*3+2]]
//     seg = { x1 - x2, y1 - y2, z1 - z2 }

//     Math.sqrt( (seg[1]*other[2]-seg[2]*other[1])**2 +
//                (seg[2]*other[0]-seg[0]*other[2])**2 +
//                (seg[0]*other[1]-seg[1]*other[0])**2 ) / Math.sqrt( seg[0]**2 +
//                                                                    seg[1]**2 +
//                                                                    seg[2]**2 )
// }

// func (vb *VertexBuffer) distanceFromLineSegment(i1 int,
// 											    x1 float64,
// 											    y1 float64,
// 											    z1 float64,
// 											    x2 float64,
// 											    y2 float64,
// 											    z2 float64)
// 						(float64) {
//   return 0
// }

// Accepts set of vertex indices and calculates the average vertex.
func (vb *VertexBuffer) Average(indices ...int) [3]float64 {
	sum := [3]float64{}
	for _, i := range indices {
		sum[0] += vb.Buffer[i*3]
		sum[1] += vb.Buffer[i*3+1]
		sum[2] += vb.Buffer[i*3+2]
	}
	sum[0] /= float64(len(indices))
	sum[1] /= float64(len(indices))
	sum[2] /= float64(len(indices))
	return sum
}

// Accepts a triangle in the form of three vertex indices.
// If the triangle is right angled then it returns a slice containing the two
//  vertex indices defining the line segment of the hypotenuse and true.
// If the triangle is not right angled then it returns a slice containing two
//  zeros and false.
func (vb *VertexBuffer) Hypotenuse(a, b, c int) ([2]int, bool) {
	FLOAT_EQUALITY_THRESHOLD := 0.0000001

	sqrA := (math.Pow(vb.Buffer[a*3]-vb.Buffer[b*3], 2) +
		math.Pow(vb.Buffer[a*3+1]-vb.Buffer[b*3+1], 2) +
		math.Pow(vb.Buffer[a*3+2]-vb.Buffer[b*3+2], 2))
	sqrB := (math.Pow(vb.Buffer[b*3]-vb.Buffer[c*3], 2) +
		math.Pow(vb.Buffer[b*3+1]-vb.Buffer[c*3+1], 2) +
		math.Pow(vb.Buffer[b*3+2]-vb.Buffer[c*3+2], 2))
	sqrC := (math.Pow(vb.Buffer[c*3]-vb.Buffer[a*3], 2) +
		math.Pow(vb.Buffer[c*3+1]-vb.Buffer[a*3+1], 2) +
		math.Pow(vb.Buffer[c*3+2]-vb.Buffer[a*3+2], 2))

	if math.Abs(sqrA-sqrB+sqrC) < FLOAT_EQUALITY_THRESHOLD {
		return [2]int{a, b}, true
	} else if math.Abs(sqrB-sqrC+sqrA) < FLOAT_EQUALITY_THRESHOLD {
		return [2]int{b, c}, true
	} else if math.Abs(sqrC-sqrA+sqrB) < FLOAT_EQUALITY_THRESHOLD {
		return [2]int{c, a}, true
	} else {
		return [2]int{}, false
	}
}
