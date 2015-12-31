package mesh

import (
	"container/list"
	"math"
	"sort"
)

import cb "github.com/nat-n/gomesh/cuboid"
import tr "github.com/nat-n/gomesh/transformation"
import tb "github.com/nat-n/gomesh/triplebuffer"

type Mesh struct {
	Name  string
	Verts tb.VertexBuffer
	Norms tb.VectorBuffer
	Faces tb.TriangleBuffer
}

// Constructor
func New(name string) *Mesh {
	return &Mesh{
		Name:  name,
		Verts: tb.NewVertexBuffer(),
		Norms: tb.NewVectorBuffer(),
		Faces: tb.NewTriangleBuffer(),
	}
}

// Applies the given transformation to every vertex.
func (m *Mesh) Transform(t tr.Transformation) {
	for i := 0; i < m.Verts.Len(); i++ {
		t.Apply(m.Verts.Buffer[i*3 : i*3+3])
	}
}

// Applies the given transformation to each vertex in indices.
func (m *Mesh) TransformSubset(indices []int, t tr.Transformation) {
	for _, i := range indices {
		t.Apply(m.Verts.Buffer[i*3 : i*3+3])
	}
}

func (m *Mesh) BoundingBox() *cb.Cuboid {
	var minX, maxX, minY, maxY, minZ, maxZ float64
	var index int
	minX = math.Inf(1)
	minY = math.Inf(1)
	minZ = math.Inf(1)
	maxX = math.Inf(-1)
	maxY = math.Inf(-1)
	maxZ = math.Inf(-1)
	for i := 0; i < m.Verts.Len(); i++ {
		index = i * 3
		minX = math.Min(minX, m.Verts.Buffer[index])
		minY = math.Min(minY, m.Verts.Buffer[index+1])
		minZ = math.Min(minZ, m.Verts.Buffer[index+2])
		maxX = math.Max(maxX, m.Verts.Buffer[index])
		maxY = math.Max(maxY, m.Verts.Buffer[index+1])
		maxZ = math.Max(maxZ, m.Verts.Buffer[index+2])
	}
	return cb.New(minX, minY, minZ, maxX, maxY, maxZ)
}

func (m *Mesh) SubsetBoundingBox(subset_indices []int) *cb.Cuboid {
	var minX, maxX, minY, maxY, minZ, maxZ float64
	var index int
	minX = math.Inf(1)
	minY = math.Inf(1)
	minZ = math.Inf(1)
	maxX = math.Inf(-1)
	maxY = math.Inf(-1)
	maxZ = math.Inf(-1)
	for _, i := range subset_indices {
		index = i * 3
		minX = math.Min(minX, m.Verts.Buffer[index])
		minY = math.Min(minY, m.Verts.Buffer[index+1])
		minZ = math.Min(minZ, m.Verts.Buffer[index+2])
		maxX = math.Max(maxX, m.Verts.Buffer[index])
		maxY = math.Max(maxY, m.Verts.Buffer[index+1])
		maxZ = math.Max(maxZ, m.Verts.Buffer[index+2])
	}
	return cb.New(minX, minY, minZ, maxX, maxY, maxZ)
}

// Identifies border vertices and returns an array of arrays representing closed
// loops of border vertices.
// Border vertices are identified as including a face which includes an edge
// which is only included in that one face.
func (m *Mesh) IdentifyBoundaries() (boundaries [][]int) {
	boundary_edges_slice := make([][2]int, 0)

	// Build up boundary_edges as a sequences of pairs of vertex indices
	// representing boundary edges
	m.Faces.UpdateIndex()
	for i := 0; i < m.Verts.Len(); i++ {
		face_triples := m.Faces.TriplesWith(i)
		vertex_counts := make(map[int]int)
		for _, v := range face_triples {
			// discard occurances of i
			// discard vertices less than i because they should already have been picked
			//  up in a previous iteration of the outerloop
			if v > i {
				vertex_counts[v]++
			}
		}
		for v, count := range vertex_counts {
			if count == 1 {
				boundary_edges_slice = append(boundary_edges_slice, [2]int{i, v})
			}
		}
	}

	// The purpose of the intermediate boundary_edges_slice is so the following
	// intermediate boundary_edges list will be sorted so that this function can
	// be idempotent.
	boundary_edges := list.New()
	sort.Sort(byIndices(boundary_edges_slice))
	for _, boundary_edge := range boundary_edges_slice {
		boundary_edges.PushBack(boundary_edge)
	}

	// Transform boundary_edges into one or more closed loops of connected
	// vertices using partials for termporary storage
	partials := list.New()
	for boundary_edges.Len() > 0 {
		latest_partial := make([]int, 0)
		first_edge := boundary_edges.Front().Value.([2]int)
		latest_partial = append(latest_partial, first_edge[0], first_edge[1])
		_ = boundary_edges.Remove(boundary_edges.Front())
		for true {
			head := latest_partial[0]
			tail := latest_partial[len(latest_partial)-1]
			if head == tail {
				break
			}
			for edge := boundary_edges.Front(); edge != nil; edge = edge.Next() {
				if edge.Value.([2]int)[0] == tail {
					latest_partial = append(latest_partial, edge.Value.([2]int)[1])
					boundary_edges.Remove(edge)
					break
				} else if edge.Value.([2]int)[1] == tail {
					latest_partial = append(latest_partial, edge.Value.([2]int)[0])
					boundary_edges.Remove(edge)
					break
				}
			}
		}
		partials.PushBack(latest_partial[:len(latest_partial)-1])
	}

	// Copy completed boundaries from partials over into boundaries
	boundaries = make([][]int, 0, partials.Len())
	for el := partials.Front(); el != nil; el = el.Next() {
		complete_boundary := make([]int, len(el.Value.([]int)), len(el.Value.([]int)))
		for i, boundary_vert := range el.Value.([]int) {
			complete_boundary[i] = boundary_vert
		}
		boundaries = append(boundaries, complete_boundary)
	}

	return
}

// for sorting the slice of edges... TODO: tidy this up somewhere

type byIndices [][2]int

func (v byIndices) Len() int      { return len(v) }
func (v byIndices) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v byIndices) Less(i, j int) bool {
	return v[i][0] < v[j][0] || (v[i][0] == v[j][0] && v[i][1] < v[j][1])
}
