gomesh
======

A go package for processing 3D meshes.

## Example Usage

```go
package main

import "fmt"
import "github.com/nat-n/gomesh/mesh"
import trans "github.com/nat-n/gomesh/transformation"

func main() {
	m, uhoh := mesh.ReadOBJFile("my_mesh.obj")
	if uhoh != nil {
		fmt.Println(uhoh)
	}

	// find center point of mesh
	center := m.BoundingBox().Center()

	// scale in place by factor, and rotate by rad
	factor := 2.5
	rad := 1.57

	// apply multiple transformations to my mesh
	// transformations will take effect in reverse order
	ts := []trans.Transformation{
		trans.Translation(center.X, center.Y, center.Z),
		trans.Rotation(rad, 0, 0, 1),
		trans.Scale(factor),
		trans.Translation(-center.X, -center.Y, -center.Z),
	}
	m.Transform(ts[0].Multiply(ts[1:]...))

	// create mesh of bounding box of my_mesh
	bb_mesh := mesh.NewFromCuboid(*m.BoundingBox())

	// write new obj files for the transformed mesh and its bounding box
	uhoh = m.WriteOBJFile("my_mesh_transformed.obj")
	if uhoh != nil {
		fmt.Println(uhoh)
	}
	uhoh = bb_mesh.WriteOBJFile("my_mesh_transformed_bb.obj")
	if uhoh != nil {
		fmt.Println(uhoh)
	}
}
```
