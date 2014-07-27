package triplebuffer

import "math"

type VectorBuffer struct {
	tripleFloatBuffer
}

func NewVectorBuffer() VectorBuffer {
	return VectorBuffer{tripleFloatBuffer{
		[]float64{},
		make(map[float64][]int),
	}}
}

// Calculates the length of vector i (from the origin)
func (vb *VectorBuffer) length(i int) float64 {
	return math.Sqrt(
		math.Pow(vb.Buffer[i*3], 2) +
			math.Pow(vb.Buffer[i*3+1], 2) +
			math.Pow(vb.Buffer[i*3+2], 2))
}

// Update vector i to have length 1
func (vb *VectorBuffer) normalize(i int) {
	l := math.Sqrt(
		math.Pow(vb.Buffer[i*3], 2) +
			math.Pow(vb.Buffer[i*3+1], 2) +
			math.Pow(vb.Buffer[i*3+2], 2))
	vb.Buffer[i*3] /= l
	vb.Buffer[i*3+1] /= l
	vb.Buffer[i*3+2] /= l
}

func (vb *VectorBuffer) crossProd(i int, x, y, z float64) (a, b, c float64) {
	index := i * 3
	a = vb.Buffer[index+1]*z - vb.Buffer[index+2]*y
	b = vb.Buffer[index+2]*x - vb.Buffer[index]*z
	c = vb.Buffer[index]*y - vb.Buffer[index+1]*x
	return
}

func (vb *VectorBuffer) dotProd(i int, x, y, z float64) float64 {
	index := i * 3
	return vb.Buffer[index]*x + vb.Buffer[index+1]*y + vb.Buffer[index+2]*z
}

// func (vb *VectorBuffer) averageNormal(indices ...int) float64 {
// 	return 0
// }
