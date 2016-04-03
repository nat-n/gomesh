package mesh

import (
	"errors"
	"github.com/nat-n/geom"
)

type MeshI interface {
	GetName() string
	GetVertices() VertexCollection
	GetFaces() FaceCollection
	ReindexVerticesAndFaces()
}

type Mesh struct {
	Name     string
	Vertices VertexCollection
	Faces    FaceCollection
}

// Constructor
func New(name string) *Mesh {
	return &Mesh{
		Name:     name,
		Vertices: &VertexSlice{make([]VertexI, 0)},
		Faces:    &FaceSlice{make([]FaceI, 0)},
	}
}

type VertexCollection interface {
	Len() int
	Get(...int) []VertexI
	Update(int, VertexI)
	GetAll() []VertexI
	Append(...VertexI)
	Remove(...int)
	Filter(func(VertexI) bool)
	Each(func(VertexI))
	EachWithIndex(func(int, VertexI))
	IsEmpty() bool
	Average() geom.Vec3
	ToString() string
	PositionsAsCSV() string
	NormalsAsCSV() string
}

type FaceCollection interface {
	Len() int
	Get(...int) []FaceI
	Update(int, FaceI)
	GetAll() []FaceI
	Append(...FaceI)
	Remove(...int)
	Filter(func(FaceI) bool)
	Each(func(FaceI))
	EachWithIndex(func(int, FaceI))
	IsEmpty() bool
	ToString() string
	IndicesAsCSV() string
}

func (m *Mesh) GetName() string {
	return m.Name
}

func (m *Mesh) GetVertices() VertexCollection {
	return m.Vertices
}

func (m *Mesh) GetFaces() FaceCollection {
	return m.Faces
}

func (m *Mesh) ReindexVerticesAndFaces() {
	m.Vertices.EachWithIndex(func(i int, v VertexI) {
		v.ForgetLocationInMeshByName(m.GetName())
		v.SetLocationInMesh(*m, i)
	})
	m.Faces.EachWithIndex(func(i int, f FaceI) { f.SetMeshLocation(*m, i) })
}

// Accepts two vertices with identical locations and moves all faces from the
// secondaries to the primary
// This method assumes vertex Indices are accurate
func MergeSharedVertices(vprime VertexI, vsecs ...VertexI) (err error) {
	for _, vsec := range vsecs {
		if vprime.GetX() != vsec.GetX() ||
			vprime.GetY() != vsec.GetY() ||
			vprime.GetZ() != vsec.GetZ() {
			err = errors.New("Cannot merge vertices with different locations: " +
				vprime.ToString() + " " + vsec.ToString())
			return
		}
	}

	for _, vsec := range vsecs {
		// connect faces from vsec to vprime
		vsec.EachFace(func(f FaceI) {
			f.ReplaceVertex(vsec, vprime)
			err = vprime.AddFace(f)
			if err != nil {
				return
			}
		})
		// clear face references from vsec
		vsec.RemoveAllFaces()

		// update vertex/mesh relationship to replace vsec with vprime
		vsec_mesh, vsec_i := vsec.GetMeshLocation()
		if vsec_mesh.GetVertices().Get(vsec_i)[0] != vsec {
			panic("Method assumption violated: vertex index inaccurate")
		}
		vsec_mesh.GetVertices().Update(vsec_i, vprime)
		vprime.SetLocationInMesh(vsec_mesh, vsec_i)
	}
	return
}
