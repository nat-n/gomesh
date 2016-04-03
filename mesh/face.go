package mesh

import (
	"github.com/nat-n/geom"
	"strconv"
)

type FaceI interface {
	GetA() VertexI
	GetB() VertexI
	GetC() VertexI
	SetA(v VertexI)
	SetB(v VertexI)
	SetC(v VertexI)
	GetMeshLocation() (Mesh, int)
	SetMeshLocation(Mesh, int)
	AsTriangle() geom.Triangle
	ReferencesVertex(VertexI) bool
	EachVertex(func(VertexI))
	ReplaceVertex(VertexI, VertexI)
	ToString() string
}

type Face struct {
	Vertices [3]VertexI
	Mesh     Mesh
	Index    int
}

func (f *Face) GetA() VertexI { return f.Vertices[0] }
func (f *Face) GetB() VertexI { return f.Vertices[1] }
func (f *Face) GetC() VertexI { return f.Vertices[2] }

func (f *Face) SetA(v VertexI) { f.Vertices[0] = v }
func (f *Face) SetB(v VertexI) { f.Vertices[1] = v }
func (f *Face) SetC(v VertexI) { f.Vertices[2] = v }

func (f *Face) GetMeshLocation() (Mesh, int) { return f.Mesh, f.Index }
func (f *Face) SetMeshLocation(m Mesh, i int) {
	f.Mesh = m
	f.Index = i
}

func (f *Face) AsTriangle() geom.Triangle {
	return geom.Triangle{
		f.Vertices[0],
		f.Vertices[1],
		f.Vertices[2],
	}
}

func (f *Face) ReferencesVertex(v VertexI) bool {
	return f.Vertices[0] == v || f.Vertices[1] == v || f.Vertices[2] == v
}

func (f *Face) EachVertex(cb func(VertexI)) {
	for _, f := range f.Vertices {
		cb(f)
	}
}

func (f *Face) ReplaceVertex(old_vert, new_vert VertexI) {
	if f.Vertices[0] == old_vert {
		f.Vertices[0] = new_vert
	} else if f.Vertices[1] == old_vert {
		f.Vertices[1] = new_vert
	} else if f.Vertices[2] == old_vert {
		f.Vertices[2] = new_vert
	} else {
		panic("didn't find old_vert to replace in face")
	}
}

func (f *Face) ToString() string {
	return "{Face " +
		strconv.Itoa(f.Vertices[0].GetLocationInMesh(f.Mesh)) + " " +
		strconv.Itoa(f.Vertices[1].GetLocationInMesh(f.Mesh)) + " " +
		strconv.Itoa(f.Vertices[2].GetLocationInMesh(f.Mesh)) + "}"
}
