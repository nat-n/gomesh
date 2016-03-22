package mesh

import (
	"errors"
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
	Validate() error
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

	// assert("ReplaceVertex of face called validly",
	// 	// new_vert shouldn't already be references by the face
	// 	f.Verts[0] != new_vert && f.Verts[1] != new_vert && f.Verts[2] != new_vert &&
	// 		// old_vert should be referenced by the face
	// 		(f.Verts[0] == old_vert || f.Verts[1] == old_vert || f.Verts[2] == old_vert))

	// defer func() {
	// 	assert("ReplaceVertex of face succeeded",
	// 		// f.Verts no longer references old_vert
	// 		f.Verts[0] != old_vert && f.Verts[1] != old_vert && f.Verts[2] != old_vert &&
	// 			// f.Verts references new_vert exactly once
	// 			((f.Verts[0] == new_vert && f.Verts[1] != new_vert && f.Verts[2] != new_vert) ||
	// 				(f.Verts[0] != new_vert && f.Verts[1] == new_vert && f.Verts[2] != new_vert) ||
	// 				(f.Verts[0] != new_vert && f.Verts[1] != new_vert && f.Verts[2] == new_vert)))
	// }()

	// var replaced_i int
	// replaced_i = -1

	if f.Vertices[0] == old_vert {
		f.Vertices[0] = new_vert
		// replaced_i = 0
	} else if f.Vertices[1] == old_vert {
		f.Vertices[1] = new_vert
		// replaced_i = 1
	} else if f.Vertices[2] == old_vert {
		f.Vertices[2] = new_vert
		// replaced_i = 2
	} else {
		panic("didn't find old_vert to replace in face")
	}

	// if os.Getenv("verbosity") == "2" {
	// 	fmt.Println("$$1 ReplaceVertex f.Verts: ", f.Verts, "Replaced:", replaced_i)
	// }
}

func (f *Face) ToString() string {
	return "{Face " +
		strconv.Itoa(f.Vertices[0].GetLocationInMesh(f.Mesh)) + " " +
		strconv.Itoa(f.Vertices[1].GetLocationInMesh(f.Mesh)) + " " +
		strconv.Itoa(f.Vertices[2].GetLocationInMesh(f.Mesh)) + "}"
}

func (f *Face) Validate() (err error) {
	if false {
		errors.New("false!")
	}

	// has three different vertices...

	return
}
