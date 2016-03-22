package mesh

import (
	"errors"
	"fmt"
	"github.com/nat-n/geom"
	"strconv"
)

type VertexI interface {
	// Inherited by embedding Vec3I
	GetX() float64
	GetY() float64
	GetZ() float64
	SetX(float64)
	SetY(float64)
	SetZ(float64)
	Clone() geom.Vec3
	Magnitude() float64
	Normalized() geom.Vec3
	Inverse() geom.Vec3
	Add(geom.Vec3I) geom.Vec3
	Sum(...geom.Vec3I) geom.Vec3
	Subtract(geom.Vec3I) geom.Vec3
	Multiply(geom.Vec3I) geom.Vec3
	Divide(geom.Vec3I) geom.Vec3
	AddScalar(float64) geom.Vec3
	SubtractScalar(float64) geom.Vec3
	MultiplyScalar(float64) geom.Vec3
	DivideScalar(float64) geom.Vec3
	Mean(...geom.Vec3I) geom.Vec3
	CrossProd(geom.Vec3I) geom.Vec3
	DotProd(geom.Vec3I) float64
	Angle(geom.Vec3I) float64
	LessThan(geom.Vec3I) bool
	// Vertex methods
	GetMeshLocation() (Mesh, int)
	GetLocationInMesh(Mesh) int
	SetLocationInMesh(Mesh, int)
	ForgetLocationInMeshByName(string)
	OccursInMesh(Mesh) bool
	EachMeshLocation(func(Mesh, int))
	CountOccurances() int
	AddFace(fs ...FaceI) error
	RemoveFace(FaceI) error
	RemoveAllFaces()
	ReferencesFace(FaceI) bool
	EachFace(func(FaceI))
	NeighborCounts() map[VertexI]int
	GetNormal() *geom.Vec3
	SetNormal(*geom.Vec3)
	CalculateNormal()
	ToString() string
	Validate() (err error)
}

type Vertex struct {
	geom.Vec3
	Faces  []FaceI
	Normal *geom.Vec3
	Meshes map[Mesh]int
}

func MakeMeshesMap(m Mesh, i int) map[Mesh]int {
	result := make(map[Mesh]int)
	result[m] = i
	return result
}

func ConvertVertexSliceToVec3ISlice(input []VertexI) []geom.Vec3I {
	result := make([]geom.Vec3I, 0)
	for _, v := range input {
		result = append(result, v)
	}
	return result
}

// Assumes there is only one mesh, panic's otherwise
func (v *Vertex) GetMeshLocation() (Mesh, int) {
	if len(v.Meshes) != 1 {
		fmt.Println(v, v.Meshes)
		for m, _ := range v.Meshes {
			fmt.Println(m.GetName())
		}
		panic("Cannot call GetMeshLocation() on Border Vertex") // FIXME: anti-pattern
	}
	for m, id := range v.Meshes {
		return m, id
	}
	panic("didn't find vert in a mesh...")
	// return nil, 0
}

func (v *Vertex) GetLocationInMesh(m Mesh) int {
	if v.OccursInMesh(m) {
		// fmt.Println("location was set!")
		return v.Meshes[m]
	}
	m2, l2 := v.GetMeshLocation()
	fmt.Println(&m, &m2, m == m2, l2, v.Meshes[m])
	panic("No location is set for vertex in mesh")
}

func (v *Vertex) SetLocationInMesh(m Mesh, i int) {
	v.Meshes[m] = i
}

func (v *Vertex) OccursInMesh(m Mesh) bool {
	_, result := v.Meshes[m]
	return result
}

func (v *Vertex) ForgetLocationInMeshByName(mesh_name string) {
	for m, _ := range v.Meshes {
		if m.GetName() == mesh_name {
			delete(v.Meshes, m)
		}
	}
}

func (v *Vertex) EachMeshLocation(cb func(Mesh, int)) {
	for m, i := range v.Meshes {
		cb(m, i)
	}
}

func (v *Vertex) CountOccurances() int {
	return len(v.Meshes)
}

func (v *Vertex) ReferencesFace(f1 FaceI) bool {
	for _, f2 := range v.Faces {
		if f1 == f2 {
			return true
		}
	}
	return false
}

func (v *Vertex) AddFace(fs ...FaceI) (err error) {
	for _, f := range fs {
		if v.ReferencesFace(f) {
			err = errors.New("Can't add face to vertex twice")
			return
		}
	}
	v.Faces = append(v.Faces, fs...)
	return
}

func (v *Vertex) RemoveFace(f1 FaceI) (err error) {
	for i, f2 := range v.Faces {
		if f1 == f2 {
			v.Faces = append(v.Faces[:i], v.Faces[i+1:]...)
			return
		}
	}
	err = errors.New("Can't remove face from vertex that doesn't reference it")
	return
}

func (v *Vertex) RemoveAllFaces() {
	v.Faces = make([]FaceI, 0)
}

func (v *Vertex) EachFace(cb func(FaceI)) {
	for _, f := range v.Faces {
		cb(f)
	}
}

func (v *Vertex) NeighborCounts() (counts map[VertexI]int) {
	for _, f := range v.Faces {
		f.EachVertex(func(v2 VertexI) {
			if v != v2 {
				counts[v2] += 1
			}
		})
	}
	return
}

func (v *Vertex) GetNormal() *geom.Vec3  { return v.Normal }
func (v *Vertex) SetNormal(n *geom.Vec3) { v.Normal = n }

func (v *Vertex) CalculateNormal() {
	acc := geom.Vec3{0, 0, 0}
	for _, f := range v.Faces {
		t := f.AsTriangle()
		n := t.Normal()
		acc = acc.Add(geom.Vec3I(&n))
	}
	result := acc.DivideScalar(float64(len(v.Faces)))
	v.Normal = &result
}

func (v *Vertex) ToString() string {
	return "{Vertex " +
		strconv.FormatFloat(v.X, 'f', -1, 64) + " " +
		strconv.FormatFloat(v.Y, 'f', -1, 64) + " " +
		strconv.FormatFloat(v.Z, 'f', -1, 64) + "}"
}

func (v *Vertex) Validate() (err error) {
	if false {
		errors.New("false!")
	}
	return
}

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
