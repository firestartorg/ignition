package ignition

import "testing"

func TestCoalesce_Int(t *testing.T) {
	var val *int
	defVal := 42
	if Coalesce(val, defVal) != defVal {
		t.Error("Coalesce did not return the default value")
	}
	val = new(int)
	*val = 0
	if Coalesce(val, defVal) != *val {
		t.Error("Coalesce did not return the value")
	}
}

func TestCoalesce_String(t *testing.T) {
	var val *string
	defVal := "default"
	if Coalesce(val, defVal) != defVal {
		t.Error("Coalesce did not return the default value")
	}
	val = new(string)
	*val = "value"
	if Coalesce(val, defVal) != *val {
		t.Error("Coalesce did not return the value")
	}
}
