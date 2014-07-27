package triplebuffer

type TriangleBuffer struct {
	tripleIntBuffer
}

func NewTriangleBuffer() TriangleBuffer {
	return TriangleBuffer{tripleIntBuffer{
		[]int{},
		make(map[int][]int),
	}}
}

// func (vb *TriangleBuffer) area(i int) [3]int {
//   return [0,0,0]
// }

// func (vb *TriangleBuffer) EnsureUniqueness(indices ...int) {
// 	return
// }

func (vb *TriangleBuffer) ApplyMap(mapp map[int]int) {
	for n, i := range vb.Buffer {
		if _, ok := mapp[n]; ok {
			vb.Buffer[i] = mapp[n]
		}
	}
}
