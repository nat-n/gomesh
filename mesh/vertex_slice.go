package mesh

import (
	"github.com/nat-n/geom"
	"strconv"
	"strings"
)

type VertexSlice struct {
	slice []VertexI
}

func (vs *VertexSlice) Len() int {
	return len(vs.slice)
}

func (vs *VertexSlice) Get(indices ...int) (r []VertexI) {
	for _, i := range indices {
		r = append(r, vs.slice[i])
	}
	return
}

func (vs *VertexSlice) Update(i int, v VertexI) {
	vs.slice[i] = v
}

func (vs *VertexSlice) GetAll() []VertexI {
	r := make([]VertexI, len(vs.slice), len(vs.slice))
	copy(r, vs.slice)
	return r
}

func (vs *VertexSlice) Append(vertices ...VertexI) {
	vs.slice = append(vs.slice, vertices...)
}

func (vs *VertexSlice) Remove(indices ...int) {
	for _, i := range indices {
		vs.slice = append(vs.slice[:i], vs.slice[i+1:]...)
	}
}

func (vs *VertexSlice) Filter(cb func(VertexI) bool) {
	result := make([]VertexI, 0)
	for _, v := range vs.slice {
		if cb(v) {
			result = append(result, v)
		}
	}
	vs.slice = result
}

func (vs *VertexSlice) Each(cb func(VertexI)) {
	for _, v := range vs.slice {
		cb(v)
	}
}

func (vs *VertexSlice) EachWithIndex(cb func(int, VertexI)) {
	for i, v := range vs.slice {
		cb(i, v)
	}
}

func (vs *VertexSlice) IsEmpty() bool {
	return len(vs.slice) == 0
}

func (vs *VertexSlice) Average() (avg geom.Vec3) {
	length := vs.Len()
	for i := 0; i < length; i++ {
		v := vs.slice[i]
		avg.X += v.GetX()
		avg.Y += v.GetY()
		avg.Z += v.GetZ()
	}
	avg.X /= float64(length)
	avg.Y /= float64(length)
	avg.Z /= float64(length)
	return
}

func (vs *VertexSlice) ToString() string {
	stringFloats := make([]string, len(vs.slice), len(vs.slice))
	for i := 0; i < len(vs.slice); i++ {
		stringFloats[i] = "(" +
			strconv.FormatFloat(vs.slice[i].GetX(), 'f', -1, 64) + " " +
			strconv.FormatFloat(vs.slice[i].GetY(), 'f', -1, 64) + " " +
			strconv.FormatFloat(vs.slice[i].GetZ(), 'f', -1, 64) + ")"
	}
	return "{VertexSlice " + strings.Join(stringFloats, "") + "}"
}

func (vs *VertexSlice) PositionsAsCSV() string {
	stringFloats := make([]string, len(vs.slice), len(vs.slice))
	for i := 0; i < len(vs.slice); i++ {
		stringFloats[i] = "" +
			strconv.FormatFloat(vs.slice[i].GetX(), 'f', -1, 64) + "," +
			strconv.FormatFloat(vs.slice[i].GetY(), 'f', -1, 64) + "," +
			strconv.FormatFloat(vs.slice[i].GetZ(), 'f', -1, 64)
	}
	return strings.Join(stringFloats, ",")
}

func (vs *VertexSlice) NormalsAsCSV() string {
	stringFloats := make([]string, len(vs.slice), len(vs.slice))
	for i := 0; i < len(vs.slice); i++ {
		normal := vs.slice[i].GetNormal()
		if normal == nil {
			return ""
		}
		stringFloats[i] = "" +
			strconv.FormatFloat(normal.GetX(), 'f', -1, 64) + "," +
			strconv.FormatFloat(normal.GetY(), 'f', -1, 64) + "," +
			strconv.FormatFloat(normal.GetZ(), 'f', -1, 64)
	}
	return strings.Join(stringFloats, ",")
}
