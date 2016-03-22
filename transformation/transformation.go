package transformation

import "math"
import "github.com/nat-n/geom"

type Transformation [16]float64

func (t *Transformation) Cell(x int, y int) float64 {
	return t[y*4+x]
}

func (t *Transformation) Col(x int) [4]float64 {
	return [4]float64{t[x], t[x+4], t[x+8], t[x+12]}
}

func (t *Transformation) Row(y int) [4]float64 {
	return [4]float64{t[y*4], t[y*4+1], t[y*4+2], t[y*4+3]}
}

// Performs standard matrix multiplication t . m and returns the result as a new
//  Transform.
func (t *Transformation) Multiply(matrices ...Transformation) (n Transformation) {
	n = *t
	for _, m := range matrices {
		n = Transformation{
			n[0]*m[0] + n[1]*m[4] + n[2]*m[8] + n[3]*m[12],
			n[0]*m[1] + n[1]*m[5] + n[2]*m[9] + n[3]*m[13],
			n[0]*m[2] + n[1]*m[6] + n[2]*m[10] + n[3]*m[14],
			n[0]*m[3] + n[1]*m[7] + n[2]*m[11] + n[3]*m[15],

			n[4]*m[0] + n[5]*m[4] + n[6]*m[8] + n[7]*m[12],
			n[4]*m[1] + n[5]*m[5] + n[6]*m[9] + n[7]*m[13],
			n[4]*m[2] + n[5]*m[6] + n[6]*m[10] + n[7]*m[14],
			n[4]*m[3] + n[5]*m[7] + n[6]*m[11] + n[7]*m[15],

			n[8]*m[0] + n[9]*m[4] + n[10]*m[8] + n[11]*m[12],
			n[8]*m[1] + n[9]*m[5] + n[10]*m[9] + n[11]*m[13],
			n[8]*m[2] + n[9]*m[6] + n[10]*m[10] + n[11]*m[14],
			n[8]*m[3] + n[9]*m[7] + n[10]*m[11] + n[11]*m[15],

			n[12]*m[0] + n[13]*m[4] + n[14]*m[8] + n[15]*m[12],
			n[12]*m[1] + n[13]*m[5] + n[14]*m[9] + n[15]*m[13],
			n[12]*m[2] + n[13]*m[6] + n[14]*m[10] + n[15]*m[14],
			n[12]*m[3] + n[13]*m[7] + n[14]*m[11] + n[15]*m[15],
		}
	}
	return
}

// Performs standard matrix addition t + m and returns the result as a new
//  Transform.
func (t *Transformation) Add(matrices ...Transformation) (n Transformation) {
	n = *t
	for _, m := range matrices {
		n = Transformation{
			n[0] + m[0], n[1] + m[1], n[2] + m[2], n[3] + m[3],
			n[4] + m[4], n[5] + m[5], n[6] + m[6], n[7] + m[7],
			n[8] + m[8], n[9] + m[9], n[10] + m[10], n[11] + m[11],
			n[12] + m[12], n[13] + m[13], n[14] + m[14], n[15] + m[15],
		}
	}
	return
}

// Apply a translation to a given pointer of vector represented by a slice of
//  length 3.
func (t *Transformation) Apply(v []float64) {
	if len(v) != 3 {
		panic("Transformations can only be applied to slices of length 3.")
	}
	// multiply t with the given vector to produce a quaternion
	q := [4]float64{
		t[0]*v[0] + t[1]*v[1] + t[2]*v[2] + t[3],
		t[4]*v[0] + t[5]*v[1] + t[6]*v[2] + t[7],
		t[8]*v[0] + t[9]*v[1] + t[10]*v[2] + t[11],
		t[12]*v[0] + t[13]*v[1] + t[14]*v[2] + t[15],
	}

	// Complete the transform and update v
	v[0] = q[0] / q[3]
	v[1] = q[1] / q[3]
	v[2] = q[2] / q[3]
}

func (t *Transformation) ApplyToVec3(vs ...geom.Vec3I) {
	// multiply t with the given vector to produce a quaternion
	for _, v := range vs {
		q := [4]float64{
			t[0]*v.GetX() + t[1]*v.GetY() + t[2]*v.GetZ() + t[3],
			t[4]*v.GetX() + t[5]*v.GetY() + t[6]*v.GetZ() + t[7],
			t[8]*v.GetX() + t[9]*v.GetY() + t[10]*v.GetZ() + t[11],
			t[12]*v.GetX() + t[13]*v.GetY() + t[14]*v.GetZ() + t[15],
		}

		// Complete the transform and update v
		v.SetX(q[0] / q[3])
		v.SetY(q[1] / q[3])
		v.SetZ(q[2] / q[3])
	}
}

func Translation(x, y, z float64) Transformation {
	return Transformation{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
}

func Scale(f float64) Transformation {
	return Transformation{
		f, 0, 0, 0,
		0, f, 0, 0,
		0, 0, f, 0,
		0, 0, 0, 1,
	}
}

func ScaleDimensions(x, y, z float64) Transformation {
	return Transformation{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

// Generates a transformation matrix for rotation by theta radians around the
//  (presumed) unit vector [px, py, pz]
func Rotation(theta, px, py, pz float64) Transformation {
	FLOAT_EQUALITY_THRESHOLD := 0.0000001

	// unpack reusable values
	pxsqr := px * px
	pysqr := py * py
	pzsqr := pz * pz
	ct := math.Cos(theta)
	mct := 1 - ct
	st := math.Sin(theta)

	// ensure pivot vector has a length of 1
	pivotVectorLength := math.Sqrt(pxsqr + pysqr + pzsqr)
	if math.Abs(pivotVectorLength-1) > FLOAT_EQUALITY_THRESHOLD {
		px /= pivotVectorLength
		py /= pivotVectorLength
		pz /= pivotVectorLength
	}

	return Transformation{
		ct + pxsqr*mct, px*py*mct - pz*st, px*pz*mct + py*st, 0,
		py*px*mct + pz*st, ct + pysqr*mct, py*pz*mct - px*st, 0,
		pz*px*mct - py*st, pz*py*mct + px*st, ct + pzsqr*mct, 0,
		0, 0, 0, 1,
	}
}
