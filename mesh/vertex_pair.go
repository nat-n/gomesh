package mesh

type VertexPair struct {
	V1 VertexI
	V2 VertexI
}

func MakeVertexPair(v1, v2 VertexI) VertexPair {
	if v1.LessThan(v2) {
		return VertexPair{v1, v2}
	} else {
		return VertexPair{v2, v1}
	}
}

func (vp *VertexPair) LessThan(vp2 *VertexPair) bool {
	return (vp.V1.LessThan(vp2.V1) ||
		(!vp2.V1.LessThan(vp.V1) && vp.V2.LessThan(vp2.V2)))
}

type SortableVertexPairs []VertexPair

func (v SortableVertexPairs) Len() int           { return len(v) }
func (v SortableVertexPairs) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v SortableVertexPairs) Less(i, j int) bool { return v[i].LessThan(&v[j]) }
