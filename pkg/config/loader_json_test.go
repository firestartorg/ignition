package config

import "testing"

func Test_unpackFixName(t *testing.T) {
	if unpackFixName("test", true) != "Test" {
		t.Error("unpackFixName did not return the correct value")
	}
	if unpackFixName("Test", true) != "Test" {
		t.Error("unpackFixName did not return the correct value")
	}
	if unpackFixName("test-test", true) != "TestTest" {
		t.Error("unpackFixName did not return the correct value")
	}
	if unpackFixName("test_test", true) != "TestTest" {
		t.Error("unpackFixName did not return the correct value")
	}

	if unpackFixName("test-test", false) != "test-test" {
		t.Error("unpackFixName did not return the correct value")
	}
	if unpackFixName("test_test", false) != "test_test" {
		t.Error("unpackFixName did not return the correct value")
	}
	if unpackFixName("test", false) != "test" {
		t.Error("unpackFixName did not return the correct value")
	}

}
