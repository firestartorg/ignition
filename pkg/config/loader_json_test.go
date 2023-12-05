package config

import "testing"

func Test_unpackFixName(t *testing.T) {
	if fixName("test", true) != "Test" {
		t.Error("fixName did not return the correct value")
	}
	if fixName("Test", true) != "Test" {
		t.Error("fixName did not return the correct value")
	}
	if fixName("test-test", true) != "TestTest" {
		t.Error("fixName did not return the correct value")
	}
	if fixName("test_test", true) != "TestTest" {
		t.Error("fixName did not return the correct value")
	}

	if fixName("test-test", false) != "test-test" {
		t.Error("fixName did not return the correct value")
	}
	if fixName("test_test", false) != "test_test" {
		t.Error("fixName did not return the correct value")
	}
	if fixName("test", false) != "test" {
		t.Error("fixName did not return the correct value")
	}

}
