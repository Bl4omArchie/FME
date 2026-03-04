package fme

import (
	"testing"
)


// Test schema + combination validation
func TestSchema(t *testing.T) {
	schema := InitSchema()

	_ = schema.Add("req", "a", "b")
	_ =  schema.Add("req", "b", "c")
	err := schema.Add("interfer", "c", "a")
	if err == nil {
		t.Fatalf("invalid schema accepted : %v", err)
	}

	c := schema.ValidateCombination([]string{"a", "b", "c"})
	if c.Error == nil {
		t.Fatalf("invalid combination accepted : %v", err)
	}
}


// Test incorrect interference cases
// i.e : a->b ; b->c ; c-/->a
func TestInterference(t *testing.T) {
	schema := InitSchema()

	_ = schema.Add("req", "a", "b")
	_ = schema.Add("req", "b", "c")
	if err := schema.Add("interfer", "c", "a"); err == nil {
		t.Fatalf("TestConflict #1 : interference was accepted")
	}
}


// Test cycle depedency
func TestCycleDependency(t *testing.T) {
	schema := InitSchema()
	
	_ = schema.Add("req", "a", "b")
	if err := schema.Add("req", "b", "a"); err == nil {
		t.Fatalf("TestConflict #2 : cycle dependency was accepted")
	}
}
