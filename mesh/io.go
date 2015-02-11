package mesh

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)
import tb "github.com/nat-n/gomesh/triplebuffer"

// Populate this Mesh from the given OBJ file
func LoadOBJ(obj_reader *io.Reader) (m *Mesh, err error) {
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
func (m *Mesh) WriteOBJ(obj_writer io.Writer) (err error) {
	// Write Vertices
	for i := 0; i < m.Verts.Len(); i++ {
		_, err = obj_writer.Write([]byte(
			"v " + strconv.FormatFloat(m.Verts.Buffer[i*3], 'f', -1, 64) +
				" " + strconv.FormatFloat(m.Verts.Buffer[i*3+1], 'f', -1, 64) +
				" " + strconv.FormatFloat(m.Verts.Buffer[i*3+2], 'f', -1, 64) +
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

	// Write normals, unless the Norms buffer is empty
	if !m.Norms.IsEmpty() {
		for i := 0; i < m.Norms.Len(); i++ {
			_, err = obj_writer.Write([]byte(
				"vn " + strconv.FormatFloat(m.Norms.Buffer[i*3], 'f', -1, 64) +
					" " + strconv.FormatFloat(m.Norms.Buffer[i*3+1], 'f', -1, 64) +
					" " + strconv.FormatFloat(m.Norms.Buffer[i*3+2], 'f', -1, 64) +
					"\n",
			))
		}
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
		_, err = obj_writer.Write([]byte(
			"f " + strconv.Itoa(m.Faces.Buffer[i*3]+1) +
				" " + strconv.Itoa(m.Faces.Buffer[i*3+1]+1) +
				" " + strconv.Itoa(m.Faces.Buffer[i*3+2]+1) +
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
