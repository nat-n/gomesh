package mesh

import (
	"bufio"
	"container/list"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
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

func New(name string) *Mesh {
	return &Mesh{
		Name:  name,
		Verts: tb.NewVertexBuffer(),
		Norms: tb.NewVectorBuffer(),
		Faces: tb.NewTriangleBuffer(),
	}
}

// apply the given transformation to every vertex
func (m *Mesh) Transform(t tr.Transformation) {
	for i := 0; i < m.Verts.Len(); i++ {
		t.Apply(m.Verts.Buffer[i*3 : i*3+3])
	}
}

// apply the given transformation to each vertex in indices
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
//  loops of border vertices.
// Border vertices are identified as including a face which includes and edge
//  which is only included in that one face.
func (m *Mesh) IdentifyBoundaries() (boundaries [][]int) {
	partials := list.New()

	boundary_edges := list.New()

	// Build up boundary_edges as a sequences of pairs of vertex indices
	//  representing boundary edges
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
				boundary_edges.PushBack([2]int{i, v})
			}
		}
	}

	// Transform boundary_edges into one or more closed loops of connected
	//  vertices using partials for termporary storage
	for boundary_edges.Len() > 0 {
		latest_partial := make([]int, 0, 1000)
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

// Populate this Mesh from the given OBJ file
func LoadOBJ(obj_path string) (m *Mesh, err error) {
	// prepare for data
	m = New("")
	m.Verts = tb.NewVertexBuffer()
	m.Norms = tb.NewVectorBuffer()
	m.Faces = tb.NewTriangleBuffer()

	// setup for parsing
	var (
		line  string
		words []string
	)
	line_no := -1

	// open and parse file
	file, _ := os.Open(obj_path)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line_no++
		// trim leading and trailing whitespace
		line = strings.TrimSpace(scanner.Text())
		// firstly discard anything on this line after a #
		if comment_start := strings.Index(line, "#"); comment_start >= 0 {
			line = line[:comment_start]
		}
		// ignore empty lines
		if len(line) == 0 {
			continue
		}
		words = strings.Fields(line)
		switch words[0] {
		case "v":
			// read in a vertex
			floats, parseErr := parse3Floats(words[1:])
			if parseErr != nil {
				err = errors.New(
					"Error parsing OBJ file on line: " +
						strconv.Itoa(line_no))
				return
			}
			m.Verts.Append(floats[0], floats[1], floats[2])
		case "vn":
			// read in a vertex normal
			floats, parseErr := parse3Floats(words[1:])
			if parseErr != nil {
				err = errors.New(
					"Error parsing OBJ file on line: " +
						strconv.Itoa(line_no))
				return
			}
			m.Norms.Append(floats[0], floats[1], floats[2])
		case "f":
			// read in a faces
			ints, parseErr := parse3Ints(words[1:])
			if parseErr != nil {
				err = errors.New(
					"Error parsing OBJ file on line: " +
						strconv.Itoa(line_no))
				return
			}
			m.Faces.Append(ints[0]-1, ints[1]-1, ints[2]-1)
		default:
			err = errors.New(
				"Error parsing OBJ file on line: " +
					strconv.Itoa(line_no))
			return
		}

	}
	return
}

// Write this mesh to a new obj file.
func (m *Mesh) WriteOBJ(obj_path string) (problem error) {
	f, err := os.Create(obj_path)
	if err != nil {
		problem = errors.New(
			"Error occured when attempting to write obj file: " + obj_path,
		)
		return
	}
	for i := 0; i < m.Verts.Len(); i++ {
		_, err = f.Write([]byte(
			"v " + strconv.FormatFloat(m.Verts.Buffer[i*3], 'f', -1, 64) +
				" " + strconv.FormatFloat(m.Verts.Buffer[i*3+1], 'f', -1, 64) +
				" " + strconv.FormatFloat(m.Verts.Buffer[i*3+2], 'f', -1, 64) +
				"\n",
		))
	}
	// It seems improbably that an error would occur for writing vertex but not
	//  the last one so only check the last one.
	if err != nil {
		problem = errors.New(
			"Error occured when attempting to write vertices to obj file: " +
				obj_path,
		)
		return
	}

	// Write out normals, unless the norms buffer is empty
	if !m.Norms.IsEmpty() {
		for i := 0; i < m.Norms.Len(); i++ {
			_, err = f.WriteString(
				"vn " + strconv.FormatFloat(m.Norms.Buffer[i*3], 'f', -1, 64) +
					" " + strconv.FormatFloat(m.Norms.Buffer[i*3+1], 'f', -1, 64) +
					" " + strconv.FormatFloat(m.Norms.Buffer[i*3+2], 'f', -1, 64) +
					"\n",
			)
		}
	}
	// It seems improbably that an error would occur for writing vertex but not
	// the last one so only check the last one.
	if err != nil {
		problem = errors.New(
			"Error occured when attempting to write vertex normals obj file: " +
				obj_path,
		)
		return
	}

	// Write out faces
	for i := 0; i < m.Faces.Len(); i++ {
		_, err = f.Write([]byte(
			"f " + strconv.Itoa(m.Faces.Buffer[i*3]+1) +
				" " + strconv.Itoa(m.Faces.Buffer[i*3+1]+1) +
				" " + strconv.Itoa(m.Faces.Buffer[i*3+2]+1) +
				"\n",
		))
	}
	// It seems improbably that an error would occur for writing vertex but not
	// the last one so only check the last one.
	if err != nil {
		problem = errors.New(
			"Error occured when attempting to write faces to obj file: " + obj_path,
		)
		return
	}
	return
}

// Read in STL file as new mesh. Not yet implemented
func LoadSTL(stl_path string) (m *Mesh, err error) {
	err = errors.New(
		"STL reading not yet implemented",
	)
	return
}

// Write mesh to STL file. Not yet implemented
func (m *Mesh) WriteSTL(stl_path string) (problem error) {
	problem = errors.New(
		"STL writing not yet implemented",
	)
	return
}

// Constructs a new Mesh from a cuboid.
func FromCuboid(c *cb.Cuboid) (m Mesh) {
	m = Mesh{
		Name:  "Cuboid",
		Verts: tb.NewVertexBuffer(),
		Norms: tb.NewVectorBuffer(),
		Faces: tb.NewTriangleBuffer(),
	}
	m.Verts.Buffer = []float64{
		c.OriginX, c.OriginY, c.OriginZ,
		c.TerminusX, c.OriginY, c.OriginZ,
		c.TerminusX, c.OriginY, c.TerminusZ,
		c.OriginX, c.OriginY, c.TerminusZ,
		c.OriginX, c.TerminusY, c.OriginZ,
		c.TerminusX, c.TerminusY, c.OriginZ,
		c.TerminusX, c.TerminusY, c.TerminusZ,
		c.OriginX, c.TerminusY, c.TerminusZ,
	}
	m.Faces.Buffer = []int{
		0, 2, 1,
		0, 3, 2,
		0, 5, 1,
		0, 4, 5,
		1, 6, 2,
		1, 5, 6,
		2, 7, 3,
		2, 6, 7,
		3, 4, 0,
		3, 7, 4,
		5, 7, 6,
		5, 4, 7,
	}
	return
}
