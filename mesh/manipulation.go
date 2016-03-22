package mesh

import (
	"container/list"
	"fmt"
	"github.com/nat-n/geom"
	"math"
	"sort"
)

import cb "github.com/nat-n/gomesh/cuboid"
import tr "github.com/nat-n/gomesh/transformation"

// Applies the given transformation to every vertex.
func (m *Mesh) Transform(t tr.Transformation) {
	t.ApplyToVec3(ConvertVertexSliceToVec3ISlice(m.Vertices.GetAll())...)
}

// Applies the given transformation to each vertex in indices.
func (m *Mesh) TransformSubset(indices []int, t tr.Transformation) {
	t.ApplyToVec3(ConvertVertexSliceToVec3ISlice(m.Vertices.Get(indices...))...)
}

func (m *Mesh) BoundingBox() *cb.Cuboid {
	var minX, maxX, minY, maxY, minZ, maxZ float64
	minX = math.Inf(1)
	minY = math.Inf(1)
	minZ = math.Inf(1)
	maxX = math.Inf(-1)
	maxY = math.Inf(-1)
	maxZ = math.Inf(-1)
	for i := 0; i < m.Vertices.Len(); i++ {
		v := m.Vertices.Get(i)[0]
		minX = math.Min(minX, v.GetX())
		minY = math.Min(minY, v.GetY())
		minZ = math.Min(minZ, v.GetZ())
		maxX = math.Max(maxX, v.GetX())
		maxY = math.Max(maxY, v.GetY())
		maxZ = math.Max(maxZ, v.GetZ())
	}
	return cb.New(minX, minY, minZ, maxX, maxY, maxZ)
}

func (m *Mesh) SubsetBoundingBox(subset interface{}) *cb.Cuboid {
	var subset_indices []int

	if verts, ok := subset.([]VertexI); ok {
		subset_indices = make([]int, len(verts), len(verts))
		for i, vert := range verts {
			subset_indices[i] = m.GetIndexOf(vert)
		}
	} else {
		subset_indices = subset.([]int)
	}

	var minX, maxX, minY, maxY, minZ, maxZ float64
	minX = math.Inf(1)
	minY = math.Inf(1)
	minZ = math.Inf(1)
	maxX = math.Inf(-1)
	maxY = math.Inf(-1)
	maxZ = math.Inf(-1)
	for _, i := range subset_indices {
		v := m.Vertices.Get(i)[0]
		minX = math.Min(minX, v.GetX())
		minY = math.Min(minY, v.GetY())
		minZ = math.Min(minZ, v.GetZ())
		maxX = math.Max(maxX, v.GetX())
		maxY = math.Max(maxY, v.GetY())
		maxZ = math.Max(maxZ, v.GetZ())
	}
	return cb.New(minX, minY, minZ, maxX, maxY, maxZ)
}

// Constructs a new Mesh from a cuboid.
func NewFromCuboid(c cb.Cuboid) (m *Mesh) {
	m = New("Cuboid")
	normalComponent := 1.0 / math.Sqrt(1)
	verts := []VertexI{
		VertexI(&Vertex{
			Vec3:   geom.Vec3{c.OriginX, c.OriginY, c.OriginZ},
			Normal: &geom.Vec3{-normalComponent, -normalComponent, -normalComponent},
			Meshes: MakeMeshesMap(*m, 0),
		}),
		&Vertex{
			Vec3:   geom.Vec3{c.TerminusX, c.OriginY, c.OriginZ},
			Normal: &geom.Vec3{normalComponent, -normalComponent, -normalComponent},
			Meshes: MakeMeshesMap(*m, 1),
		},
		&Vertex{
			Vec3:   geom.Vec3{c.TerminusX, c.OriginY, c.TerminusZ},
			Normal: &geom.Vec3{normalComponent, -normalComponent, normalComponent},
			Meshes: MakeMeshesMap(*m, 2),
		},
		&Vertex{
			Vec3:   geom.Vec3{c.OriginX, c.OriginY, c.TerminusZ},
			Normal: &geom.Vec3{-normalComponent, -normalComponent, normalComponent},
			Meshes: MakeMeshesMap(*m, 3),
		},
		&Vertex{
			Vec3:   geom.Vec3{c.OriginX, c.TerminusY, c.OriginZ},
			Normal: &geom.Vec3{-normalComponent, normalComponent, -normalComponent},
			Meshes: MakeMeshesMap(*m, 4),
		},
		&Vertex{
			Vec3:   geom.Vec3{c.TerminusX, c.TerminusY, c.OriginZ},
			Normal: &geom.Vec3{normalComponent, normalComponent, -normalComponent},
			Meshes: MakeMeshesMap(*m, 5),
		},
		&Vertex{
			Vec3:   geom.Vec3{c.TerminusX, c.TerminusY, c.TerminusZ},
			Normal: &geom.Vec3{normalComponent, normalComponent, normalComponent},
			Meshes: MakeMeshesMap(*m, 6),
		},
		&Vertex{
			Vec3:   geom.Vec3{c.OriginX, c.TerminusY, c.TerminusZ},
			Normal: &geom.Vec3{-normalComponent, normalComponent, normalComponent},
			Meshes: MakeMeshesMap(*m, 7),
		},
	}
	faces := []FaceI{
		&Face{
			Vertices: [3]VertexI{verts[0], verts[2], verts[1]},
			Mesh:     *m,
			Index:    0,
		},
		&Face{
			Vertices: [3]VertexI{verts[0], verts[3], verts[2]},
			Mesh:     *m,
			Index:    1,
		},
		&Face{
			Vertices: [3]VertexI{verts[0], verts[5], verts[1]},
			Mesh:     *m,
			Index:    2,
		},
		&Face{
			Vertices: [3]VertexI{verts[0], verts[4], verts[5]},
			Mesh:     *m,
			Index:    3,
		},
		&Face{
			Vertices: [3]VertexI{verts[1], verts[6], verts[2]},
			Mesh:     *m,
			Index:    4,
		},
		&Face{
			Vertices: [3]VertexI{verts[1], verts[5], verts[6]},
			Mesh:     *m,
			Index:    5,
		},
		&Face{
			Vertices: [3]VertexI{verts[2], verts[7], verts[3]},
			Mesh:     *m,
			Index:    6,
		},
		&Face{
			Vertices: [3]VertexI{verts[2], verts[6], verts[7]},
			Mesh:     *m,
			Index:    7,
		},
		&Face{
			Vertices: [3]VertexI{verts[3], verts[4], verts[0]},
			Mesh:     *m,
			Index:    8,
		},
		&Face{
			Vertices: [3]VertexI{verts[3], verts[7], verts[4]},
			Mesh:     *m,
			Index:    9,
		},
		&Face{
			Vertices: [3]VertexI{verts[5], verts[7], verts[6]},
			Mesh:     *m,
			Index:    10,
		},
		&Face{
			Vertices: [3]VertexI{verts[5], verts[4], verts[7]},
			Mesh:     *m,
			Index:    11,
		},
	}
	for _, f := range faces {
		f.EachVertex(func(v VertexI) {
			v.AddFace(f)
		})
	}
	m.Vertices.Append(verts...)
	m.Faces.Append(faces...)
	return
}

// Identifies border vertices and returns an array of arrays representing closed
// loops of border vertices.
// Border vertices are identified as including a face which includes an edge
// which is only included in that one face.
func (m *Mesh) IdentifyBoundaries() (boundaries [][]VertexI) {
	edge_counts := make(map[VertexPair]int)
	boundary_edges_slice := make(sortableVertexPairs, 0)
	boundary_edges := list.New()

	m.Faces.Each(func(f FaceI) {
		edge_counts[MakeVertexPair(f.GetA(), f.GetB())]++
		edge_counts[MakeVertexPair(f.GetB(), f.GetC())]++
		edge_counts[MakeVertexPair(f.GetC(), f.GetA())]++
	})

	for edge, count := range edge_counts {
		if count == 1 {
			boundary_edges_slice = append(boundary_edges_slice, edge)
		}
	}

	// The purpose of the intermediate boundary_edges_slice is so the following
	// intermediate boundary_edges list will be sorted so that this function can
	// be idempotent.
	sort.Sort(boundary_edges_slice)
	for _, boundary_edge := range boundary_edges_slice {
		boundary_edges.PushBack(boundary_edge)
	}

	// Transform boundary_edges into one or more closed loops of connected
	// vertices using partials for termporary storage
	partials := list.New()
	for boundary_edges.Len() > 0 {
		latest_partial := make([]VertexI, 0)
		seed_edge := boundary_edges.Front().Value.(VertexPair)
		_ = boundary_edges.Remove(boundary_edges.Front())
		latest_partial = append(latest_partial, seed_edge.V1, seed_edge.V2)

		for true {
			head := latest_partial[0]
			tail := latest_partial[len(latest_partial)-1]
			if head == tail {
				// border is complete
				break
			}
			found_next_node := false
			for node := boundary_edges.Front(); node != nil; node = node.Next() {
				// fmt.Println("node", boundary_edges.Len())
				if node.Value.(VertexPair).V1 == tail {
					latest_partial = append(latest_partial, node.Value.(VertexPair).V2)
					boundary_edges.Remove(node)
					found_next_node = true
					break
				} else if node.Value.(VertexPair).V2 == tail {
					latest_partial = append(latest_partial, node.Value.(VertexPair).V1)
					boundary_edges.Remove(node)
					found_next_node = true
					break
				}
			}
			if !found_next_node {
				// TODO: make better
				fmt.Println(m.GetName())
				fmt.Println("got stuck, with remaining", boundary_edges.Len())
				panic(":(")
			}
		}
		partials.PushBack(latest_partial[:len(latest_partial)-1])
	}

	// Copy completed boundaries from partials over into boundaries
	boundaries = make([][]VertexI, 0, partials.Len())
	for el := partials.Front(); el != nil; el = el.Next() {
		boundary_len := len(el.Value.([]VertexI))
		complete_boundary := make([]VertexI, boundary_len, boundary_len)
		for i, boundary_vert := range el.Value.([]VertexI) {
			complete_boundary[i] = boundary_vert
		}
		boundaries = append(boundaries, complete_boundary)
	}

	return
}

// for sorting the slice of edges... TODO: tidy this up somewhere

type sortableVertexPairs []VertexPair

func (v sortableVertexPairs) Len() int           { return len(v) }
func (v sortableVertexPairs) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v sortableVertexPairs) Less(i, j int) bool { return v[i].LessThan(&v[j]) }
