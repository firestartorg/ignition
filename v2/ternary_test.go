package ignition

import "testing"

func TestTernary(t *testing.T) {
	if Ternary(true, 1, 2) != 1 {
		t.Error("Ternary(true, 1, 2) != 1")
	}
	if Ternary(false, 1, 2) != 2 {
		t.Error("Ternary(false, 1, 2) != 2")
	}

	if Ternary(true, "a", "b") != "a" {
		t.Error(`Ternary(true, "a", "b") != "a"`)
	}
	if Ternary(false, "a", "b") != "b" {
		t.Error(`Ternary(false, "a", "b") != "b"`)
	}

	if Ternary(true, true, false) != true {
		t.Error("Ternary(true, true, false) != true")
	}
	if Ternary(false, true, false) != false {
		t.Error("Ternary(false, true, false) != false")
	}
}
