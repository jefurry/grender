package gobox

import (
	"testing"
)

func TestBoxMinMax(t *testing.T) {
	var x, y, z float64

	x = 4.352
	y = 3.25678
	z = -12.45242
	if z != min(x, y, z) {
		t.Errorf("min testing failed.")
	}
	if x != max(x, y, z) {
		t.Errorf("max testing failed.")
	}

	x = -43.135343
	y = 4.248103
	z = -3461.32
	if z != min(x, y, z) {
		t.Errorf("min testing failed.")
	}
	if y != max(x, y, z) {
		t.Errorf("max testing failed.")
	}

	x = 34.1246
	y = 125.09832
	z = 0.9876
	if z != min(x, y, z) {
		t.Errorf("min testing failed.")
	}
	if y != max(x, y, z) {
		t.Errorf("max testing failed.")
	}
}

func TestBoxRound(t *testing.T) {
	var v float64
	var r float64

	v = 123.567428
	r, _ = Round(v, 2)
	if 123.57 != r {
		t.Errorf("round testing failed.")
	}
	r, _ = Round(v, 4)
	if 123.5674 != r {
		t.Errorf("round testing failed.")
	}

	v = -453.126736
	r, _ = Round(v, 5)
	if -453.12674 != r {
		t.Errorf("round testing failed.")
	}
	r, _ = Round(v, 1)
	if -453.1 != r {
		t.Errorf("round testing failed.")
	}
	r, _ = Round(v, 2)
	if -453.13 != r {
		t.Errorf("round testing failed.")
	}
}
