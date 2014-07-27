package cuboid

import "math"

/*
 * Represents axis-aligned cuboids. Ideal for 3d bounding boxes.
 */

type Cuboid struct {
	OriginX, OriginY, OriginZ, TerminusX, TerminusY, TerminusZ float64
}

// Constructs a new cuboid from the coordinates of the minimum and maximum
// corners.
func New(x1, y1, z1, x2, y2, z2 float64) *Cuboid {
	return &Cuboid{
		OriginX:   math.Min(x1, x2),
		OriginY:   math.Min(y1, y2),
		OriginZ:   math.Min(z1, z2),
		TerminusX: math.Max(x1, x2),
		TerminusY: math.Max(y1, y2),
		TerminusZ: math.Max(z1, z2),
	}
}

// Returns a new cuboid that is `distance` units larger in all 6 directions.
func (c *Cuboid) Expanded(distance float64) *Cuboid {
	return &Cuboid{
		OriginX:   c.OriginX - distance,
		OriginY:   c.OriginY - distance,
		OriginZ:   c.OriginZ - distance,
		TerminusX: c.TerminusX + distance,
		TerminusY: c.TerminusY + distance,
		TerminusZ: c.TerminusZ + distance,
	}
}

// Find the center point of the cuboid.
func (c *Cuboid) Center() (x, y, z float64) {
	x = c.OriginX + (c.TerminusX-c.OriginX)/2
	y = c.OriginY + (c.TerminusY-c.OriginY)/2
	z = c.OriginZ + (c.TerminusZ-c.OriginZ)/2
	return
}

// Checks if the cuboid contains a specific point.
func (c *Cuboid) Contains(x, y, z float64) bool {
	return x >= c.OriginX && x <= c.TerminusX &&
		y >= c.OriginY && y <= c.TerminusY &&
		z >= c.OriginZ && z <= c.TerminusZ
}

// Check if the cuboid shares some volume (or area) with another Cuboid
func (c *Cuboid) Intersects(other *Cuboid) bool {
	return !(c.TerminusY <= other.OriginY ||
		c.OriginY >= other.TerminusY ||
		c.OriginX >= other.TerminusX ||
		c.TerminusX <= other.OriginX ||
		c.TerminusZ <= other.OriginZ ||
		c.OriginZ >= other.TerminusZ)
}

// Returns a new cuboid which represents the bounding box of the cuboid and
// others.
func (c *Cuboid) Union(others ...Cuboid) (merged Cuboid) {
	merged = Cuboid{
		c.OriginX,
		c.OriginY,
		c.OriginZ,
		c.TerminusX,
		c.TerminusY,
		c.TerminusZ,
	}
	for _, other := range others {
		merged.OriginX = math.Min(merged.OriginX, other.OriginX)
		merged.OriginY = math.Min(merged.OriginY, other.OriginY)
		merged.OriginZ = math.Min(merged.OriginZ, other.OriginZ)
		merged.TerminusX = math.Max(merged.TerminusX, other.TerminusX)
		merged.TerminusY = math.Max(merged.TerminusY, other.TerminusY)
		merged.TerminusZ = math.Max(merged.TerminusZ, other.TerminusZ)
	}
	return
}
