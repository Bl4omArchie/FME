package fme

import (
	"testing"
)

func TestSchema(t *testing.T) {
	schema := NewSchema()
	schema.AddConstraint(&Require{}, &Interfer{})

	schema.Constraints[0].Add("a", "b")
	schema.Constraints[0].Add("b", "c")
	schema.Constraints[1].Add("c", "a")


	if err := schema.ValidateSchema(); err != nil {
		t.Fatalf("%v", err)
	}
}
