package simplification

type Quadric [10]float64

func (q1 *Quadric) Add(q2 *Quadric) {
	q1[0] += q2[0]
	q1[1] += q2[1]
	q1[2] += q2[2]
	q1[3] += q2[3]
	q1[4] += q2[4]
	q1[5] += q2[5]
	q1[6] += q2[6]
	q1[7] += q2[7]
	q1[8] += q2[8]
	q1[9] += q2[9]
}

// 0  => 0,  1  => 1,  2  => 2,  3 => 3,
// 4  => 1,  5  => 4,  6  => 5,  7 => 6,
// 8  => 2,  9  => 5,  10 => 7, 11 => 8,
// 12 => 3,  13 => 6,  14 => 8, 15 => 9,

func (q *Quadric) Inverse() (result_ref *Quadric, can_invert bool) {
	result := Quadric{}
	result_ref = &result
	can_invert = true
	det := q.Determinant()
	if det == 0 {
		can_invert = false
		return
	}
	result[0] = (q[5]*q[8]*q[6] - q[6]*q[7]*q[6] + q[6]*q[5]*q[8] - q[4]*q[8]*q[8] - q[5]*q[5]*q[9] + q[4]*q[7]*q[9]) / det
	result[1] = (q[3]*q[7]*q[6] - q[2]*q[8]*q[6] - q[3]*q[5]*q[8] + q[1]*q[8]*q[8] + q[2]*q[5]*q[9] - q[1]*q[7]*q[9]) / det
	result[2] = (q[2]*q[6]*q[6] - q[3]*q[5]*q[6] + q[3]*q[4]*q[8] - q[1]*q[6]*q[8] - q[2]*q[4]*q[9] + q[1]*q[5]*q[9]) / det
	result[3] = (q[3]*q[5]*q[5] - q[2]*q[6]*q[5] - q[3]*q[4]*q[7] + q[1]*q[6]*q[7] + q[2]*q[4]*q[8] - q[1]*q[5]*q[8]) / det
	result[4] = (q[2]*q[8]*q[3] - q[3]*q[7]*q[3] + q[3]*q[2]*q[8] - q[0]*q[8]*q[8] - q[2]*q[2]*q[9] + q[0]*q[7]*q[9]) / det
	result[5] = (q[3]*q[5]*q[3] - q[2]*q[6]*q[3] - q[3]*q[1]*q[8] + q[0]*q[6]*q[8] + q[2]*q[1]*q[9] - q[0]*q[5]*q[9]) / det
	result[6] = (q[2]*q[6]*q[2] - q[3]*q[5]*q[2] + q[3]*q[1]*q[7] - q[0]*q[6]*q[7] - q[2]*q[1]*q[8] + q[0]*q[5]*q[8]) / det
	result[7] = (q[1]*q[6]*q[3] - q[3]*q[4]*q[3] + q[3]*q[1]*q[6] - q[0]*q[6]*q[6] - q[1]*q[1]*q[9] + q[0]*q[4]*q[9]) / det
	result[8] = (q[3]*q[4]*q[2] - q[1]*q[6]*q[2] - q[3]*q[1]*q[5] + q[0]*q[6]*q[5] + q[1]*q[1]*q[8] - q[0]*q[4]*q[8]) / det
	result[9] = (q[1]*q[5]*q[2] - q[2]*q[4]*q[2] + q[2]*q[1]*q[5] - q[0]*q[5]*q[5] - q[1]*q[1]*q[7] + q[0]*q[4]*q[7]) / det
	return
}

func (q *Quadric) Determinant() float64 {
	return q[3]*q[5]*q[5]*q[3] - q[2]*q[6]*q[5]*q[3] - q[3]*q[4]*q[7]*q[3] + q[1]*q[6]*q[7]*q[3] +
		q[2]*q[4]*q[8]*q[3] - q[1]*q[5]*q[8]*q[3] - q[3]*q[5]*q[2]*q[6] + q[2]*q[6]*q[2]*q[6] +
		q[3]*q[1]*q[7]*q[6] - q[0]*q[6]*q[7]*q[6] - q[2]*q[1]*q[8]*q[6] + q[0]*q[5]*q[8]*q[6] +
		q[3]*q[4]*q[2]*q[8] - q[1]*q[6]*q[2]*q[8] - q[3]*q[1]*q[5]*q[8] + q[0]*q[6]*q[5]*q[8] +
		q[1]*q[1]*q[8]*q[8] - q[0]*q[4]*q[8]*q[8] - q[2]*q[4]*q[2]*q[9] + q[1]*q[5]*q[2]*q[9] +
		q[2]*q[1]*q[5]*q[9] - q[0]*q[5]*q[5]*q[9] - q[1]*q[1]*q[7]*q[9] + q[0]*q[4]*q[7]*q[9]
}

func (q *Quadric) VertexError(x, y, z float64) float64 {
	// v(transpose) * q * v
	return x*x*q[0] + 2*x*y*q[1] + 2*x*z*q[2] + 2*x*q[3] +
		y*y*q[4] + 2*y*z*q[5] + 2*y*q[6] +
		z*z*q[7] + 2*z*q[8] +
		q[9]
}

func (q *Quadric) Clone() *Quadric {
	new_q := Quadric{}
	copy(new_q[:], q[:])
	return &new_q
}

func (q *Quadric) Determinant2(
	a11, a12, a13,
	a21, a22, a23,
	a31, a32, a33 int) float64 {
	return q[a11]*q[a22]*q[a33] + q[a13]*q[a21]*q[a32] + q[a12]*q[a23]*q[a31] -
		q[a13]*q[a22]*q[a31] - q[a11]*q[a23]*q[a32] - q[a12]*q[a21]*q[a33]

}
