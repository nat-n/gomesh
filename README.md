gomesh
======

A go package for processing 3D meshes.

## Example Usage

    package main
    
    import "fmt"
    import "github.com/nat-n/gomesh/mesh"
    import trans "github.com/nat-n/gomesh/transformation"

    func main() {
    	m, uhoh := mesh.LoadOBJ("path/to/my_mesh.obj")
    	if uhoh != nil {
    		fmt.Println(uhoh)
    	}
    	
    	// find center point of mesh
    	cx, cy, cz := m.BoundingBox().Center()
    
    	// scale in place by f, and rotate by r
    	f := 2.5
    	r := 1.57
    	
    	// apply multiple transformations to my mesh
		// transformations will take effect in reverse order
    	ts := []trans.Transformation{
    		trans.Translation(cx, cy, cz),
    		trans.Rotation(r, 0, 0, 1),
    		trans.Scale(f),
    		trans.Translation(-cx, -cy, -cz),
    	}
    	m.Transform(ts[0].Multiply(ts[1:]...))
    
    	// create mesh of bounding box of my_mesh
    	bb_mesh := mesh.FromCuboid(m.BoundingBox())
    
        // write new obj files for the transformed mesh and its bounding box
    	uhoh = m.WriteOBJ("my_mesh_transformed.obj")
    	if uhoh != nil {
    		fmt.Println(uhoh)
    	}
    	
    	uhoh = bb_mesh.WriteOBJ("my_mesh_bb.obj")
    	if uhoh != nil {
    		fmt.Println(uhoh)
    	}
    }
