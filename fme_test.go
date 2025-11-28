package fme

import (
	"testing"
)

// Test Conflict #1
//
// a -> b
// b -> c
// c -/-> a
func TestConflict(t *testing.T) {
	schema := NewSchema()
	schema.Require("a", "b")
	schema.Require("b", "c")
	if ok, _ := schema.Interfer("c", "a"); ok == true {
		t.Fatalf("TestConflict #1 : interference was accepted")
	}
}

// Test Conflict #2
//
// a -> b
// b -> a
func TestCycleDependency(t *testing.T) {
	schema := NewSchema()
	schema.Require("a", "b")

	if ok, _ := schema.Require("b", "a"); ok == true {
		t.Fatalf("TestConflict #2 : cycle dependency")
	}
	
}
