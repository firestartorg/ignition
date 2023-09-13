package ignition

import "testing"

func TestCapitalizeString(t *testing.T) {
	if CapitalizeString("test") != "Test" {
		t.Error("CapitalizeString did not capitalize the string")
	}
	if CapitalizeString("fooBar") != "Foobar" {
		t.Error("CapitalizeString did not capitalize the string")
	}
	if CapitalizeString("TEST") != "Test" {
		t.Error("CapitalizeString did not lowercase the string")
	}
	if CapitalizeString(" Test ") != "Test" {
		t.Error("CapitalizeString did not trim the string")
	}
	if CapitalizeString("") != "" {
		t.Error("CapitalizeString did not return an empty string")
	}
}
