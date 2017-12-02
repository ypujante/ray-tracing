package main

import "testing"

func TestVec3_Scale(t *testing.T) {
	var tests = []struct {
		orig Vec3
		factor float64
		expected Vec3
	} {
		{Vec3{1.0, 2.0, 3.0}, 1.0, Vec3{1.0, 2.0, 3.0}},
		{Vec3{1.0, 2.0, 3.0}, 2.0, Vec3{2.0, 4.0, 6.0}},
		{Vec3{1.0, 2.0, 3.0}, 0.5, Vec3{0.5, 1.0, 1.5}},
		{Vec3{}, 1, Vec3{0, 0, 0}},
		{Vec3{1.0, 2.0, 3.0}, 0, Vec3{}},
	}

	for idx, test := range tests {
		s := test.orig.Scale(test.factor)
		if s.X != test.expected.X || s.Y != test.expected.Y || s.Z != test.expected.Z {
			t.Errorf("%v/%v/%v expected got %v/%v/%v instead [test %v]", test.expected.X, test.expected.Y, test.expected.Z, s.X, s.Y, s.Z, idx)
		}
	}
}

type RndMock struct {
	floats []float64
	idx int
}

func (rnd *RndMock) Float64() float64 {
	idx := rnd.idx
	rnd.idx++
	return rnd.floats[idx]
}

var EPSILON = 0.00000001

func floatEquals(a, b float64) bool {
	if (a - b) < EPSILON && (b - a) < EPSILON {
		return true
	}
	return false
}

func TestRandomInUnitSphere(t *testing.T) {
	rndMock := RndMock{floats: []float64{0.8, 0.7, 0.6}}

	r := randomInUnitSphere(&rndMock)

	if !floatEquals(r.X, 2.0 * 0.8 - 1.0) || !floatEquals(r.Y,  2.0 * 0.7 - 1.0) || !floatEquals(r.Z, 2.0 * 0.6 - 1.0) {
		t.Errorf("unexpected vector %v/%v/%v", r.X, r.Y, r.Z)
	}

	rndMock = RndMock{floats: []float64{0.99, 0.99, 0.99, 0.8, 0.7, 0.6}}

	r = randomInUnitSphere(&rndMock)

	if !floatEquals(r.X, 2.0 * 0.8 - 1.0) || !floatEquals(r.Y,  2.0 * 0.7 - 1.0) || !floatEquals(r.Z, 2.0 * 0.6 - 1.0) {
		t.Errorf("unexpected vector %v/%v/%v", r.X, r.Y, r.Z)
	}


}