package segtree

import (
	"testing"
)

func TestLog2(t *testing.T) {
	test := func(num, expected int) {
		if result := log2(num); result != expected {
			t.Errorf("Log₂(%d) should be %d, got %d", num, expected, result)
		}
	}

	// Log₂ of a negative integer is impossible, but instead it returns -1.
	// I did not want to bother with NaN errors.
	test(-10, -1)
	// Log₂ of 0 is indefined, but it's limit aprouches -∞. It returs NegInf,
	// which is equal to the minimal value of an integer.
	// I did not want to bother with Inf errors.
	test(1, 0)
	test(2, 1)
	test(3, 1)
	test(4, 2)
	test(9, 3)
	test(10, 3)
	test(1024, 10)

}
