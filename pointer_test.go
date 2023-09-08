package ignition

import "testing"

func TestDereference_Int(t *testing.T) {
	var val *int
	if Dereference(val) != 0 {
		t.Error("Dereference did not return the default value")
	}
	val = new(int)
	*val = 42
	if Dereference(val) != *val {
		t.Error("Dereference did not return the value")
	}
}

func TestDereference_String(t *testing.T) {
	var val *string
	if Dereference(val) != "" {
		t.Error("Dereference did not return the default value")
	}
	val = new(string)
	*val = "value"
	if Dereference(val) != *val {
		t.Error("Dereference did not return the value")
	}
}

func TestReference_Int(t *testing.T) {
	val := 42
	if *Reference(val) != val {
		t.Error("Reference did not return the value")
	}
}

func TestReference_String(t *testing.T) {
	val := "value"
	if *Reference(val) != val {
		t.Error("Reference did not return the value")
	}
}
