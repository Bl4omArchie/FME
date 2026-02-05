package fme

import (
	"testing"
)


// Test schema + combination validation
func TestSchema(t *testing.T) {
	schema := NewSchema(map[string]Constraint{"req": &Require{}, "interfer": &Interfer{}})

	_ = schema.AddConstraint("req", "a", "b")
	_ =  schema.AddConstraint("req", "b", "c")
	err := schema.AddConstraint("interfer", "c", "a")
	if err != nil {
		t.Fatalf("incorrect schema : %v", err)
	}

	_, err = NewCombination([]string{"a", "b", "c"}, schema)
	if err != nil {
		t.Fatalf("incorrect combination : %v", err)
	}
}


// Test incorrect interference cases
// i.e : a->b ; b->c ; c-/->a
func TestInterference(t *testing.T) {
	schema := NewSchema(map[string]Constraint{"req": &Require{}, "interfer": &Interfer{}})

	_ = schema.AddConstraint("req", "a", "b")
	_ = schema.AddConstraint("req", "b", "c")
	if err := schema.AddConstraint("interfer", "c", "a"); err == nil {
		t.Fatalf("TestConflict #1 : interference was accepted")
	}
}


// Test cycle depedency
func TestCycleDependency(t *testing.T) {
	schema := NewSchema(map[string]Constraint{"req": &Require{}, "interfer": &Interfer{}})
	
	_ = schema.AddConstraint("req", "a", "b")
	if err := schema.AddConstraint("req", "b", "a"); err == nil {
		t.Fatalf("TestConflict #2 : cycle dependency was accepted")
	}
}
