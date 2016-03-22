package mesh

import (
	"bufio"
	"errors"
	"github.com/nat-n/geom"
	"io"
	"os"
	"strconv"
	"strings"
)

// Populate this Mesh from the given OBJ file
func LoadOBJ(obj_reader *io.Reader) (m *Mesh, err error) {
	// prepare for data
	m = New("")

	// setup for parsing
	var (
		line  string
		words []string
	)
	line_no := -1

	normalsBuffer := make([]*geom.Vec3, 0)
	facesBuffer := make([][3]int, 0)

	// open and parse file
	scanner := bufio.NewScanner(*obj_reader)
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
				err = newParseError("OBJ", line_no)
				return
			}
			m.Vertices.Append(&Vertex{
				Vec3:   geom.Vec3{floats[0], floats[1], floats[2]},
				Meshes: make(map[Mesh]int)})
		case "vn":
			// read in a vertex normal
			floats, parseErr := parse3Floats(words[1:])
			if parseErr != nil {
				err = newParseError("OBJ", line_no)
				return
			}
			normalsBuffer = append(normalsBuffer,
				&geom.Vec3{floats[0], floats[1], floats[2]})
		case "f":
			// read in a faces
			ints, parseErr := parse3Ints(words[1:])
			if parseErr != nil {
				err = newParseError("OBJ", line_no)
				return
			}
			facesBuffer = append(facesBuffer,
				[3]int{ints[0] - 1, ints[1] - 1, ints[2] - 1})
		default:
			err = newParseError("OBJ", line_no)
			return
		}
	}

	for i, n := range normalsBuffer {
		m.Vertices.Get(i)[0].SetNormal(n)
	}
	for i, f := range facesBuffer {
		abc := m.Vertices.Get(f[0], f[1], f[2])
		a, b, c := abc[0], abc[1], abc[2]
		m.Faces.Append(&Face{Vertices: [3]VertexI{a, b, c}, Mesh: *m, Index: i})
	}

	return
}

// Write this mesh to a new obj file.
func (m *Mesh) WriteOBJ(obj_writer io.Writer) (err error) {
	// track where vertices were written
	vert_lookup := make(map[VertexI]int)

	// Write Vertices
	for i := 0; i < m.Vertices.Len(); i++ {
		v := m.Vertices.Get(i)[0]
		vert_lookup[v] = i
		_, err = obj_writer.Write([]byte(
			"v " + strconv.FormatFloat(v.GetX(), 'f', -1, 64) +
				" " + strconv.FormatFloat(v.GetY(), 'f', -1, 64) +
				" " + strconv.FormatFloat(v.GetZ(), 'f', -1, 64) +
				"\n",
		))
	}
	// It seems improbably that an error would occur for writing vertex but not
	// the last one so only check the last one.
	if err != nil {
		err = errors.New(
			"Error occured when attempting to write vertices to obj file.",
		)
		return
	}

	// Write normals, unless first vector has no normal
	for i := 0; i < m.Vertices.Len(); i++ {
		v := m.Vertices.Get(i)[0]
		n := v.GetNormal()
		if n == nil {
			v.CalculateNormal()
			n = v.GetNormal()
		}
		_, err = obj_writer.Write([]byte(
			"vn " + strconv.FormatFloat(n.GetX(), 'f', -1, 64) +
				" " + strconv.FormatFloat(n.GetY(), 'f', -1, 64) +
				" " + strconv.FormatFloat(n.GetZ(), 'f', -1, 64) +
				"\n",
		))
	}
	// It seems improbably that an error would occur for writing vertex but not
	// the last one so only check the last one.
	if err != nil {
		err = errors.New(
			"Error occured when attempting to write vertex normals obj file.",
		)
		return
	}

	// Write faces
	for i := 0; i < m.Faces.Len(); i++ {
		f := m.Faces.Get(i)[0]
		_, err = obj_writer.Write([]byte(
			"f " + strconv.Itoa(vert_lookup[f.GetA()]+1) +
				" " + strconv.Itoa(vert_lookup[f.GetB()]+1) +
				" " + strconv.Itoa(vert_lookup[f.GetC()]+1) +
				"\n",
		))
	}
	// It seems improbable that an error would occur for writing vertex but not
	// the last one so only check the last one.
	if err != nil {
		err = errors.New(
			"Error occured when attempting to write faces to obj file.",
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
func (m *Mesh) WriteSTL(stl_path string) (err error) {
	err = errors.New(
		"STL writing not yet implemented",
	)
	return
}

func ReadOBJFile(input_path string) (m *Mesh, err error) {
	// Open file
	input_file, err := os.Open(input_path)
	if err != nil {
		return
	}
	defer input_file.Close()

	// Read from file
	mesh_reader := io.Reader(input_file)
	m, err = LoadOBJ(&mesh_reader)

	return
}

func (m *Mesh) WriteOBJFile(output_path string) (err error) {
	// Serialized JSON and stream to a file
	output_file, err := os.Create(output_path)
	if err != nil {
		return
	}
	defer output_file.Close()
	err = m.WriteOBJ(io.Writer(output_file))

	return
}

func newParseError(fileType string, line_no int) error {
	return errors.New(
		"Error parsing " + fileType + " file on line: " + strconv.Itoa(line_no))
}
