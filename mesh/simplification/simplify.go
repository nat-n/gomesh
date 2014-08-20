package simplification

import (
	"container/heap"
	"fmt"
	"github.com/nat-n/gomesh/mesh"
	"math"
)

import tb "github.com/nat-n/gomesh/triplebuffer"

type vertex struct {
	Coords     [3]float64
	Faces      []*face
	Q          *Quadric
	Edges      []*edge
	FinalIndex int
	Collapsed  bool
}

type face struct {
	Verts     [3]*vertex
	Kp        *Quadric
	Collapsed bool
	Edges     []*edge // should be [3]*edge
}

type edge struct {
	V1             *vertex
	V2             *vertex
	CollapseTarget [3]float64
	Q              *Quadric
	Error          float64
	Faces          []*face // should be [2]*edge??
	Removed        bool
}

// Compute Euclidean Distance between the two vertices of the edge
func (e *edge) Length() float64 {
	dx := e.V1.Coords[0] - e.V2.Coords[0]
	dy := e.V1.Coords[1] - e.V2.Coords[1]
	dz := e.V1.Coords[2] - e.V2.Coords[2]
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// Interface and Convenience functions to for our heap of edges
type edgeHeap []*edge

func (h edgeHeap) Len() int           { return len(h) }
func (h edgeHeap) Less(i, j int) bool { return h[i].Error < h[j].Error }
func (h edgeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *edgeHeap) Push(x interface{}) {
	e := x.(*edge)
	*h = append(*h, e)
}

func (h *edgeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *edgeHeap) Fix(indices ...int) {
	for _, index := range indices {
		heap.Fix(h, index)
	}
}

// Find and call fix on each of the affected edges,
// this is inefficient, but not too bad, and I'm not sure how to avoid it
func (h *edgeHeap) UpdateEdges(affected_edges []*edge) {
	possible_matches := make([]*edge, 0, len(affected_edges))
	copy(possible_matches, affected_edges)
	eh := *h
	for i, e := range eh {
		i++
		for j, possible_match := range possible_matches {
			if e == possible_match {
				heap.Fix(h, i)
				possible_matches = append(
					possible_matches[:j],
					possible_matches[j+1:]...,
				)
				break
			}
		}
	}
}

// Not sure if this is needed
func (f *face) IncludesVertex(v1 *vertex) bool {
	for _, v2 := range f.Verts {
		if v1 == v2 {
			return true
		}
	}
	return false
}

// Calculate the Optimal collapse location for this edge and the associated
// error, and update the edge with these values.
func (e *edge) calculateError() {
	// Calculate error quadric for this edge as sum of vertex error quadrics
	Q := Quadric{}
	e.Q = &Q
	e.Q.Add(e.V1.Q)
	e.Q.Add(e.V2.Q)

	det := Q.Determinant()
	// if det != 0 && false { // BLOCKED
	e.CollapseTarget[0] = (Q[3]*Q[5]*Q[5] - Q[2]*Q[6]*Q[5] - Q[3]*Q[4]*Q[7] + Q[1]*Q[6]*Q[7] + Q[2]*Q[4]*Q[8] - Q[1]*Q[5]*Q[8]) / det
	e.CollapseTarget[1] = (Q[2]*Q[6]*Q[2] - Q[3]*Q[5]*Q[2] + Q[3]*Q[1]*Q[7] - Q[0]*Q[6]*Q[7] - Q[2]*Q[1]*Q[8] + Q[0]*Q[5]*Q[8]) / det
	e.CollapseTarget[2] = (Q[3]*Q[4]*Q[2] - Q[1]*Q[6]*Q[2] - Q[3]*Q[1]*Q[5] + Q[0]*Q[6]*Q[5] + Q[1]*Q[1]*Q[8] - Q[0]*Q[4]*Q[8]) / det

	// } else {

	// Determine which is best, V1, V2 or their midpoint
	midpoint := [3]float64{
		(e.V1.Coords[0] + e.V2.Coords[0]) / 2,
		(e.V1.Coords[1] + e.V2.Coords[1]) / 2,
		(e.V1.Coords[2] + e.V2.Coords[2]) / 2,
	}
	v1_error := Q.VertexError(e.V1.Coords[0], e.V1.Coords[1], e.V1.Coords[2])
	v2_error := Q.VertexError(e.V2.Coords[0], e.V2.Coords[1], e.V2.Coords[2])
	midpoint_error := Q.VertexError(midpoint[0], midpoint[1], midpoint[2])

	e.Error = Q.VertexError(
		e.CollapseTarget[0],
		e.CollapseTarget[1],
		e.CollapseTarget[2],
	)

	if e.Error < v1_error && e.Error < v2_error && e.Error < midpoint_error {
		fmt.Println("THIS NEVER HAPPENS :( ")
		fmt.Println("The present implementation of the textbook approach is somehow very broken")
		fmt.Println("in that the resulting vertex locations are all on a line.")
	} else if v1_error < v2_error {
		if v1_error < midpoint_error {
			e.CollapseTarget = e.V1.Coords
		} else {
			e.CollapseTarget = midpoint
		}
	} else {
		if v2_error < midpoint_error {
			e.CollapseTarget = e.V2.Coords
		} else {
			e.CollapseTarget = midpoint
		}
	}

	// }

	e.Error = Q.VertexError(
		e.CollapseTarget[0],
		e.CollapseTarget[1],
		e.CollapseTarget[2],
	)
}

func (e *edge) collapse(threshold float64) (did_collapse bool) {
	did_collapse = false

	if e.Removed {
		return
	}
	e.Removed = true

	// Lazily compute edge length, and compare to threshold
	if e.Length() > threshold {
		return
	}

	// Skip edges where a neighbor from v1 is also a neighbor of v2
	// but there's is no face shared by the three vertices.
	for _, v1_edge := range e.V1.Edges {
		var v1_other *vertex
		if v1_edge.V1 == e.V1 {
			v1_other = e.V2
		} else {
			v1_other = e.V1
		}
		for _, v2_edge := range e.V1.Edges {
			var v2_other *vertex
			if v2_edge.V1 == e.V1 {
				v2_other = e.V2
			} else {
				v2_other = e.V1
			}
			if v1_other == v2_other {
				// check if there is a face shared by all three
				found_face := false
				for _, v1_face := range e.V1.Faces {
					if (e.V1 == v1_face.Verts[0] || e.V1 == v1_face.Verts[1] || e.V1 == v1_face.Verts[2]) &&
						(e.V2 == v1_face.Verts[0] || e.V2 == v1_face.Verts[1] || e.V2 == v1_face.Verts[2]) &&
						(v1_other == v1_face.Verts[0] || v1_other == v1_face.Verts[1] || v1_other == v1_face.Verts[2]) {
						found_face = true
					}
				}
				if !found_face {
					fmt.Println("skipping super-triangle!")

					// in this case maybe we should actually considering replacing the
					// three faces with one (depending on distance of center vertex from
					// the average of the three corner vertices)

					return
				}
			}
		}
	}

	// temporary... seems to reduce artifacts
	for _, v1_edge := range e.V1.Edges {
		v1_edge.Removed = true
	}
	for _, v2_edge := range e.V2.Edges {
		v2_edge.Removed = true
	}

	// Update V1 to the new location and Q
	e.V1.Coords[0] = e.CollapseTarget[0]
	e.V1.Q = e.Q

	// Mark V2 as collapsed
	e.V2.Collapsed = true

	// Update faces of V2
	// - Mark faces on edge `e` as collapsed
	// - For other faces:
	//  - Update references to V2 to V1
	//  - Register faces with V1
	for _, f := range e.V2.Faces {
		for _, f_edge := range f.Edges {
			if f_edge == e {
				f.Collapsed = true
				break
			}
		}
		if !f.Collapsed {
			if f.Verts[0] == e.V2 {
				f.Verts[0] = e.V1
			} else if f.Verts[1] == e.V2 {
				f.Verts[1] = e.V1
			} else if f.Verts[2] == e.V2 {
				f.Verts[2] = e.V1
			} else {
				// if this face doesn't reference V1 then continue without updating
				// e.V1.Faces
				continue
			}
			e.V1.Faces = append(e.V1.Faces, f)
		}
	}

	// RETRYING THE STUPID WAY
	// IDENTIFY NEIGHBORS OF V1
	// REMAP EDGE FROM V2 TO V1 UNLESS THEY INCLUDE ANOTHER NEIGHBOR OF V1
	V1_neighbors := make([]*vertex, 0, len(e.V1.Edges)-1)
	for _, v1_edge := range e.V1.Edges {
		if v1_edge == e {
			continue
		}
		if v1_edge.V1 == e.V1 {
			V1_neighbors = append(V1_neighbors, v1_edge.V2)
		} else if v1_edge.V2 == e.V1 {
			V1_neighbors = append(V1_neighbors, v1_edge.V1)
		} else {
			panic("V1 edge ref error!")
		}
	}
	for _, v2_edge := range e.V2.Edges {
		if v2_edge == e {
			continue
		}
		var V2_neighbor *vertex
		if v2_edge.V1 == e.V2 {
			v2_edge.V1 = e.V1
			V2_neighbor = v2_edge.V2
		} else if v2_edge.V2 == e.V2 {
			v2_edge.V2 = e.V1
			V2_neighbor = v2_edge.V1
		} else {
			panic("V2 edge ref error!")
		}
		for _, V1_neighbor := range V1_neighbors {
			if V1_neighbor == V2_neighbor {
				v2_edge.Removed = true
				break
			}
		}
		if !v2_edge.Removed {
			e.V1.Edges = append(e.V1.Edges, v2_edge)
		}
	}

	// Update Q for all edges of V1
	for _, v1_edge := range e.V1.Edges {
		v1_edge.calculateError()
	}

	did_collapse = true
	return
}

// Calculate the Kp fundemental error matrix of a face, quadric of plane
func (f *face) calculateKp() {
	a, b, c := tb.Normal(
		f.Verts[0].Coords[:],
		f.Verts[1].Coords[:],
		f.Verts[2].Coords[:],
	)
	// use center point of triangle is better?
	cx := (f.Verts[0].Coords[0] + f.Verts[1].Coords[0] + f.Verts[2].Coords[0]) / 3
	cy := (f.Verts[0].Coords[1] + f.Verts[1].Coords[1] + f.Verts[2].Coords[1]) / 3
	cz := (f.Verts[0].Coords[2] + f.Verts[1].Coords[2] + f.Verts[2].Coords[2]) / 3
	d := -(a*cx + b*cy + c*cz)

	f.Kp = &Quadric{
		a * a, a * b, a * c, a * d,
		b * b, b * c, b * d,
		c * c, c * d,
		d * d,
	}
}

// Quadric Edge Collapse Decimation
// threshold: is the maximum length edge that will be contracted
// target_face_count: is interpreted based on the assumption that every
// collapsed edge removes two faces.
// safer_mode: if true then at most one edge associated with each vertex will be
// collapsed, this seems to reduce artifacts for equivalent performance, though
// less can be achieved per invokation.
func QECD(m *mesh.Mesh, threshold float64, target_face_count int, safer_mode bool) {

	vertices := make([]*vertex, 0)
	faces := make([]*face, 0)
	edges := &edgeHeap{}

	// build up vertices
	m.Verts.Each(func(x, y, z float64) {
		vertices = append(vertices, &vertex{
			Coords: [3]float64{x, y, z},
			Faces:  make([]*face, 0),
			Q:      &Quadric{},
			Edges:  make([]*edge, 0),
		})
	})
	// Build up faces and update verts
	// iterate through faces and collect non-border edges
	// by counting the occurances of every edge, and keeping those with a count of 2
	edge_occurances := make(map[[2]int][]*face)
	m.Faces.Each(func(a, b, c int) {
		new_face := &face{
			Verts: [3]*vertex{vertices[a], vertices[b], vertices[c]},
		}
		faces = append(faces, new_face)
		new_face.calculateKp()

		vertices[a].Faces = append(vertices[a].Faces, new_face)
		vertices[a].Q.Add(new_face.Kp)
		vertices[b].Faces = append(vertices[b].Faces, new_face)
		vertices[b].Q.Add(new_face.Kp)
		vertices[c].Faces = append(vertices[c].Faces, new_face)
		vertices[c].Q.Add(new_face.Kp)

		var edge_description [2]int
		if a < b {
			edge_description = [2]int{a, b}
		} else {
			edge_description = [2]int{b, a}
		}
		edge_occurances[edge_description] = append(
			edge_occurances[edge_description],
			new_face,
		)
		if a < c {
			edge_description = [2]int{a, c}
		} else {
			edge_description = [2]int{c, a}
		}
		edge_occurances[edge_description] = append(
			edge_occurances[edge_description],
			new_face,
		)
		if c < b {
			edge_description = [2]int{c, b}
		} else {
			edge_description = [2]int{b, c}
		}
		edge_occurances[edge_description] = append(
			edge_occurances[edge_description],
			new_face,
		)
	})

	// First iterate through edge_occurances to identify border vertices
	boundary_vertices := make(map[int]bool)
	for e, occurances := range edge_occurances {
		if len(occurances) == 1 {
			boundary_vertices[e[0]] = true
			boundary_vertices[e[1]] = true
		}
	}

	// Iterate through edge_occurances again and build up edges
	for e, occurances := range edge_occurances {
		// Skip if edge doesn't occur in two faces (boundary or non-manifold)
		if len(occurances) == 2 {
			// Skip edge if one of the vertices is on a mesh boundary
			_, is_boundary1 := boundary_vertices[e[0]]
			_, is_boundary2 := boundary_vertices[e[1]]
			if !is_boundary1 && !is_boundary2 {
				new_edge := &edge{
					V1:      vertices[e[0]],
					V2:      vertices[e[1]],
					Faces:   occurances,
					Removed: false,
				}
				for _, occurance := range occurances {
					occurance.Edges = append(occurance.Edges, new_edge)
				}
				new_edge.calculateError()
				vertices[e[0]].Edges = append(vertices[e[0]].Edges, new_edge)
				vertices[e[1]].Edges = append(vertices[e[1]].Edges, new_edge)
				edges.Push(new_edge)
			}
		}
	}

	// Sort edges by error
	heap.Init(edges)

	// Iteratively Collapse the lowest error edges, resorting after each collapse
	// just arbitrarily half the original number of edges for starters
	edges_collapse_target := (len(faces) - target_face_count) / 2
	for len(*edges) > 0 && edges_collapse_target > 0 {
		lowest_cost_edge := heap.Pop(edges).(*edge)
		did_collapse := lowest_cost_edge.collapse(threshold)
		if did_collapse {
			edges_collapse_target--
		}
		edges.UpdateEdges(lowest_cost_edge.V1.Edges)
	}

	//
	// Update the mesh with the changes made to vertices and faces
	//
	m.Verts = tb.NewVertexBuffer()
	m.Norms = tb.NewVectorBuffer()
	m.Faces = tb.NewTriangleBuffer()

	i := 0
	for _, v := range vertices {
		if !v.Collapsed {
			v.FinalIndex = i
			m.Verts.Buffer = append(m.Verts.Buffer, v.Coords[:]...)
			i++
		}
	}

	for _, f := range faces {
		if !(f.Collapsed || f.Verts[0].Collapsed || f.Verts[1].Collapsed || f.Verts[2].Collapsed) {
			a := f.Verts[0].FinalIndex
			b := f.Verts[1].FinalIndex
			c := f.Verts[2].FinalIndex
			m.Faces.Buffer = append(m.Faces.Buffer, a, b, c)
		}
	}

}
