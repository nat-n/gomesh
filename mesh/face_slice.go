package mesh

import (
	"strconv"
	"strings"
)

type FaceSlice struct {
	slice []FaceI
}

func (fs *FaceSlice) Len() int {
	return len(fs.slice)
}

func (fs *FaceSlice) Get(indices ...int) (r []FaceI) {
	for _, i := range indices {
		r = append(r, fs.slice[i])
	}
	return
}

func (fs *FaceSlice) Update(index int, f FaceI) {
	fs.slice[index] = f
}

func (fs *FaceSlice) GetAll() []FaceI {
	r := make([]FaceI, len(fs.slice), len(fs.slice))
	copy(r, fs.slice)
	return r
}

func (fs *FaceSlice) Append(vertices ...FaceI) {
	fs.slice = append(fs.slice, vertices...)
}

func (fs *FaceSlice) Remove(indices ...int) {
	for _, i := range indices {
		fs.slice = append(fs.slice[:i], fs.slice[i+1:]...)
	}
}

func (fs *FaceSlice) Each(cb func(FaceI)) {
	for _, f := range fs.slice {
		cb(f)
	}
}

func (fs *FaceSlice) EachWithIndex(cb func(int, FaceI)) {
	for i, f := range fs.slice {
		cb(i, f)
	}
}

func (fs *FaceSlice) Filter(cb func(FaceI) bool) {
	result := make([]FaceI, 0)
	for _, f := range fs.slice {
		if cb(f) {
			result = append(result, f)
		}
	}
	fs.slice = result
}

func (fs *FaceSlice) IsEmpty() bool {
	return len(fs.slice) == 0
}

func (fs *FaceSlice) ToString() string {
	stringFloats := make([]string, len(fs.slice), len(fs.slice))
	for i := 0; i < len(fs.slice); i++ {
		stringFloats[i] = "(" +
			fs.slice[i].ToString() + " " +
			fs.slice[i].ToString() + " " +
			fs.slice[i].ToString() + ")"
	}
	return "{FaceSlice " + strings.Join(stringFloats, "") + "}"
}

func (fs *FaceSlice) IndicesAsCSV() string {
	stringFloats := make([]string, len(fs.slice), len(fs.slice))
	fs.EachWithIndex(func(i int, f FaceI) {
		m, _ := f.GetMeshLocation()
		stringFloats[i] = "" +
			strconv.Itoa(f.GetA().GetLocationInMesh(m)) + "," +
			strconv.Itoa(f.GetB().GetLocationInMesh(m)) + "," +
			strconv.Itoa(f.GetC().GetLocationInMesh(m))
	})
	return strings.Join(stringFloats, ",")
}
