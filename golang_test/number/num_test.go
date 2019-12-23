package number

import (
	"testing"
)

func TestAdd(t *testing.T) {
	calc := Num(1)

	if calc.Add(3) != 4 {
		t.Errorf("%f + %f != %f", 1.0, 3.0, calc.Add(3))
	}
}
